package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StatsResponse 统计信息响应
type StatsResponse struct {
	Timestamp string                 `json:"timestamp"`
	ASR       *ASRStats              `json:"asr,omitempty"`
	TTS       *TTSStats              `json:"tts,omitempty"`
	Sessions  *SessionStats          `json:"sessions,omitempty"`
	Resources *ResourceStats         `json:"resources,omitempty"`
}

// ASRStats ASR统计信息
type ASRStats struct {
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	AvgLatencyMs       float64 `json:"avg_latency_ms"`
	P95LatencyMs       float64 `json:"p95_latency_ms"`
	P99LatencyMs       float64 `json:"p99_latency_ms"`
	RequestsPerSecond  float64 `json:"requests_per_second"`
}

// TTSStats TTS统计信息
type TTSStats struct {
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	AvgLatencyMs       float64 `json:"avg_latency_ms"`
	P95LatencyMs       float64 `json:"p95_latency_ms"`
	P99LatencyMs       float64 `json:"p99_latency_ms"`
	RequestsPerSecond  float64 `json:"requests_per_second"`
}

// SessionStats 会话统计信息
type SessionStats struct {
	Active            int   `json:"active"`
	Total             int   `json:"total"`
	AvgDurationSeconds int64 `json:"avg_duration_seconds"`
}

// ResourceStats 资源统计信息
type ResourceStats struct {
	CPUUsagePercent  float64 `json:"cpu_usage_percent"`
	MemoryUsageMB    int64   `json:"memory_usage_mb"`
	ThreadCount      int     `json:"thread_count"`
	GoroutineCount   int     `json:"goroutine_count"`
	PoolUsagePercent float64 `json:"pool_usage_percent"`
	QueueLength      int     `json:"queue_length"`
}

// StatsGetter 统计信息获取器接口
type StatsGetter interface {
	GetASRStats() *ASRStats
	GetTTSStats() *TTSStats
	GetSessionStats() *SessionStats
	GetResourceStats() *ResourceStats
}

// StatsHandler 统计信息处理器
func StatsHandler(getter StatsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := StatsResponse{
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if getter != nil {
			response.ASR = getter.GetASRStats()
			response.TTS = getter.GetTTSStats()
			response.Sessions = getter.GetSessionStats()
			response.Resources = getter.GetResourceStats()
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    response,
		})
	}
}

