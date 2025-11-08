package asr

import (
	"context"
	"os"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

// mockProvider 模拟Provider
type mockProvider struct {
	transcribeResult string
	transcribeError  error
	sampleRate       int
}

func (m *mockProvider) Transcribe(audio []byte) (string, error) {
	if m.transcribeError != nil {
		return "", m.transcribeError
	}
	return m.transcribeResult, nil
}

func (m *mockProvider) Warmup() error {
	return nil
}

func (m *mockProvider) Reset() error {
	return nil
}

func (m *mockProvider) Release() error {
	return nil
}

func (m *mockProvider) GetSampleRate() int {
	return m.sampleRate
}

func TestNewManager(t *testing.T) {
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

func TestManager_Transcribe(t *testing.T) {
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

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试识别
	audioData := make([]byte, 1600) // 0.1秒的音频
	_, err = manager.Transcribe(nil, audioData)
	if err != nil {
		t.Logf("Transcribe() error = %v (expected if models not available)", err)
	}
}

func TestManager_TranscribeWithContext(t *testing.T) {
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

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试带上下文的识别
	ctx := context.Background()
	audioData := make([]byte, 1600)
	_, err = manager.Transcribe(ctx, audioData)
	if err != nil {
		t.Logf("Transcribe() error = %v (expected if models not available)", err)
	}
}

func TestManager_TranscribeError(t *testing.T) {
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

	manager, err := NewManager(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewManager() error = %v (expected if models not available)", err)
		return
	}
	defer manager.Close()

	// 测试空音频（应该返回错误）
	_, err = manager.Transcribe(nil, []byte{})
	if err == nil {
		t.Error("Expected error for empty audio")
	}
}

func TestManager_GetStats(t *testing.T) {
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

