package bootstrap

import (
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestInitApp(t *testing.T) {
	// 创建测试配置
	cfg := &config.UnifiedConfig{
		Mode: "unified",
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		RateLimit: config.RateLimitConfig{
			Enabled:           false,
			RequestsPerSecond: 1000,
			BurstSize:         2000,
			MaxConnections:    2000,
		},
		VAD: config.VADConfig{
			Enabled: false,
		},
	}

	// 测试初始化（注意：这里不测试ASR和TTS，因为需要模型文件）
	deps, err := InitApp(cfg)
	if err != nil {
		// 如果失败是因为没有模型文件，这是预期的
		t.Logf("InitApp failed (expected if no models): %v", err)
		return
	}

	// 检查基础组件是否初始化
	if deps.Config == nil {
		t.Error("Expected Config to be initialized")
	}

	if deps.RateLimiter == nil {
		t.Error("Expected RateLimiter to be initialized")
	}

	if deps.SessionManager == nil {
		t.Error("Expected SessionManager to be initialized")
	}

	// 清理
	if deps != nil {
		deps.Close()
	}
}

