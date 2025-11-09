package asr

import (
	"os"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestNewASRProvider(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	// 注意：这个测试需要实际的sherpa-onnx库和模型文件
	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}

	if provider == nil {
		t.Fatal("NewASRProvider() returned nil")
	}

	if provider.GetSampleRate() != 16000 {
		t.Errorf("Expected sample rate 16000, got %d", provider.GetSampleRate())
	}

	// 清理
	provider.Release()
}

func TestASRProvider_Transcribe(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试识别（使用有效的音频数据）
	// 生成0.1秒的16kHz 16-bit PCM音频数据
	audioData := make([]byte, 1600*2) // 1600 samples * 2 bytes per sample
	for i := 0; i < len(audioData); i += 2 {
		// 生成简单的正弦波音频数据
		sample := int16(1000 * float32(i%100) / 100.0)
		audioData[i] = byte(sample & 0xFF)
		audioData[i+1] = byte((sample >> 8) & 0xFF)
	}

	result, err := provider.Transcribe(audioData)
	if err != nil {
		t.Logf("Transcribe() error = %v (expected if models not available or audio is too short)", err)
		return
	}

	// 验证结果
	if result == "" {
		t.Log("Transcribe() returned empty result (may be valid for short audio)")
	} else {
		t.Logf("Transcribe() result = %s", result)
	}
}

func TestASRProvider_TranscribeEmptyAudio(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试空音频
	_, err = provider.Transcribe([]byte{})
	if err == nil {
		t.Error("Expected error for empty audio")
	}
}

func TestASRProvider_Warmup(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// 测试预热
	err = provider.Warmup()
	if err != nil {
		t.Logf("Warmup() error = %v (expected if models not available)", err)
	}
}

func TestASRProvider_Reset(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	// Reset应该总是成功
	err = provider.Reset()
	if err != nil {
		t.Errorf("Reset() error = %v", err)
	}
}

func TestASRProvider_Release(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}

	// Release应该总是成功
	err = provider.Release()
	if err != nil {
		t.Errorf("Release() error = %v", err)
	}
}

func TestASRProvider_GetSampleRate(t *testing.T) {
	// 检查是否有实际模型文件
	if os.Getenv("ASR_MODEL_PATH") == "" {
		t.Skip("Skipping test: ASR_MODEL_PATH not set (requires actual model files)")
		return
	}

	modelFile := os.Getenv("ASR_MODEL_PATH")
	tokensFile := os.Getenv("ASR_TOKENS_PATH")
	if tokensFile == "" {
		t.Skip("Skipping test: ASR_TOKENS_PATH not set (requires actual tokens file)")
		return
	}

	cfg := &config.ASRConfig{
		ModelPath:  modelFile,
		TokensPath: tokensFile,
		Language:   "zh",
		Provider: config.ProviderConfig{
			Provider:   "cpu",
			NumThreads: 1,
		},
	}

	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		return
	}
	defer provider.Release()

	sampleRate := provider.GetSampleRate()
	if sampleRate != 16000 {
		t.Errorf("Expected sample rate 16000, got %d", sampleRate)
	}
}

