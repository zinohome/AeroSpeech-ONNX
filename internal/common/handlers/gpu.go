package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GPUInfoResponse GPU信息响应
type GPUInfoResponse struct {
	GPUAvailable bool      `json:"gpu_available"`
	GPUCount     int       `json:"gpu_count"`
	GPUs         []GPUInfo `json:"gpus,omitempty"`
}

// GPUInfo GPU设备信息
type GPUInfo struct {
	DeviceID      int    `json:"device_id"`
	Name          string `json:"name"`
	MemoryTotalMB int64  `json:"memory_total_mb"`
	MemoryFreeMB int64  `json:"memory_free_mb"`
	MemoryUsedMB  int64  `json:"memory_used_mb"`
	CUDAVersion   string `json:"cuda_version,omitempty"`
	DriverVersion string `json:"driver_version,omitempty"`
}

// GPUInfoGetter GPU信息获取器接口
type GPUInfoGetter interface {
	GetGPUInfo() *GPUInfoResponse
}

// GPUInfoHandler GPU信息处理器
func GPUInfoHandler(getter GPUInfoGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response *GPUInfoResponse
		if getter != nil {
			response = getter.GetGPUInfo()
		} else {
			response = &GPUInfoResponse{
				GPUAvailable: false,
				GPUCount:     0,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    response,
		})
	}
}

// GPUBenchmarkRequest GPU性能测试请求
type GPUBenchmarkRequest struct {
	Provider   string `json:"provider"`
	DeviceID   int    `json:"device_id"`
	Iterations int    `json:"iterations"`
}

// GPUBenchmarkResponse GPU性能测试响应
type GPUBenchmarkResponse struct {
	Provider      string  `json:"provider"`
	DeviceID      int     `json:"device_id"`
	Iterations    int     `json:"iterations"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	P95LatencyMs  float64 `json:"p95_latency_ms"`
	P99LatencyMs  float64 `json:"p99_latency_ms"`
	ThroughputQPS float64 `json:"throughput_qps"`
}

// GPUBenchmarkRunner GPU性能测试运行器接口
type GPUBenchmarkRunner interface {
	RunBenchmark(provider string, deviceID int, iterations int) (*GPUBenchmarkResponse, error)
}

// GPUBenchmarkHandler GPU性能测试处理器
func GPUBenchmarkHandler(runner GPUBenchmarkRunner) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GPUBenchmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "invalid request",
				"error": gin.H{
					"type":    "INVALID_PARAMS",
					"details": err.Error(),
				},
			})
			return
		}

		if runner == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    503,
				"message": "GPU benchmark not available",
			})
			return
		}

		result, err := runner.RunBenchmark(req.Provider, req.DeviceID, req.Iterations)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "benchmark failed",
				"error": gin.H{
					"type":    "BENCHMARK_ERROR",
					"details": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    result,
		})
	}
}

