package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/handlers"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/router"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/ws"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/tts"
)

func main() {
	// 加载配置
	configPath := os.Getenv("TTS_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/tts-config.json"
	}

	cfg, err := config.LoadTTSConfig(configPath)
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

	logger.Info("Starting TTS server...")

	// 创建TTS管理器
	poolSize := 5 // TTS默认池大小较小
	ttsManager, err := tts.NewManager(&cfg.TTS, poolSize)
	if err != nil {
		logger.Errorf("Failed to create TTS manager: %v", err)
		os.Exit(1)
	}
	defer ttsManager.Close()

	// 创建会话管理器
	sessionManager := session.NewManager(1000, 30*time.Minute)

	// 创建路由
	r := router.NewRouter()
	r.SetupMiddleware()
	r.SetupStaticFiles("web/static", "web/templates")

	// 创建处理器
	ttsHandler := handlers.NewTTSHandler(ttsManager, cfg)

	// 创建WebSocket升级器
	upgrader := ws.NewUpgrader(
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		54*time.Second,
		60*time.Second,
		int64(cfg.WebSocket.MaxMessageSize),
		cfg.WebSocket.EnableCompression,
	)

	// 创建TTS WebSocket处理器
	ttsWSHandler := ws.NewTTSHandler(sessionManager, ttsManager, cfg)

	// 设置路由
	r.SetupRoutes(func(ginEngine *gin.Engine) {
		// API路由
		api := ginEngine.Group("/api/v1")
		{
			// 健康检查
			api.GET("/health", handlers.HealthHandler(nil, &handlers.ProviderInfo{
				TTS:          cfg.TTS.Provider.Provider,
				GPUAvailable: cfg.TTS.Provider.Provider == "cuda",
				GPUDeviceID:  cfg.TTS.Provider.DeviceID,
			}))

			// TTS API
			ttsAPI := api.Group("/tts")
			{
				ttsAPI.POST("/synthesize", ttsHandler.Synthesize)
				ttsAPI.POST("/batch", ttsHandler.BatchSynthesize)
				ttsAPI.GET("/speakers", ttsHandler.GetSpeakers)
				ttsAPI.GET("/config", ttsHandler.GetConfig)
				ttsAPI.GET("/stats", ttsHandler.GetStats)
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
			ttsWSHandler.HandleConnection(conn)
		})

		// 静态页面
		ginEngine.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "tts-test.html", nil)
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
		logger.Infof("TTS server listening on %s", addr)
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

