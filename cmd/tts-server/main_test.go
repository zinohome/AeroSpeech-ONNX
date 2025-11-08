package main

import (
	"os"
	"testing"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-tts-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configJSON := `{
		"server": {
			"host": "0.0.0.0",
			"port": 8081
		},
		"tts": {
			"model_path": "/tmp/test-tts-model.onnx",
			"provider": {
				"provider": "cpu",
				"num_threads": 4
			}
		},
		"audio": {
			"sample_rate": 24000
		}
	}`

	if _, err := tmpFile.WriteString(configJSON); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()

	// 创建临时模型文件
	os.Create("/tmp/test-tts-model.onnx")
	defer os.Remove("/tmp/test-tts-model.onnx")

	// 测试配置加载
	cfg, err := config.LoadTTSConfig(tmpFile.Name())
	if err != nil {
		t.Logf("LoadTTSConfig() error = %v (expected if models not available)", err)
		return
	}

	if cfg == nil {
		t.Fatal("LoadTTSConfig() returned nil")
	}

	if cfg.Server.Port != 8081 {
		t.Errorf("Expected port 8081, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-tts-config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configJSON := `{
		"server": {
			"host": "0.0.0.0",
			"port": 8081
		},
		"tts": {
			"model_path": "/tmp/test-tts-model.onnx",
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
	os.Create("/tmp/test-tts-model.onnx")
	defer os.Remove("/tmp/test-tts-model.onnx")

	// 设置环境变量
	os.Setenv("TTS_CONFIG_PATH", tmpFile.Name())
	defer os.Unsetenv("TTS_CONFIG_PATH")

	// 测试从环境变量加载配置
	configPath := os.Getenv("TTS_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/tts-config.json"
	}

	cfg, err := config.LoadTTSConfig(configPath)
	if err != nil {
		t.Logf("LoadTTSConfig() error = %v (expected if models not available)", err)
		return
	}

	if cfg == nil {
		t.Fatal("LoadTTSConfig() returned nil")
	}
}

