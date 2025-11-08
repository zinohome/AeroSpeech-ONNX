package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/asr"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/handlers"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/router"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/ws"
)

func main() {
	// 加载配置
	configPath := os.Getenv("STT_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/stt-config.json"
	}

	cfg, err := config.LoadSTTConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(logger.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		FilePath:   cfg.Logging.FilePath,
		MaxSize:    cfg.Logging.MaxSize,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,
		Compress:   cfg.Logging.Compress,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting STT server...")

	// 创建ASR管理器
	poolSize := 4 // 默认池大小
	asrManager, err := asr.NewManager(&cfg.ASR, poolSize)
	if err != nil {
		logger.Errorf("Failed to create ASR manager: %v", err)
		os.Exit(1)
	}
	defer asrManager.Close()

	// 创建会话管理器
	sessionManager := session.NewManager(1000, 30*time.Minute)

	// 创建路由
	r := router.NewRouter()
	r.SetupMiddleware()
	r.SetupStaticFiles("web/static", "web/templates")

	// 创建处理器
	sttHandler := handlers.NewSTTHandler(asrManager, cfg)

	// 创建WebSocket升级器
	upgrader := ws.NewUpgrader(
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		54*time.Second,
		60*time.Second,
		int64(cfg.WebSocket.MaxMessageSize),
		cfg.WebSocket.EnableCompression,
	)

	// 创建STT WebSocket处理器
	sttWSHandler := ws.NewSTTHandler(sessionManager, asrManager, cfg)

	// 设置路由
	r.SetupRoutes(func(ginEngine *gin.Engine) {
		// API路由
		api := ginEngine.Group("/api/v1")
		{
			// 健康检查
			api.GET("/health", handlers.HealthHandler(nil, &handlers.ProviderInfo{
				ASR:          cfg.ASR.Provider.Provider,
				GPUAvailable: cfg.ASR.Provider.Provider == "cuda",
				GPUDeviceID:  cfg.ASR.Provider.DeviceID,
			}))

			// STT API
			stt := api.Group("/stt")
			{
				stt.POST("/recognize", sttHandler.Recognize)
				stt.POST("/batch", sttHandler.BatchRecognize)
				stt.GET("/config", sttHandler.GetConfig)
				stt.GET("/stats", sttHandler.GetStats)
			}

			// 统计信息
			api.GET("/stats", handlers.StatsHandler(nil))
			api.GET("/monitor", handlers.MonitorHandler(nil))
		}

		// WebSocket路由
		ginEngine.GET("/ws", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				logger.Errorf("WebSocket upgrade failed: %v", err)
				return
			}
			sttWSHandler.HandleConnection(conn)
		})

		// 静态页面
		ginEngine.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "stt-test.html", nil)
		})
	})

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r.GetEngine(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.ReadTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Infof("STT server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("Server stopped")
}

