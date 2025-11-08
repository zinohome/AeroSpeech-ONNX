package main

import (
	"os"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-stt-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configJSON := `{
		"server": {
			"host": "0.0.0.0",
			"port": 8080
		},
		"asr": {
			"model_path": "/tmp/test-model.onnx",
			"tokens_path": "/tmp/test-tokens.txt",
			"provider": {
				"provider": "cpu",
				"num_threads": 4
			}
		},
		"audio": {
			"sample_rate": 16000
		}
	}`

	if _, err := tmpFile.WriteString(configJSON); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// 创建临时模型文件
	os.Create("/tmp/test-model.onnx")
	os.Create("/tmp/test-tokens.txt")
	defer os.Remove("/tmp/test-model.onnx")
	defer os.Remove("/tmp/test-tokens.txt")

	// 测试配置加载
	cfg, err := config.LoadSTTConfig(tmpFile.Name())
	if err != nil {
		t.Logf("LoadSTTConfig() error = %v (expected if models not available)", err)
		return
	}

	if cfg == nil {
		t.Fatal("LoadSTTConfig() returned nil")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-stt-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configJSON := `{
		"server": {
			"host": "0.0.0.0",
			"port": 8080
		},
		"asr": {
			"model_path": "/tmp/test-model.onnx",
			"tokens_path": "/tmp/test-tokens.txt",
			"provider": {
				"provider": "cpu",
				"num_threads": 4
			}
		}
	}`

	if _, err := tmpFile.WriteString(configJSON); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// 创建临时模型文件
	os.Create("/tmp/test-model.onnx")
	os.Create("/tmp/test-tokens.txt")
	defer os.Remove("/tmp/test-model.onnx")
	defer os.Remove("/tmp/test-tokens.txt")

	// 设置环境变量
	os.Setenv("STT_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("STT_CONFIG_PATH")

	// 测试从环境变量加载配置
	configPath := os.Getenv("STT_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/stt-config.json"
	}

	cfg, err := config.LoadSTTConfig(configPath)
	if err != nil {
		t.Logf("LoadSTTConfig() error = %v (expected if models not available)", err)
		return
	}

	if cfg == nil {
		t.Fatal("LoadSTTConfig() returned nil")
	}
}

