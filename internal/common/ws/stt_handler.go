package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
)

// STTMessage STT消息结构
type STTMessage struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// STTHandler STT WebSocket处理器
type STTHandler struct {
	sessionManager *session.Manager
	asrManager     ASRManager
	config         *config.STTConfig
}

// ASRManager ASR管理器接口
type ASRManager interface {
	Transcribe(ctx interface{}, audio []byte) (string, error)
	GetStats() interface{}
	GetAvgLatency() interface{}
	GetPoolUsage() float64
}

// NewSTTHandler 创建STT处理器
func NewSTTHandler(sessionManager *session.Manager, asrManager ASRManager, cfg *config.STTConfig) *STTHandler {
	return &STTHandler{
		sessionManager: sessionManager,
		asrManager:     asrManager,
		config:         cfg,
	}
}

// HandleConnection 处理WebSocket连接
func (h *STTHandler) HandleConnection(conn *websocket.Conn) {
	// 创建会话
	sess, err := h.sessionManager.CreateSession(conn, h.config.Session.SendQueueSize)
	if err != nil {
		logger.Errorf("Failed to create session: %v", err)
		conn.Close()
		return
	}
	defer func() {
		h.sessionManager.RemoveSession(sess.ID)
		sess.Close()
	}()

	// 发送连接确认消息
	configMsg := STTMessage{
		Type:      "connection",
		SessionID: sess.ID,
		Data: map[string]interface{}{
			"status":     "connected",
			"session_id": sess.ID,
			"config": map[string]interface{}{
				"sample_rate":   h.config.Audio.SampleRate,
				"chunk_size":    h.config.Audio.ChunkSize,
				"format":        "pcm_s16le",
				"provider":      h.config.ASR.Provider.Provider,
				"gpu_available": h.config.ASR.Provider.Provider == "cuda",
				"gpu_device_id": h.config.ASR.Provider.DeviceID,
			},
		},
	}

	if err := sess.Send(configMsg); err != nil {
		logger.Errorf("Failed to send connection message: %v", err)
		return
	}

	// 设置Pong处理器
	SetPongHandler(conn, time.Duration(h.config.WebSocket.ReadTimeout)*time.Second)

	// 创建音频处理通道
	// 使用带缓冲的通道，避免处理慢阻塞读取
	audioChan := make(chan []byte, 10)

	// 启动处理协程
	go h.processPump(sess, audioChan)

	// 启动读取循环
	h.readPump(sess, conn, audioChan)
}

// readPump 读取WebSocket消息
func (h *STTHandler) readPump(sess *session.Session, conn *websocket.Conn, audioChan chan []byte) {
	defer close(audioChan)

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

		switch messageType {
		case websocket.BinaryMessage:
			// 音频数据，发送到处理通道
			select {
			case audioChan <- message:
			default:
				// 通道已满，丢弃数据或报错
				logger.Warnf("Audio buffer full for session %s, dropping packet", sess.ID)
				sess.Send(STTMessage{
					Type:      "warning",
					SessionID: sess.ID,
					Error:     "server busy, audio packet dropped",
				})
			}

		case websocket.TextMessage:
			// 文本消息（控制消息）
			var msg STTMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				logger.Warnf("Failed to parse message: %v", err)
				continue
			}

			switch msg.Type {
			case "reset":
				// 重置识别 - 这里其实很难真正重置正在处理的音频，
				// 但可以清空通道中尚未处理的数据（如果channel支持清空的话，但channel不支持直接清空）
				// 简单的做法是发送一个重置信号到processPump，或者由客户端控制
				sess.Send(STTMessage{
					Type:      "reset",
					SessionID: sess.ID,
					Data:      map[string]string{"status": "ok"},
				})

			case "ping":
				// 心跳响应
				sess.Send(STTMessage{
					Type:      "pong",
					SessionID: sess.ID,
				})
			}
		}
	}
}

// processPump 处理音频数据
func (h *STTHandler) processPump(sess *session.Session, audioChan <-chan []byte) {
	audioBuffer := make([]byte, 0, h.config.Audio.ChunkSize*2)

	for audioData := range audioChan {
		audioBuffer = append(audioBuffer, audioData...)

		// 当缓冲区达到一定大小时，进行识别
		if len(audioBuffer) >= h.config.Audio.ChunkSize {
			h.processAudio(sess, audioBuffer)
			audioBuffer = audioBuffer[:0] // 清空缓冲区
		}
	}

	// 处理剩余的音频数据
	if len(audioBuffer) > 0 {
		h.processAudio(sess, audioBuffer)
	}
}

// processAudio 执行ASR识别
func (h *STTHandler) processAudio(sess *session.Session, audio []byte) {
	// 执行识别
	result, err := h.asrManager.Transcribe(nil, audio)
	if err != nil {
		logger.Errorf("ASR transcription failed: %v", err)
		sess.Send(STTMessage{
			Type:      "error",
			SessionID: sess.ID,
			Error:     err.Error(),
		})
		return
	}

	// 发送识别结果
	sess.Send(STTMessage{
		Type:      "result",
		SessionID: sess.ID,
		Data: map[string]interface{}{
			"text":      result,
			"timestamp": time.Now().Unix(),
		},
	})
}
