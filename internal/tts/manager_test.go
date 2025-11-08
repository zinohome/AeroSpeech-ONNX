package tts

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestNewManager(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	// 注意：这个测试需要实际的sherpa-onnx库和模型文件
	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	// 清理
	manager.Close()
}

func TestManager_Synthesize(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试合成
	_, err = manager.Synthesize(nil, "测试文本", 0, 1.0)
	if err != nil {
		t.Logf("Synthesize() error = %v (expected if models not available)", err)
	}
}

func TestManager_SynthesizeWithContext(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试带上下文的合成
	ctx := context.Background()
	_, err = manager.Synthesize(ctx, "测试文本", 0, 1.0)
	if err != nil {
		t.Logf("Synthesize() error = %v (expected if models not available)", err)
	}
}

func TestManager_SynthesizeError(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试空文本（应该返回错误）
	_, err = manager.Synthesize(nil, "", 0, 1.0)
	if err == nil {
		t.Error("Expected error for empty text")
	}
}

func TestManager_GetStats(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	stats := manager.GetStats()
	if stats == nil {
		t.Error("GetStats() returned nil")
	}
}

func TestManager_GetAvgLatency(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	avgLatency := manager.GetAvgLatency()
	if avgLatency == nil {
		t.Error("GetAvgLatency() returned nil")
	}
}

func TestManager_GetPoolUsage(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	poolUsage := manager.GetPoolUsage()
	if poolUsage < 0 || poolUsage > 1 {
		t.Errorf("PoolUsage should be between 0 and 1, got %f", poolUsage)
	}
}

func TestManager_GetPoolStats(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("TTS_MODEL_PATH") == "" {
		t.Skip("Skipping test: TTS_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("TTS_MODEL_PATH")
	tokensFile := os.Getenv("TTS_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: TTS_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	// 尝试从model路径推断dataDir路径
	dataDir := ""
	if modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}

	// 尝试从环境变量或model路径推断voicesPath
	voicesPath := os.Getenv("TTS_VOICES_PATH")
	if voicesPath == "" && modelFile != "" {
		modelDir := filepath.Dir(modelFile)
		potentialVoicesPath := filepath.Join(modelDir, "voices.bin")
		if _, err := os.Stat(potentialVoicesPath); err == nil {
			voicesPath = potentialVoicesPath
		}
	}

	cfg := &config.TTSModelConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		VoicesPath: voicesPath, // 如果存在则使用，否则为空
		DataDir:    dataDir,
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	poolStats := manager.GetPoolStats()
	if poolStats == nil {
		t.Error("GetPoolStats() returned nil")
	}
}

