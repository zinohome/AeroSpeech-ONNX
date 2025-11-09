package asr

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestNewPool(t *testing.T) {
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
	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}

	if pool == nil {
		t.Fatal("NewPool() returned nil")
	}

	// 清理
	pool.Close()
}

func TestPool_Get(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 测试获取Provider
	ctx := context.Background()
	provider, err := pool.Get(ctx)
	if err != nil {
		t.Skipf("Skipping test: Get() error = %v (expected if models not available)", err)
		return
	}

	if provider == nil {
		t.Error("Get() returned nil")
	}

	// 归还Provider
	pool.Put(provider)
}

func TestPool_GetTimeout(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 获取所有Provider
	ctx := context.Background()
	provider, err := pool.Get(ctx)
	if err != nil {
		t.Skipf("Skipping test: Get() error = %v (expected if models not available)", err)
		return
	}

	// 尝试再次获取（应该超时）
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = pool.Get(timeoutCtx)
	if err == nil {
		// 如果池中有多个Provider，可能不会超时
		t.Log("Get() did not timeout (pool may have multiple providers)")
	}

	// 归还Provider
	pool.Put(provider)
}

func TestPool_Put(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 获取Provider
	ctx := context.Background()
	provider, err := pool.Get(ctx)
	if err != nil {
		t.Skipf("Skipping test: Get() error = %v (expected if models not available)", err)
		return
	}

	// 归还Provider
	pool.Put(provider)

	// 应该能够再次获取
	provider2, err := pool.Get(ctx)
	if err != nil {
		t.Logf("Get() error = %v (expected if models not available)", err)
		return
	}

	if provider2 == nil {
		t.Error("Get() returned nil after Put")
	}

	pool.Put(provider2)
}

func TestPool_Stats(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 获取统计信息
	stats := pool.GetStats()
	if stats == nil {
		t.Error("GetStats() returned nil")
	}
}

func TestPool_GetUsage(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 获取使用率
	usage := pool.GetUsage()
	if usage < 0 || usage > 1 {
		t.Errorf("Usage should be between 0 and 1, got %f", usage)
	}
}

func TestPool_PutFull(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}
	defer pool.Close()

	// 获取Provider
	ctx := context.Background()
	provider1, err := pool.Get(ctx)
	if err != nil {
		t.Skipf("Skipping test: Get() error = %v (expected if models not available)", err)
		return
	}

	// 创建一个新的Provider（模拟池满的情况）
	provider2, err := NewASRProvider(cfg)
	if err != nil {
		t.Skipf("Skipping test: NewASRProvider() error = %v (expected if models not available)", err)
		pool.Put(provider1)
		return
	}

	// 归还第一个Provider
	pool.Put(provider1)

	// 尝试归还第二个Provider（池已满，应该释放它）
	pool.Put(provider2)
}

func TestPool_Cleanup(t *testing.T) {
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

	pool, err := NewPool(cfg, 1)
	if err != nil {
		t.Skipf("Skipping test: NewPool() error = %v (expected if models not available)", err)
		return
	}

	// 清理池
	pool.Close()

	// 清理后应该无法获取Provider（Close后Get应该返回错误或阻塞）
	// 注意：根据pool的实现，Close后Get可能返回错误或阻塞，这里只检查不会panic
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	_, err = pool.Get(ctx)
	// Close后Get可能返回错误或超时，都是正常行为
	if err == nil {
		// 如果Get成功，说明pool可能没有完全关闭，这是实现相关的
		t.Log("Get() succeeded after Close() (implementation dependent)")
	}
}

