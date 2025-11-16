package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/pkg/utils"
)

// STTManager STT管理器接口
type STTManager interface {
	Transcribe(ctx interface{}, audio []byte) (string, error)
	GetStats() interface{}
	GetAvgLatency() interface{}
	GetPoolUsage() float64
	GetPoolStats() map[string]interface{}
}

// STTHandler STT API处理器
type STTHandler struct {
	manager STTManager
	config  *config.STTConfig
}

// NewSTTHandler 创建STT处理器
func NewSTTHandler(manager STTManager, cfg *config.STTConfig) *STTHandler {
	return &STTHandler{
		manager: manager,
		config:  cfg,
	}
}

// RecognizeRequest 识别请求
type RecognizeRequest struct {
	Audio []byte `json:"audio,omitempty"`
	URL   string `json:"url,omitempty"`
}

// RecognizeResponse 识别响应
type RecognizeResponse struct {
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

// Recognize 文件上传识别
// @Summary      文件上传识别
// @Description  上传音频文件进行语音识别
// @Tags         STT
// @Accept       multipart/form-data
// @Produce      json
// @Param        audio  formData  file  true  "音频文件"
// @Success      200    {object}  map[string]interface{}  "识别成功"
// @Failure      400    {object}  map[string]interface{}  "请求参数错误"
// @Failure      500    {object}  map[string]interface{}  "服务器错误"
// @Router       /stt/recognize [post]
func (h *STTHandler) Recognize(c *gin.Context) {
	// 从multipart form获取文件
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request",
			"error": gin.H{
				"type":    "INVALID_PARAMS",
				"details": "audio file is required",
			},
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to open file",
			"error": gin.H{
				"type":    "FILE_ERROR",
				"details": err.Error(),
			},
		})
		return
	}
	defer src.Close()

	// 读取音频数据
	audioData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to read file",
			"error": gin.H{
				"type":    "FILE_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// 执行识别
	result, err := h.manager.Transcribe(nil, audioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "recognition failed",
			"error": gin.H{
				"type":    "RECOGNITION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": RecognizeResponse{
			Text:      result,
			Timestamp: time.Now().Unix(),
		},
	})
}

// BatchRecognizeRequest 批量识别请求
type BatchRecognizeRequest struct {
	Files []string `json:"files" binding:"required"`
}

// BatchRecognizeResponse 批量识别响应
type BatchRecognizeResponse struct {
	Results []RecognizeResponse `json:"results"`
}

// BatchRecognize 批量识别
// @Summary      批量识别
// @Description  批量识别多个音频文件
// @Tags         STT
// @Accept       json
// @Produce      json
// @Param        request  body      BatchRecognizeRequest  true  "批量识别请求"
// @Success      200      {object}  map[string]interface{}  "识别成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Failure      500      {object}  map[string]interface{}  "服务器错误"
// @Router       /stt/batch [post]
func (h *STTHandler) BatchRecognize(c *gin.Context) {
	var req BatchRecognizeRequest
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

	// 批量识别（简化实现）
	results := make([]RecognizeResponse, 0, len(req.Files))
	for _, filePath := range req.Files {
		// 读取文件
		audioData, err := utils.ReadFile(filePath)
		if err != nil {
			results = append(results, RecognizeResponse{
				Text:      "",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		// 执行识别
		result, err := h.manager.Transcribe(nil, audioData)
		if err != nil {
			results = append(results, RecognizeResponse{
				Text:      "",
				Timestamp: time.Now().Unix(),
			})
			continue
		}

		results = append(results, RecognizeResponse{
			Text:      result,
			Timestamp: time.Now().Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": BatchRecognizeResponse{
			Results: results,
		},
	})
}

// GetConfig 获取配置
// @Summary      获取STT配置
// @Description  获取语音识别服务的配置信息
// @Tags         STT
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "配置信息"
// @Router       /stt/config [get]
func (h *STTHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"sample_rate":     h.config.Audio.SampleRate,
			"chunk_size":      h.config.Audio.ChunkSize,
			"format":          "pcm_s16le",
			"provider":        h.config.ASR.Provider.Provider,
			"gpu_available":  h.config.ASR.Provider.Provider == "cuda",
			"gpu_device_id":   h.config.ASR.Provider.DeviceID,
			"num_threads":     h.config.ASR.Provider.NumThreads,
		},
	})
}

// GetStats 获取统计信息
// @Summary      获取STT统计信息
// @Description  获取语音识别服务的统计信息
// @Tags         STT
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "统计信息"
// @Router       /stt/stats [get]
func (h *STTHandler) GetStats(c *gin.Context) {
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

