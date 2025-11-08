package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	components := map[string]string{
		"asr": "ok",
		"tts": "ok",
	}
	
	provider := &ProviderInfo{
		ASR:          "cpu",
		TTS:          "cpu",
		GPUAvailable: false,
		GPUDeviceID:  0,
	}
	
	handler := HealthHandler(components, provider)
	
	router := gin.New()
	router.GET("/health", handler)
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthHandlerWithoutComponents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handler := HealthHandler(nil, nil)
	
	router := gin.New()
	router.GET("/health", handler)
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthHandlerWithGPU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	provider := &ProviderInfo{
		ASR:          "cuda",
		TTS:          "cuda",
		GPUAvailable: true,
		GPUDeviceID:  0,
	}
	
	handler := HealthHandler(nil, provider)
	
	router := gin.New()
	router.GET("/health", handler)
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

