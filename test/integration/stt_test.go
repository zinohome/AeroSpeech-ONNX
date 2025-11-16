package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/handlers"
)

func TestSTTAPI(t *testing.T) {
	// 创建测试配置
	cfg := &config.STTConfig{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		ASR: config.ASRConfig{
			ModelPath:  "/tmp/test-model.onnx",
			TokensPath: "/tmp/test-tokens.txt",
			Provider: config.ProviderConfig{
				Provider:   "cpu",
				NumThreads: 1,
			},
		},
	}

	// 创建测试处理器（使用mock manager）
	handler := handlers.NewSTTHandler(nil, cfg)

	// 创建Gin引擎
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/stt/recognize", handler.Recognize)
	router.GET("/api/v1/stt/config", handler.GetConfig)

	// 测试配置接口
	t.Run("GetConfig", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/stt/config", nil)
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

