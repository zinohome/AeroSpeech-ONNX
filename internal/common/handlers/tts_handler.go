package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
)

// TTSManager TTS管理器接口
type TTSManager interface {
	Synthesize(ctx interface{}, text string, speakerID int, speed float32) ([]byte, error)
	GetStats() interface{}
	GetAvgLatency() interface{}
	GetPoolUsage() float64
	GetPoolStats() map[string]interface{}
}

// TTSHandler TTS API处理器
type TTSHandler struct {
	manager TTSManager
	config  *config.TTSConfig
}

// NewTTSHandler 创建TTS处理器
func NewTTSHandler(manager TTSManager, cfg *config.TTSConfig) *TTSHandler {
	return &TTSHandler{
		manager: manager,
		config:  cfg,
	}
}

// SynthesizeRequest 合成请求
type SynthesizeRequest struct {
	Text      string  `json:"text" binding:"required"`
	SpeakerID int     `json:"speaker_id,omitempty"`
	Speed    float32 `json:"speed,omitempty"`
}

// Synthesize 文本合成
func (h *TTSHandler) Synthesize(c *gin.Context) {
	var req SynthesizeRequest
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

	if req.Speed == 0 {
		req.Speed = 1.0
	}

	// 执行合成
	audio, err := h.manager.Synthesize(nil, req.Text, req.SpeakerID, req.Speed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "synthesis failed",
			"error": gin.H{
				"type":    "SYNTHESIS_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// 返回音频数据
	c.Data(http.StatusOK, "audio/wav", audio)
}

// BatchSynthesizeRequest 批量合成请求
type BatchSynthesizeRequest struct {
	Texts []SynthesizeRequest `json:"texts" binding:"required"`
}

// BatchSynthesize 批量合成
func (h *TTSHandler) BatchSynthesize(c *gin.Context) {
	var req BatchSynthesizeRequest
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

	// 批量合成（简化实现）
	results := make([]map[string]interface{}, 0, len(req.Texts))
	for _, textReq := range req.Texts {
		if textReq.Speed == 0 {
			textReq.Speed = 1.0
		}

		audio, err := h.manager.Synthesize(nil, textReq.Text, textReq.SpeakerID, textReq.Speed)
		if err != nil {
			results = append(results, map[string]interface{}{
				"text":  textReq.Text,
				"error": err.Error(),
			})
			continue
		}

		results = append(results, map[string]interface{}{
			"text":      textReq.Text,
			"audio":     audio,
			"timestamp": time.Now().Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    results,
	})
}

// GetSpeakers 获取说话人列表
func (h *TTSHandler) GetSpeakers(c *gin.Context) {
	// 返回默认说话人列表（简化实现）
	speakers := []map[string]interface{}{
		{"id": 0, "name": "默认"},
		{"id": 45, "name": "小贝"},
		{"id": 46, "name": "小妮"},
		{"id": 47, "name": "小小"},
		{"id": 48, "name": "小艺"},
		{"id": 49, "name": "云健"},
		{"id": 50, "name": "云溪"},
		{"id": 51, "name": "云霞"},
		{"id": 52, "name": "云阳"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    speakers,
	})
}

// GetConfig 获取配置
func (h *TTSHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"sample_rate":    h.config.Audio.SampleRate,
			"format":         "pcm_s16le",
			"provider":       h.config.TTS.Provider.Provider,
			"gpu_available":  h.config.TTS.Provider.Provider == "cuda",
			"gpu_device_id":  h.config.TTS.Provider.DeviceID,
			"num_threads":    h.config.TTS.Provider.NumThreads,
		},
	})
}

// GetStats 获取统计信息
func (h *TTSHandler) GetStats(c *gin.Context) {
	stats := h.manager.GetStats()
	poolStats := h.manager.GetPoolStats()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"stats":      stats,
			"pool_stats": poolStats,
			"pool_usage": h.manager.GetPoolUsage(),
		},
	})
}

