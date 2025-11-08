package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockGPUInfoGetter 模拟GPU信息获取器
type mockGPUInfoGetter struct {
	info *GPUInfoResponse
}

func (m *mockGPUInfoGetter) GetGPUInfo() *GPUInfoResponse {
	return m.info
}

func TestGPUInfoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	getter := &mockGPUInfoGetter{
		info: &GPUInfoResponse{
			GPUAvailable: false,
			GPUCount:     0,
		},
	}
	
	handler := GPUInfoHandler(getter)
	
	router := gin.New()
	router.GET("/gpu", handler)
	
	req := httptest.NewRequest("GET", "/gpu", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGPUInfoHandlerNil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handler := GPUInfoHandler(nil)
	
	router := gin.New()
	router.GET("/gpu", handler)
	
	req := httptest.NewRequest("GET", "/gpu", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

