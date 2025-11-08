package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	config := Config{
		Level:    "info",
		Format:   "text",
		Output:   "console",
		FilePath: "",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger() returned nil")
	}
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer

	config := Config{
		Level:  "info",
		Format: "text",
		Output: "console",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	logger.SetOutput(&buf)

	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not be logged at info level")
	}

	buf.Reset()
	logger.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info message should be logged at info level")
	}

	buf.Reset()
	logger.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn message should be logged at info level")
	}

	buf.Reset()
	logger.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error message should be logged at info level")
	}
}

func TestLogFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		check  func(string) bool
	}{
		{
			name:   "text format",
			format: "text",
			check: func(s string) bool {
				return len(s) > 0 && !contains(s, "{")
			},
		},
		{
			name:   "json format",
			format: "json",
			check: func(s string) bool {
				return contains(s, "{") && contains(s, "timestamp")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			config := Config{
				Level:  "info",
				Format: tt.format,
				Output: "console",
			}

			err := InitLogger(config)
			if err != nil {
				t.Fatalf("InitLogger() error = %v", err)
			}

			logger := GetLogger()
			logger.SetOutput(&buf)

			logger.Info("test message")

			if !tt.check(buf.String()) {
				t.Errorf("Log format check failed for %s", tt.format)
			}
		})
	}
}

func TestFileOutput(t *testing.T) {
	tmpFile := "/tmp/test-logger.log"
	defer os.Remove(tmpFile)

	config := Config{
		Level:    "info",
		Format:   "text",
		Output:   "file",
		FilePath: tmpFile,
		MaxSize:  100,
		MaxBackups: 5,
		MaxAge:   30,
		Compress: false,
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	logger.Info("test message")

	// 检查文件是否存在
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

