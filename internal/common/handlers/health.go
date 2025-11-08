package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string                 `json:"status"`
	Version   string                 `json:"version"`
	Uptime    int64                  `json:"uptime_seconds"`
	Timestamp string                 `json:"timestamp"`
	Components map[string]string     `json:"components,omitempty"`
	Provider  *ProviderInfo          `json:"provider,omitempty"`
}

// ProviderInfo Provider信息
type ProviderInfo struct {
	ASR          string `json:"asr,omitempty"`
	TTS          string `json:"tts,omitempty"`
	GPUAvailable bool   `json:"gpu_available,omitempty"`
	GPUDeviceID  int    `json:"gpu_device_id,omitempty"`
}

var (
	startTime = time.Now()
	version   = "1.0.0"
)

// HealthHandler 健康检查处理器
func HealthHandler(components map[string]string, provider *ProviderInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Version:   version,
			Uptime:    int64(time.Since(startTime).Seconds()),
			Timestamp: time.Now().Format(time.RFC3339),
			Components: components,
			Provider:  provider,
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    response,
		})
	}
}

