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
	"github.com/zhangjun/AeroSpeech-ONNX/internal/asr"
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
	configPath := os.Getenv("SPEECH_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/speech-config.json"
	}

	cfg, err := config.LoadUnifiedConfig(configPath)
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

	logger.Infof("Starting speech server in %s mode...", cfg.Mode)

	// 创建管理器
	var asrManager *asr.Manager
	var ttsManager *tts.Manager

	// 初始化STT管理器
	if cfg.STT != nil {
		poolSize := 4 // 默认池大小
		asrManager, err = asr.NewManager(cfg.STT, poolSize)
		if err != nil {
			logger.Errorf("Failed to create ASR manager: %v", err)
			os.Exit(1)
		}
		defer asrManager.Close()
		logger.Info("ASR manager initialized")
	}

	// 初始化TTS管理器
	if cfg.TTS != nil {
		poolSize := 5 // TTS默认池大小较小
		ttsManager, err = tts.NewManager(cfg.TTS, poolSize)
		if err != nil {
			logger.Errorf("Failed to create TTS manager: %v", err)
			os.Exit(1)
		}
		defer ttsManager.Close()
		logger.Info("TTS manager initialized")
	}

	// 创建会话管理器
	sessionManager := session.NewManager(1000, 30*time.Minute)

	// 创建路由
	r := router.NewRouter()
	r.SetupMiddleware()
	r.SetupStaticFiles("web/static", "web/templates")

	// 创建处理器
	var sttHandler *handlers.STTHandler
	var ttsHandler *handlers.TTSHandler

	if asrManager != nil {
		// 创建STT配置（用于处理器）
		sttCfg := &config.STTConfig{
			Server:    cfg.Server,
			ASR:       *cfg.STT,
			Audio:     cfg.Audio,
			WebSocket: cfg.WebSocket,
			Session:   cfg.Session,
			Logging:   cfg.Logging,
		}
		sttHandler = handlers.NewSTTHandler(asrManager, sttCfg)
	}

	if ttsManager != nil {
		// 创建TTS配置（用于处理器）
		ttsCfg := &config.TTSConfig{
			Server:    cfg.Server,
			TTS:       *cfg.TTS,
			Audio:     cfg.Audio,
			WebSocket: cfg.WebSocket,
			Session:   cfg.Session,
			Logging:   cfg.Logging,
		}
		ttsHandler = handlers.NewTTSHandler(ttsManager, ttsCfg)
	}

	// 创建WebSocket升级器
	upgrader := ws.NewUpgrader(
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		time.Duration(cfg.WebSocket.ReadTimeout)*time.Second,
		54*time.Second,
		60*time.Second,
		int64(cfg.WebSocket.MaxMessageSize),
		cfg.WebSocket.EnableCompression,
	)

	// 创建WebSocket处理器
	var sttWSHandler *ws.STTHandler
	var ttsWSHandler *ws.TTSHandler

	if asrManager != nil {
		sttCfg := &config.STTConfig{
			Server:    cfg.Server,
			ASR:       *cfg.STT,
			Audio:     cfg.Audio,
			WebSocket: cfg.WebSocket,
			Session:   cfg.Session,
			Logging:   cfg.Logging,
		}
		sttWSHandler = ws.NewSTTHandler(sessionManager, asrManager, sttCfg)
	}

	if ttsManager != nil {
		ttsCfg := &config.TTSConfig{
			Server:    cfg.Server,
			TTS:       *cfg.TTS,
			Audio:     cfg.Audio,
			WebSocket: cfg.WebSocket,
			Session:   cfg.Session,
			Logging:   cfg.Logging,
		}
		ttsWSHandler = ws.NewTTSHandler(sessionManager, ttsManager, ttsCfg)
	}

	// 设置路由
	r.SetupRoutes(func(ginEngine *gin.Engine) {
		// API路由
		api := ginEngine.Group("/api/v1")
		{
			// 健康检查
			providerInfo := &handlers.ProviderInfo{}
			if asrManager != nil {
				providerInfo.ASR = cfg.STT.Provider.Provider
				providerInfo.GPUAvailable = cfg.STT.Provider.Provider == "cuda"
				providerInfo.GPUDeviceID = cfg.STT.Provider.DeviceID
			}
			if ttsManager != nil {
				providerInfo.TTS = cfg.TTS.Provider.Provider
				if !providerInfo.GPUAvailable {
					providerInfo.GPUAvailable = cfg.TTS.Provider.Provider == "cuda"
				}
				if providerInfo.GPUDeviceID == 0 {
					providerInfo.GPUDeviceID = cfg.TTS.Provider.DeviceID
				}
			}
			api.GET("/health", handlers.HealthHandler(nil, providerInfo))

			// STT API
			if sttHandler != nil {
				stt := api.Group("/stt")
				{
					stt.POST("/recognize", sttHandler.Recognize)
					stt.POST("/batch", sttHandler.BatchRecognize)
					stt.GET("/config", sttHandler.GetConfig)
					stt.GET("/stats", sttHandler.GetStats)
				}
			}

			// TTS API
			if ttsHandler != nil {
				ttsAPI := api.Group("/tts")
				{
					ttsAPI.POST("/synthesize", ttsHandler.Synthesize)
					ttsAPI.POST("/batch", ttsHandler.BatchSynthesize)
					ttsAPI.GET("/speakers", ttsHandler.GetSpeakers)
					ttsAPI.GET("/config", ttsHandler.GetConfig)
					ttsAPI.GET("/stats", ttsHandler.GetStats)
				}
			}

			// 统计信息
			api.GET("/stats", handlers.StatsHandler(nil))
			api.GET("/monitor", handlers.MonitorHandler(nil))
		}

		// WebSocket路由
		if sttWSHandler != nil {
			ginEngine.GET("/ws/stt", func(c *gin.Context) {
				conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					logger.Errorf("WebSocket upgrade failed: %v", err)
					return
				}
				sttWSHandler.HandleConnection(conn)
			})
		}

		if ttsWSHandler != nil {
			ginEngine.GET("/ws/tts", func(c *gin.Context) {
				conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					logger.Errorf("WebSocket upgrade failed: %v", err)
					return
				}
				ttsWSHandler.HandleConnection(conn)
			})
		}

		// 兼容旧的路由（统一模式）
		if cfg.Mode == "unified" {
			if sttWSHandler != nil {
				ginEngine.GET("/ws", func(c *gin.Context) {
					// 根据查询参数或路径判断是STT还是TTS
					serviceType := c.Query("type")
					if serviceType == "tts" && ttsWSHandler != nil {
						conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
						if err != nil {
							logger.Errorf("WebSocket upgrade failed: %v", err)
							return
						}
						ttsWSHandler.HandleConnection(conn)
					} else if sttWSHandler != nil {
						// 默认是STT
						conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
						if err != nil {
							logger.Errorf("WebSocket upgrade failed: %v", err)
							return
						}
						sttWSHandler.HandleConnection(conn)
					}
				})
			}
		}

		// 静态页面
		ginEngine.GET("/", func(c *gin.Context) {
			// 根据模式显示不同的页面
			if cfg.Mode == "unified" {
				// 统一模式：显示选择页面或默认STT页面
				c.HTML(http.StatusOK, "stt-test.html", nil)
			} else if asrManager != nil && ttsManager == nil {
				// 只有STT
				c.HTML(http.StatusOK, "stt-test.html", nil)
			} else if asrManager == nil && ttsManager != nil {
				// 只有TTS
				c.HTML(http.StatusOK, "tts-test.html", nil)
			} else {
				// 默认STT页面
				c.HTML(http.StatusOK, "stt-test.html", nil)
			}
		})

		// STT测试页面
		if asrManager != nil {
			ginEngine.GET("/stt", func(c *gin.Context) {
				c.HTML(http.StatusOK, "stt-test.html", nil)
			})
		}

		// TTS测试页面
		if ttsManager != nil {
			ginEngine.GET("/tts", func(c *gin.Context) {
				c.HTML(http.StatusOK, "tts-test.html", nil)
			})
		}

		// 监控面板
		ginEngine.GET("/monitor", func(c *gin.Context) {
			c.HTML(http.StatusOK, "monitor.html", nil)
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
		logger.Infof("Speech server listening on %s (mode: %s)", addr, cfg.Mode)
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

