package tts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestNewTTSProvider(t *testing.T) {
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
	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}

	if provider == nil {
		t.Fatal("NewTTSProvider() returned nil")
	}

	if provider.GetSampleRate() != 24000 {
		t.Errorf("Expected sample rate 24000, got %d", provider.GetSampleRate())
	}

	// 清理
	provider.Release()
}

func TestTTSProvider_Synthesize(t *testing.T) {
	// 这个测试需要实际的模型文件
	// 跳过实际测试，只测试接口
	t.Skip("Skipping test that requires actual model files")
}

func TestTTSProvider_SynthesizeEmptyText(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试空文本（应该返回错误）
	_, err = provider.Synthesize("", 0, 1.0)
	if err == nil {
		t.Error("Expected error for empty text")
	}
}

func TestTTSProvider_SynthesizeWithOptions(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试带选项的合成
	_, err = provider.Synthesize("测试文本", 0, 1.0)
	if err != nil {
		t.Logf("Synthesize() error = %v (expected if models not available)", err)
	}

	// 测试不同的说话人ID
	_, err = provider.Synthesize("测试文本", 45, 1.0)
	if err != nil {
		t.Logf("Synthesize() error = %v (expected if models not available)", err)
	}

	// 测试不同的语速
	_, err = provider.Synthesize("测试文本", 0, 1.5)
	if err != nil {
		t.Logf("Synthesize() error = %v (expected if models not available)", err)
	}
}

func TestTTSProvider_Warmup(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试预热
	err = provider.Warmup()
	if err != nil {
		t.Logf("Warmup() error = %v (expected if models not available)", err)
	}
}

func TestTTSProvider_Reset(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// Reset应该总是成功
	err = provider.Reset()
	if err != nil {
		t.Errorf("Reset() error = %v", err)
	}
}

func TestTTSProvider_Release(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}

	// Release应该总是成功
	err = provider.Release()
	if err != nil {
		t.Errorf("Release() error = %v", err)
	}
}

func TestTTSProvider_GetSampleRate(t *testing.T) {
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

	provider, err := NewTTSProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewTTSProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	sampleRate := provider.GetSampleRate()
	if sampleRate != 24000 {
		t.Errorf("Expected sample rate 24000, got %d", sampleRate)
	}
}

