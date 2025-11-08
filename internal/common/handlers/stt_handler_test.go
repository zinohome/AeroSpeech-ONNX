package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/pkg/utils"
)

// mockSTTManager 模拟STT管理器
type mockSTTManager struct {
	transcribeResult string
	transcribeError  error
	stats            interface{}
	avgLatency       interface{}
	poolUsage        float64
	poolStats        map[string]interface{}
}

func (m *mockSTTManager) Transcribe(ctx interface{}, audio []byte) (string, error) {
	if m.transcribeError != nil {
		return "", m.transcribeError
	}
	return m.transcribeResult, nil
}

func (m *mockSTTManager) GetStats() interface{} {
	return m.stats
}

func (m *mockSTTManager) GetAvgLatency() interface{} {
	return m.avgLatency
}

func (m *mockSTTManager) GetPoolUsage() float64 {
	return m.poolUsage
}

func (m *mockSTTManager) GetPoolStats() map[string]interface{} {
	return m.poolStats
}

func TestSTTHandler_Recognize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{
		transcribeResult: "测试文本",
		transcribeError:  nil,
	}

	cfg := &config.STTConfig{
		Audio: config.AudioConfig{
			SampleRate: 16000,
		},
	}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.POST("/recognize", handler.Recognize)

	// 创建multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", "test.wav")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("fake audio data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/recognize", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSTTHandler_RecognizeInvalidFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{}
	cfg := &config.STTConfig{}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.POST("/recognize", handler.Recognize)

	req := httptest.NewRequest("POST", "/recognize", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSTTHandler_RecognizeError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{
		transcribeError: utils.NewAppError(
			utils.ErrCodeRecognitionError,
			"recognition failed",
			"test error",
			nil,
		),
	}

	cfg := &config.STTConfig{}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.POST("/recognize", handler.Recognize)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", "test.wav")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("fake audio data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/recognize", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestSTTHandler_BatchRecognize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{
		transcribeResult: "测试文本",
	}

	cfg := &config.STTConfig{}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.POST("/batch", handler.BatchRecognize)

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test-audio-*.wav")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte("fake audio data"))
	tmpFile.Close()

	reqBody := `{"files": ["` + tmpFile.Name() + `"]}`
	req := httptest.NewRequest("POST", "/batch", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSTTHandler_BatchRecognizeInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{}
	cfg := &config.STTConfig{}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.POST("/batch", handler.BatchRecognize)

	req := httptest.NewRequest("POST", "/batch", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSTTHandler_GetConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{}
	cfg := &config.STTConfig{
		Audio: config.AudioConfig{
			SampleRate: 16000,
			ChunkSize:  4096,
		},
		ASR: config.ASRConfig{
			Provider: config.ProviderConfig{
				Provider:   "cpu",
				DeviceID:   0,
				NumThreads: 4,
			},
		},
	}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.GET("/config", handler.GetConfig)

	req := httptest.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSTTHandler_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockSTTManager{
		stats: map[string]interface{}{
			"total_requests": 100,
		},
		avgLatency: 50.0,
		poolUsage:  60.0,
		poolStats: map[string]interface{}{
			"active": 5,
			"total":  10,
		},
	}

	cfg := &config.STTConfig{}

	handler := NewSTTHandler(manager, cfg)

	router := gin.New()
	router.GET("/stats", handler.GetStats)

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

