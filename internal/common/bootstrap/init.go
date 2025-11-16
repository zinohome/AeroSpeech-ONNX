package bootstrap

import (
	"fmt"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/asr"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config/hotreload"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/middleware"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/tts"
	"time"
)

// AppDependencies 应用依赖
type AppDependencies struct {
	Config         *config.UnifiedConfig
	ASRManager     *asr.Manager
	TTSManager     *tts.Manager
	SessionManager *session.Manager
	RateLimiter    *middleware.RateLimiter
	HotReloadMgr   *hotreload.HotReloadManager
}

// InitApp 初始化应用
func InitApp(cfg *config.UnifiedConfig) (*AppDependencies, error) {
	logger.Info("Initializing components...")

	deps := &AppDependencies{
		Config: cfg,
	}

	// 初始化配置热加载管理器
	logger.Info("Initializing hot reload manager...")
	hotReloadMgr, err := hotreload.NewHotReloadManager()
	if err != nil {
		logger.Errorf("Failed to initialize hot reload manager: %v", err)
		// 热重载失败不影响主流程
	} else {
		deps.HotReloadMgr = hotReloadMgr
		// 注册配置变更回调
		registerHotReloadCallbacks(hotReloadMgr)
	}

	// 初始化限流器
	logger.Infof("Initializing rate limiter... enabled=%v, requests_per_second=%d, max_connections=%d",
		cfg.RateLimit.Enabled, cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.MaxConnections)
	rateLimiter := middleware.NewRateLimiter(
		cfg.RateLimit.Enabled,
		cfg.RateLimit.RequestsPerSecond,
		cfg.RateLimit.BurstSize,
		cfg.RateLimit.MaxConnections,
	)
	deps.RateLimiter = rateLimiter

	// 初始化ASR管理器
	if cfg.STT != nil {
		logger.Info("Initializing ASR manager...")
		poolSize := 4 // 默认池大小
		asrManager, err := asr.NewManager(cfg.STT, poolSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create ASR manager: %w", err)
		}
		deps.ASRManager = asrManager
		logger.Info("ASR manager initialized")
	}

	// 初始化TTS管理器
	if cfg.TTS != nil {
		logger.Info("Initializing TTS manager...")
		poolSize := 5 // TTS默认池大小较小
		ttsManager, err := tts.NewManager(cfg.TTS, poolSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create TTS manager: %w", err)
		}
		deps.TTSManager = ttsManager
		logger.Info("TTS manager initialized")
	}

	// 初始化会话管理器
	logger.Info("Initializing session manager...")
	sessionManager := session.NewManager(1000, 30*time.Minute)
	deps.SessionManager = sessionManager

	logger.Info("All components initialized successfully")
	return deps, nil
}

// registerHotReloadCallbacks 注册配置热加载回调
func registerHotReloadCallbacks(hotReloadMgr *hotreload.HotReloadManager) {
	if hotReloadMgr == nil {
		return
	}

	hotReloadMgr.RegisterCallback("logging.level", func() {
		logger.Info("Log level changed, reloading...")
	})
	hotReloadMgr.RegisterCallback("rate_limit", func() {
		logger.Info("Rate limit configuration changed")
	})
	hotReloadMgr.RegisterCallback("session", func() {
		logger.Info("Session configuration changed")
	})
	logger.Info("Hot reload callbacks registered")
}

// Close 关闭应用
func (d *AppDependencies) Close() error {
	logger.Info("Shutting down application...")

	// 关闭热重载管理器
	if d.HotReloadMgr != nil {
		d.HotReloadMgr.Stop()
	}

	// 关闭ASR管理器
	if d.ASRManager != nil {
		if err := d.ASRManager.Close(); err != nil {
			logger.Errorf("Failed to close ASR manager: %v", err)
		}
	}

	// 关闭TTS管理器
	if d.TTSManager != nil {
		if err := d.TTSManager.Close(); err != nil {
			logger.Errorf("Failed to close TTS manager: %v", err)
		}
	}

	logger.Info("Application shutdown complete")
	return nil
}

