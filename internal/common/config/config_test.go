package config

import (
	"os"
	"testing"
)

func TestLoadSTTConfig(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "test-config-*.json")
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

	config, err := LoadSTTConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Server.Port)
	}

	if config.ASR.Provider.Provider != "cpu" {
		t.Errorf("Expected provider cpu, got %s", config.ASR.Provider.Provider)
	}
}

func TestResolveProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		want     string
	}{
		{
			name:     "cpu provider",
			provider: "cpu",
			want:     "cpu",
		},
		{
			name:     "auto provider without GPU",
			provider: "auto",
			want:     "cpu", // 假设没有GPU
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &ProviderConfig{
				Provider:   tt.provider,
				NumThreads: 4,
			}

			if err := resolveProvider(provider); err != nil {
				t.Fatalf("resolveProvider() error = %v", err)
			}

			// 如果没有GPU，auto应该回退到cpu
			if !isGPUAvailable() && tt.provider == "auto" {
				if provider.Provider != "cpu" {
					t.Errorf("Expected cpu, got %s", provider.Provider)
				}
			}
		})
	}
}

func TestValidateSTTConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *STTConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &STTConfig{
				ASR: ASRConfig{
					ModelPath:  "/tmp/test-model.onnx",
					TokensPath: "/tmp/test-tokens.txt",
					Provider: ProviderConfig{
						Provider: "cpu",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing model path",
			config: &STTConfig{
				ASR: ASRConfig{
					TokensPath: "/tmp/test-tokens.txt",
					Provider: ProviderConfig{
						Provider: "cpu",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: &STTConfig{
				ASR: ASRConfig{
					ModelPath:  "/tmp/test-model.onnx",
					TokensPath: "/tmp/test-tokens.txt",
					Provider: ProviderConfig{
						Provider: "invalid",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件（如果需要）
			if tt.config.ASR.ModelPath != "" {
				os.Create(tt.config.ASR.ModelPath)
				defer os.Remove(tt.config.ASR.ModelPath)
			}
			if tt.config.ASR.TokensPath != "" {
				os.Create(tt.config.ASR.TokensPath)
				defer os.Remove(tt.config.ASR.TokensPath)
			}

			err := validateSTTConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSTTConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

