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

func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer

	config := Config{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel)

	logger.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Debug message should be logged at debug level")
	}

	buf.Reset()
	logger.SetLevel(InfoLevel)
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not be logged at info level")
	}
}

func TestSetOutput(t *testing.T) {
	var buf1, buf2 bytes.Buffer

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
	logger.SetOutput(&buf1)
	logger.Info("message 1")

	if buf1.Len() == 0 {
		t.Error("Message should be written to buf1")
	}

	logger.SetOutput(&buf2)
	logger.Info("message 2")

	if buf2.Len() == 0 {
		t.Error("Message should be written to buf2")
	}
}

func TestDebugf(t *testing.T) {
	var buf bytes.Buffer

	config := Config{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	logger.SetOutput(&buf)
	logger.SetLevel(DebugLevel)

	logger.Debugf("test message: %s", "value")
	if buf.Len() == 0 {
		t.Error("Debugf message should be logged")
	}
}

func TestInfof(t *testing.T) {
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

	logger.Infof("test message: %d", 123)
	if buf.Len() == 0 {
		t.Error("Infof message should be logged")
	}
}

func TestWarnf(t *testing.T) {
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

	logger.Warnf("test message: %f", 1.23)
	if buf.Len() == 0 {
		t.Error("Warnf message should be logged")
	}
}

func TestErrorf(t *testing.T) {
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

	logger.Errorf("test message: %v", "error")
	if buf.Len() == 0 {
		t.Error("Errorf message should be logged")
	}
}

func TestInitLoggerWithFile(t *testing.T) {
	tmpFile := "/tmp/test-logger-file.log"
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

func TestInitLoggerWithBoth(t *testing.T) {
	tmpFile := "/tmp/test-logger-both.log"
	defer os.Remove(tmpFile)

	config := Config{
		Level:    "info",
		Format:   "text",
		Output:   "both",
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

func TestInitLoggerInvalidLevel(t *testing.T) {
	config := Config{
		Level:  "invalid",
		Format: "text",
		Output: "console",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() should not fail with invalid level: %v", err)
	}

	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger() returned nil")
	}
}

func TestEscapeJSON(t *testing.T) {
	// 测试JSON转义功能（通过日志输出）
	var buf bytes.Buffer

	config := Config{
		Level:  "info",
		Format: "json",
		Output: "console",
	}

	err := InitLogger(config)
	if err != nil {
		t.Fatalf("InitLogger() error = %v", err)
	}

	logger := GetLogger()
	logger.SetOutput(&buf)

	// 测试包含特殊字符的消息
	logger.Info("test message with \"quotes\" and\nnewline")
	
	output := buf.String()
	if len(output) == 0 {
		t.Error("JSON log should be written")
	}
	
	// 验证JSON格式
	if !contains(output, "{") || !contains(output, "timestamp") {
		t.Error("Output should be valid JSON")
	}
}

