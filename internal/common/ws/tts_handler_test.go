package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/session"
)

// mockTTSManager 模拟TTS管理器
type mockTTSManager struct {
	synthesizeResult []byte
	synthesizeError  error
	stats            interface{}
	avgLatency       interface{}
	poolUsage        float64
}

func (m *mockTTSManager) Synthesize(ctx interface{}, text string, speakerID int, speed float32) ([]byte, error) {
	if m.synthesizeError != nil {
		return nil, m.synthesizeError
	}
	return m.synthesizeResult, nil
}

func (m *mockTTSManager) GetStats() interface{} {
	return m.stats
}

func (m *mockTTSManager) GetAvgLatency() interface{} {
	return m.avgLatency
}

func (m *mockTTSManager) GetPoolUsage() float64 {
	return m.poolUsage
}

func TestTTSHandler_HandleConnection(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{
			synthesizeResult: []byte("fake audio data"),
		}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
					DeviceID: 0,
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		handler.HandleConnection(conn)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 等待连接确认消息
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var msg TTSMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if msg.Type != "connection" {
		t.Errorf("Expected message type 'connection', got '%s'", msg.Type)
	}
}

func TestTTSHandler_HandleSynthesize(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{
			synthesizeResult: []byte("fake audio data"),
		}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		go handler.HandleConnection(conn)
		
		time.Sleep(100 * time.Millisecond)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 读取连接确认消息
	conn.ReadMessage()

	// 发送合成请求
	synthesizeMsg := TTSMessage{
		Type: "synthesize",
		Data: map[string]interface{}{
			"text":       "测试文本",
			"speaker_id": 0,
			"speed":      1.0,
		},
	}
	msgData, _ := json.Marshal(synthesizeMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待音频数据
	time.Sleep(200 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestTTSHandler_HandlePing(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		go handler.HandleConnection(conn)
		
		time.Sleep(100 * time.Millisecond)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 读取连接确认消息
	conn.ReadMessage()

	// 发送ping消息
	pingMsg := TTSMessage{
		Type: "ping",
	}
	msgData, _ := json.Marshal(pingMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待pong响应
	time.Sleep(100 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestTTSHandler_HandleError(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{
			synthesizeError: &mockError{msg: "synthesis failed"},
		}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		go handler.HandleConnection(conn)
		
		time.Sleep(100 * time.Millisecond)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 读取连接确认消息
	conn.ReadMessage()

	// 发送合成请求（会触发错误）
	synthesizeMsg := TTSMessage{
		Type: "synthesize",
		Data: map[string]interface{}{
			"text": "测试文本",
		},
	}
	msgData, _ := json.Marshal(synthesizeMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待错误处理
	time.Sleep(200 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestTTSHandler_HandleInvalidData(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		go handler.HandleConnection(conn)
		
		time.Sleep(100 * time.Millisecond)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 读取连接确认消息
	conn.ReadMessage()

	// 发送无效数据
	synthesizeMsg := TTSMessage{
		Type: "synthesize",
		Data: "invalid data",
	}
	msgData, _ := json.Marshal(synthesizeMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待错误处理
	time.Sleep(200 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestTTSHandler_HandleEmptyText(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		sessionManager := session.NewManager(100, 30*time.Second)
		ttsManager := &mockTTSManager{}

		cfg := &config.TTSConfig{
			Audio: config.AudioConfig{
				SampleRate: 24000,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			TTS: config.TTSModelConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewTTSHandler(sessionManager, ttsManager, cfg)
		go handler.HandleConnection(conn)
		
		time.Sleep(100 * time.Millisecond)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 读取连接确认消息
	conn.ReadMessage()

	// 发送空文本
	synthesizeMsg := TTSMessage{
		Type: "synthesize",
		Data: map[string]interface{}{
			"text": "",
		},
	}
	msgData, _ := json.Marshal(synthesizeMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待错误处理
	time.Sleep(200 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

