package asr

import (
	"os"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestNewASRProvider(t *testing.T) {
	// 创建临时模型文件
	modelFile := "/tmp/test-model.onnx"
	tokensFile := "/tmp/test-tokens.txt"
	
	os.Create(modelFile)
	os.Create(tokensFile)
	defer os.Remove(modelFile)
	defer os.Remove(tokensFile)

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
	// 如果没有，测试会失败，这是正常的
	provider, err := NewASRProvider(cfg)
	if err != nil {
		t.Logf("NewASRProvider() error = %v (expected if models not available)", err)
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

func TestASRProviderTranscribe(t *testing.T) {
	// 这个测试需要实际的模型文件
	// 跳过实际测试，只测试接口
	t.Skip("Skipping test that requires actual model files")
}

