package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/handlers"
)

func TestTTSAPI(t *testing.T) {
	// 创建测试配置
	cfg := &config.TTSConfig{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: 8081,
		},
		TTS: config.TTSModelConfig{
			ModelPath: "/tmp/test-model.onnx",
			Provider: config.ProviderConfig{
				Provider:   "cpu",
				NumThreads: 1,
			},
		},
	}

	// 创建测试处理器（使用mock manager）
	handler := handlers.NewTTSHandler(nil, cfg)

	// 创建Gin引擎
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/tts/synthesize", handler.Synthesize)
	router.GET("/api/v1/tts/config", handler.GetConfig)
	router.GET("/api/v1/tts/speakers", handler.GetSpeakers)

	// 测试配置接口
	t.Run("GetConfig", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tts/config", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	// 测试说话人列表接口
	t.Run("GetSpeakers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tts/speakers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["code"].(float64) != 200 {
			t.Errorf("Expected code 200, got %v", response["code"])
		}
	})
}

