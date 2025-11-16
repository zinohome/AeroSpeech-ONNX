package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(true, 100, 200, 1000)

	if limiter == nil {
		t.Fatal("Expected limiter to be created")
	}

	if limiter.enabled != true {
		t.Error("Expected enabled to be true")
	}

	if limiter.maxConns != 1000 {
		t.Errorf("Expected maxConns to be 1000, got %d", limiter.maxConns)
	}
}

func TestRateLimiterDisabled(t *testing.T) {
	limiter := NewRateLimiter(false, 100, 200, 1000)

	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRateLimiterGetStats(t *testing.T) {
	limiter := NewRateLimiter(true, 100, 200, 1000)

	stats := limiter.GetStats()

	if stats["enabled"] != true {
		t.Error("Expected enabled to be true in stats")
	}

	if stats["max_connections"] != 1000 {
		t.Error("Expected max_connections to be 1000 in stats")
	}

	if stats["requests_per_second"] != float64(100) {
		t.Error("Expected requests_per_second to be 100 in stats")
	}
}

