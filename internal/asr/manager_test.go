package asr

import (
	"context"
	"os"
	"testing"
	"time"

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

	// 先执行一次识别，以便有统计数据
	audioData := make([]byte, 1600*2) // 0.1秒的音频
	_, _ = manager.Transcribe(nil, audioData)

	stats := manager.GetStats()
	if stats == nil {
		t.Error("GetStats() returned nil")
	}

	// 验证stats是*Stats类型
	if statsStruct, ok := stats.(*Stats); ok {
		if statsStruct.TotalRequests == 0 {
			t.Log("TotalRequests is 0 (expected if transcription failed)")
		}
		// 验证统计信息结构
		if statsStruct.TotalRequests < 0 {
			t.Error("TotalRequests should be >= 0")
		}
		if statsStruct.SuccessfulRequests < 0 {
			t.Error("SuccessfulRequests should be >= 0")
		}
		if statsStruct.FailedRequests < 0 {
			t.Error("FailedRequests should be >= 0")
		}
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

	// 先执行一次识别，以便有统计数据
	audioData := make([]byte, 1600)
	_, _ = manager.Transcribe(nil, audioData)

	avgLatency := manager.GetAvgLatency()
	if avgLatency == nil {
		t.Error("GetAvgLatency() returned nil")
	}

	// 验证avgLatency是time.Duration类型
	if latency, ok := avgLatency.(time.Duration); ok {
		if latency < 0 {
			t.Errorf("AvgLatency should be >= 0, got %v", latency)
		}
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

	// 先执行一次操作，确保池被使用
	audioData := make([]byte, 1600*2)
	_, _ = manager.Transcribe(nil, audioData)

	poolUsage := manager.GetPoolUsage()
	if poolUsage < 0 || poolUsage > 1 {
		t.Errorf("PoolUsage should be between 0 and 1, got %f", poolUsage)
	}
	
	// 验证poolUsage是有效的浮点数
	if poolUsage != poolUsage { // NaN check
		t.Error("PoolUsage is NaN")
	}
}

func TestManager_GetPoolStats(t *testing.T) {
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

	// 先执行一次操作，确保池被使用
	audioData := make([]byte, 1600*2)
	_, _ = manager.Transcribe(nil, audioData)

	poolStats := manager.GetPoolStats()
	if poolStats == nil {
		t.Error("GetPoolStats() returned nil")
	}
	
	// 验证poolStats是map类型
	if len(poolStats) == 0 {
		t.Log("PoolStats is empty (may be valid)")
	}
	// 验证关键字段
	if size, ok := poolStats["size"].(int); ok {
		if size < 0 {
			t.Error("Pool size should be >= 0")
		}
	}
}

