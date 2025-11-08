package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/pkg/utils"
)

// mockTTSManager 模拟TTS管理器
type mockTTSManager struct {
	synthesizeResult []byte
	synthesizeError  error
	stats            interface{}
	avgLatency       interface{}
	poolUsage        float64
	poolStats        map[string]interface{}
}

func (m *mockTTSManager) Synthesize(ctx interface{}, text string, speakerID int, speed float32) ([]byte, error) {
	if m.synthesizeError != nil {
		return nil, m.synthesizeError
	}
	return m.synthesizeResult, nil
}

func (m *mockTTSManager) GetStats() interface{} {
	return m.stats
}

func (m *mockTTSManager) GetAvgLatency() interface{} {
	return m.avgLatency
}

func (m *mockTTSManager) GetPoolUsage() float64 {
	return m.poolUsage
}

func (m *mockTTSManager) GetPoolStats() map[string]interface{} {
	return m.poolStats
}

func TestTTSHandler_Synthesize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{
		synthesizeResult: []byte("fake audio data"),
		synthesizeError:  nil,
	}

	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/synthesize", handler.Synthesize)

	reqBody := `{"text": "测试文本", "speaker_id": 0, "speed": 1.0}`
	req := httptest.NewRequest("POST", "/synthesize", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTTSHandler_SynthesizeInvalidText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{}
	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/synthesize", handler.Synthesize)

	req := httptest.NewRequest("POST", "/synthesize", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTTSHandler_SynthesizeEmptyText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{}
	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/synthesize", handler.Synthesize)

	reqBody := `{"text": ""}`
	req := httptest.NewRequest("POST", "/synthesize", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTTSHandler_SynthesizeError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{
		synthesizeError: utils.NewAppError(
			utils.ErrCodeSynthesisError,
			"synthesis failed",
			"test error",
			nil,
		),
	}

	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/synthesize", handler.Synthesize)

	reqBody := `{"text": "测试文本"}`
	req := httptest.NewRequest("POST", "/synthesize", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTTSHandler_BatchSynthesize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{
		synthesizeResult: []byte("fake audio data"),
	}

	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/batch", handler.BatchSynthesize)

	reqBody := `{"texts": [{"text": "测试文本1"}, {"text": "测试文本2"}]}`
	req := httptest.NewRequest("POST", "/batch", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTTSHandler_BatchSynthesizeInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{}
	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.POST("/batch", handler.BatchSynthesize)

	req := httptest.NewRequest("POST", "/batch", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTTSHandler_GetSpeakers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{}
	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.GET("/speakers", handler.GetSpeakers)

	req := httptest.NewRequest("GET", "/speakers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTTSHandler_GetConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{}
	cfg := &config.TTSConfig{
		Audio: config.AudioConfig{
			SampleRate: 24000,
		},
		TTS: config.TTSModelConfig{
			Provider: config.ProviderConfig{
				Provider:   "cpu",
				DeviceID:   0,
				NumThreads: 4,
			},
		},
	}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.GET("/config", handler.GetConfig)

	req := httptest.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTTSHandler_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	manager := &mockTTSManager{
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

	cfg := &config.TTSConfig{}

	handler := NewTTSHandler(manager, cfg)

	router := gin.New()
	router.GET("/stats", handler.GetStats)

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

