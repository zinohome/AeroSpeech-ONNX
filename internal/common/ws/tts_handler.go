package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
)

// TTSMessage TTS消息结构
type TTSMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// TTSHandler TTS WebSocket处理器
type TTSHandler struct {
	sessionManager *session.Manager
	ttsManager     TTSManager
	config         *config.TTSConfig
}

// TTSManager TTS管理器接口
type TTSManager interface {
	Synthesize(ctx interface{}, text string, speakerID int, speed float32) ([]byte, error)
	GetStats() interface{}
	GetAvgLatency() interface{}
	GetPoolUsage() float64
}

// NewTTSHandler 创建TTS处理器
func NewTTSHandler(sessionManager *session.Manager, ttsManager TTSManager, cfg *config.TTSConfig) *TTSHandler {
	return &TTSHandler{
		sessionManager: sessionManager,
		ttsManager:     ttsManager,
		config:         cfg,
	}
}

// HandleConnection 处理WebSocket连接
func (h *TTSHandler) HandleConnection(conn *websocket.Conn) {
	// 创建会话
	sess, err := h.sessionManager.CreateSession(conn, h.config.Session.SendQueueSize)
	if err != nil {
		logger.Errorf("Failed to create session: %v", err)
		conn.Close()
		return
	}

	// 发送连接确认消息
	configMsg := TTSMessage{
		Type:      "connection",
		SessionID: sess.ID,
		Data: map[string]interface{}{
			"status":      "connected",
			"session_id": sess.ID,
			"config": map[string]interface{}{
				"sample_rate":     h.config.Audio.SampleRate,
				"format":          "pcm_s16le",
				"provider":        h.config.TTS.Provider.Provider,
				"gpu_available":   h.config.TTS.Provider.Provider == "cuda",
				"gpu_device_id":   h.config.TTS.Provider.DeviceID,
			},
		},
	}

	if err := sess.Send(configMsg); err != nil {
		logger.Errorf("Failed to send connection message: %v", err)
		sess.Close()
		return
	}

	// 设置Pong处理器
	SetPongHandler(conn, time.Duration(h.config.WebSocket.ReadTimeout)*time.Second)

	// 处理消息循环
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("WebSocket error: %v", err)
			}
			break
		}

		// 更新会话活动时间
		h.sessionManager.UpdateActivity(sess.ID)

		if messageType == websocket.TextMessage {
			// 文本消息（合成请求）
			var msg TTSMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				logger.Warnf("Failed to parse message: %v", err)
				continue
			}

			switch msg.Type {
			case "synthesize":
				h.processSynthesize(sess, conn, msg)

			case "ping":
				// 心跳响应
				sess.Send(TTSMessage{
					Type:      "pong",
					SessionID: sess.ID,
				})
			}
		}
	}

	// 清理会话
	h.sessionManager.RemoveSession(sess.ID)
}

// processSynthesize 处理合成请求
func (h *TTSHandler) processSynthesize(sess *session.Session, conn *websocket.Conn, msg TTSMessage) {
	// 解析请求数据
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		sess.Send(TTSMessage{
			Type:      "error",
			SessionID: sess.ID,
			Error:     "invalid request data",
		})
		return
	}

	text, _ := data["text"].(string)
	speakerID := 0
	if sid, ok := data["speaker_id"].(float64); ok {
		speakerID = int(sid)
	}
	speed := float32(1.0)
	if s, ok := data["speed"].(float64); ok {
		speed = float32(s)
	}

	if text == "" {
		sess.Send(TTSMessage{
			Type:      "error",
			SessionID: sess.ID,
			Error:     "text is required",
		})
		return
	}

	// 执行合成
	audio, err := h.ttsManager.Synthesize(nil, text, speakerID, speed)
	if err != nil {
		logger.Errorf("TTS synthesis failed: %v", err)
		sess.Send(TTSMessage{
			Type:      "error",
			SessionID: sess.ID,
			Error:     err.Error(),
		})
		return
	}

	// 分块发送音频数据（直接使用WebSocket连接发送二进制数据）
	chunkSize := 4096
	for i := 0; i < len(audio); i += chunkSize {
		end := i + chunkSize
		if end > len(audio) {
			end = len(audio)
		}

		// 直接使用conn发送二进制音频数据
		if err := conn.WriteMessage(websocket.BinaryMessage, audio[i:end]); err != nil {
			logger.Errorf("Failed to send audio chunk: %v", err)
			return
		}
	}

	// 发送完成消息
	sess.Send(TTSMessage{
		Type:      "complete",
		SessionID: sess.ID,
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	})
}

