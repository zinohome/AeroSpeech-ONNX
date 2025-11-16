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

	// 验证文本列表不为空
	if len(req.Texts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request",
			"error": gin.H{
				"type":    "INVALID_PARAMS",
				"details": "texts list cannot be empty",
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
	// Kokoro v1.1 支持 103 个说话人
	// 参考: https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html
	speakers := []map[string]interface{}{}
	
	// 美式英语女声 (af: American female) - ID 0-1
	speakers = append(speakers,
		map[string]interface{}{"id": 0, "name": "af_maple", "gender": "female", "language": "en-US", "category": "American Female"},
		map[string]interface{}{"id": 1, "name": "af_sol", "gender": "female", "language": "en-US", "category": "American Female"},
	)
	
	// 英式英语女声 (bf: British female) - ID 2
	speakers = append(speakers,
		map[string]interface{}{"id": 2, "name": "bf_vale", "gender": "female", "language": "en-GB", "category": "British Female"},
	)
	
	// 中文女声 (zf: Chinese female) - ID 3-57 (55个)
	chineseFemaleIDs := []int{
		3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22,
		23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	}
	chineseFemaleNames := []string{
		"zf_001", "zf_002", "zf_003", "zf_004", "zf_005", "zf_006", "zf_007", "zf_008",
		"zf_017", "zf_018", "zf_019", "zf_021", "zf_022", "zf_023", "zf_024", "zf_026",
		"zf_027", "zf_028", "zf_032", "zf_036", "zf_038", "zf_039", "zf_040", "zf_042",
		"zf_043", "zf_044", "zf_046", "zf_047", "zf_048", "zf_049", "zf_051", "zf_059",
		"zf_060", "zf_067", "zf_070", "zf_071", "zf_072", "zf_073", "zf_074", "zf_075",
		"zf_076", "zf_077", "zf_078", "zf_079", "zf_083", "zf_084", "zf_085", "zf_086",
		"zf_087", "zf_088", "zf_090", "zf_092", "zf_093", "zf_094", "zf_099",
	}
	for i, id := range chineseFemaleIDs {
		speakers = append(speakers, map[string]interface{}{
			"id": id, "name": chineseFemaleNames[i], "gender": "female", "language": "zh", "category": "Chinese Female",
		})
	}
	
	// 中文男声 (zm: Chinese male) - ID 58-102 (45个)
	chineseMaleIDs := []int{
		58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75,
		76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93,
		94, 95, 96, 97, 98, 99, 100, 101, 102,
	}
	chineseMaleNames := []string{
		"zm_009", "zm_010", "zm_011", "zm_012", "zm_013", "zm_014", "zm_015", "zm_016",
		"zm_020", "zm_025", "zm_029", "zm_030", "zm_031", "zm_033", "zm_034", "zm_035",
		"zm_037", "zm_041", "zm_045", "zm_050", "zm_052", "zm_053", "zm_054", "zm_055",
		"zm_056", "zm_057", "zm_058", "zm_061", "zm_062", "zm_063", "zm_064", "zm_065",
		"zm_066", "zm_068", "zm_069", "zm_080", "zm_081", "zm_082", "zm_089", "zm_091",
		"zm_095", "zm_096", "zm_097", "zm_098", "zm_100",
	}
	for i, id := range chineseMaleIDs {
		speakers = append(speakers, map[string]interface{}{
			"id": id, "name": chineseMaleNames[i], "gender": "male", "language": "zh", "category": "Chinese Male",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"speakers": speakers,
			"total":    len(speakers),
			"info": gin.H{
				"model":        "Kokoro v1.1",
				"sample_rate":  24000,
				"reference":    "https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html",
			},
		},
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

