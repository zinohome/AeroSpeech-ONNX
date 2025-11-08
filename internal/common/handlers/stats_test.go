package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockStatsGetter 模拟统计信息获取器
type mockStatsGetter struct {
	asrStats      *ASRStats
	ttsStats      *TTSStats
	sessionStats  *SessionStats
	resourceStats *ResourceStats
}

func (m *mockStatsGetter) GetASRStats() *ASRStats {
	return m.asrStats
}

func (m *mockStatsGetter) GetTTSStats() *TTSStats {
	return m.ttsStats
}

func (m *mockStatsGetter) GetSessionStats() *SessionStats {
	return m.sessionStats
}

func (m *mockStatsGetter) GetResourceStats() *ResourceStats {
	return m.resourceStats
}

func TestStatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	getter := &mockStatsGetter{
		asrStats: &ASRStats{
			TotalRequests:      100,
			SuccessfulRequests: 95,
			FailedRequests:     5,
			AvgLatencyMs:       50.0,
		},
		sessionStats: &SessionStats{
			Active: 10,
			Total:  100,
		},
	}
	
	handler := StatsHandler(getter)
	
	router := gin.New()
	router.GET("/stats", handler)
	
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestStatsHandlerNil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handler := StatsHandler(nil)
	
	router := gin.New()
	router.GET("/stats", handler)
	
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

