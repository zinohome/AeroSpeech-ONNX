package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitorResponse 监控数据响应
type MonitorResponse struct {
	Timestamp  string                 `json:"timestamp"`
	Metrics    *MetricsData           `json:"metrics"`
	Resources  *ResourceData          `json:"resources"`
	Performance *PerformanceData      `json:"performance"`
}

// MetricsData 指标数据
type MetricsData struct {
	ActiveConnections int     `json:"active_connections"`
	RequestsPerSecond float64 `json:"requests_per_second"`
	AvgLatencyMs      float64 `json:"avg_latency_ms"`
	P95LatencyMs      float64 `json:"p95_latency_ms"`
	ErrorRate         float64 `json:"error_rate"`
}

// ResourceData 资源数据
type ResourceData struct {
	CPUUsagePercent  float64 `json:"cpu_usage_percent"`
	MemoryUsageMB    int64   `json:"memory_usage_mb"`
	ThreadCount      int     `json:"thread_count"`
	GoroutineCount   int     `json:"goroutine_count"`
	GPUUsagePercent  float64 `json:"gpu_usage_percent,omitempty"`
	GPUMemoryMB      int64   `json:"gpu_memory_mb,omitempty"`
}

// PerformanceData 性能数据
type PerformanceData struct {
	ASRPoolUsagePercent float64 `json:"asr_pool_usage_percent,omitempty"`
	TTSPoolUsagePercent float64 `json:"tts_pool_usage_percent,omitempty"`
	QueueLength         int     `json:"queue_length"`
}

// MonitorGetter 监控数据获取器接口
type MonitorGetter interface {
	GetMetrics() *MetricsData
	GetResources() *ResourceData
	GetPerformance() *PerformanceData
}

// MonitorHandler 监控数据处理器
// @Summary      获取实时监控数据
// @Description  获取服务的实时监控数据，包括指标、资源和性能数据
// @Tags         系统
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "监控数据"
// @Router       /monitor [get]
func MonitorHandler(getter MonitorGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := MonitorResponse{
			Timestamp: time.Now().Format(time.RFC3339),
		}

		if getter != nil {
			response.Metrics = getter.GetMetrics()
			response.Resources = getter.GetResources()
			response.Performance = getter.GetPerformance()
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    response,
		})
	}
}

// MonitorDataGetter 监控数据获取器接口（与MonitorGetter相同，用于兼容）
type MonitorDataGetter interface {
	GetMetrics() *MetricsData
	GetResources() *ResourceData
	GetPerformance() *PerformanceData
	GetSessions() interface{}
	GetSession(sessionID string) interface{}
}

// GetSessionsHandler 获取会话列表
func GetSessionsHandler(getter MonitorDataGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessions interface{}
		if getter != nil {
			sessions = getter.GetSessions()
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    sessions,
		})
	}
}

// GetSessionHandler 获取会话详情
func GetSessionHandler(getter MonitorDataGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("session_id")
		
		var session interface{}
		if getter != nil {
			session = getter.GetSession(sessionID)
		}

		if session == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "session not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    session,
		})
	}
}

// GetHistoryHandler 获取历史统计数据
func GetHistoryHandler(getter MonitorDataGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简化实现：返回空历史数据
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    []interface{}{},
		})
	}
}

