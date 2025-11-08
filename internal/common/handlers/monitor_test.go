package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// mockMonitorGetter 模拟监控数据获取器
type mockMonitorGetter struct {
	metrics     *MetricsData
	resources   *ResourceData
	performance *PerformanceData
}

func (m *mockMonitorGetter) GetMetrics() *MetricsData {
	return m.metrics
}

func (m *mockMonitorGetter) GetResources() *ResourceData {
	return m.resources
}

func (m *mockMonitorGetter) GetPerformance() *PerformanceData {
	return m.performance
}

// mockMonitorDataGetter 模拟监控数据获取器（包含会话方法）
type mockMonitorDataGetter struct {
	*mockMonitorGetter
	sessions map[string]interface{}
}

func (m *mockMonitorDataGetter) GetSessions() interface{} {
	return m.sessions
}

func (m *mockMonitorDataGetter) GetSession(sessionID string) interface{} {
	if session, ok := m.sessions[sessionID]; ok {
		return session
	}
	return nil
}

func TestMonitorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	getter := &mockMonitorGetter{
		metrics: &MetricsData{
			ActiveConnections: 10,
			RequestsPerSecond: 5.5,
			AvgLatencyMs:      50.0,
		},
		resources: &ResourceData{
			CPUUsagePercent: 25.5,
			MemoryUsageMB:   512,
		},
		performance: &PerformanceData{
			ASRPoolUsagePercent: 60.0,
			QueueLength:         5,
		},
	}

	handler := MonitorHandler(getter)

	router := gin.New()
	router.GET("/monitor", handler)

	req := httptest.NewRequest("GET", "/monitor", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMonitorHandlerNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := MonitorHandler(nil)

	router := gin.New()
	router.GET("/monitor", handler)

	req := httptest.NewRequest("GET", "/monitor", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSessionsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	getter := &mockMonitorDataGetter{
		mockMonitorGetter: &mockMonitorGetter{},
		sessions: map[string]interface{}{
			"session1": map[string]interface{}{
				"id":   "session1",
				"status": "active",
			},
			"session2": map[string]interface{}{
				"id":   "session2",
				"status": "idle",
			},
		},
	}

	handler := GetSessionsHandler(getter)

	router := gin.New()
	router.GET("/sessions", handler)

	req := httptest.NewRequest("GET", "/sessions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSessionsHandlerNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := GetSessionsHandler(nil)

	router := gin.New()
	router.GET("/sessions", handler)

	req := httptest.NewRequest("GET", "/sessions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSessionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	getter := &mockMonitorDataGetter{
		mockMonitorGetter: &mockMonitorGetter{},
		sessions: map[string]interface{}{
			"session1": map[string]interface{}{
				"id":   "session1",
				"status": "active",
			},
		},
	}

	handler := GetSessionHandler(getter)

	router := gin.New()
	router.GET("/sessions/:session_id", handler)

	req := httptest.NewRequest("GET", "/sessions/session1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetSessionHandlerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	getter := &mockMonitorDataGetter{
		mockMonitorGetter: &mockMonitorGetter{},
		sessions: map[string]interface{}{},
	}

	handler := GetSessionHandler(getter)

	router := gin.New()
	router.GET("/sessions/:session_id", handler)

	req := httptest.NewRequest("GET", "/sessions/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetHistoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	getter := &mockMonitorDataGetter{
		mockMonitorGetter: &mockMonitorGetter{},
	}

	handler := GetHistoryHandler(getter)

	router := gin.New()
	router.GET("/history", handler)

	req := httptest.NewRequest("GET", "/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

