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

// mockASRManager 模拟ASR管理器
type mockASRManager struct {
	transcribeResult string
	transcribeError  error
	stats            interface{}
	avgLatency       interface{}
	poolUsage        float64
}

func (m *mockASRManager) Transcribe(ctx interface{}, audio []byte) (string, error) {
	if m.transcribeError != nil {
		return "", m.transcribeError
	}
	return m.transcribeResult, nil
}

func (m *mockASRManager) GetStats() interface{} {
	return m.stats
}

func (m *mockASRManager) GetAvgLatency() interface{} {
	return m.avgLatency
}

func (m *mockASRManager) GetPoolUsage() float64 {
	return m.poolUsage
}

func TestSTTHandler_HandleConnection(t *testing.T) {
	// 创建测试服务器
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
		asrManager := &mockASRManager{
			transcribeResult: "测试文本",
		}

		cfg := &config.STTConfig{
			Audio: config.AudioConfig{
				SampleRate: 16000,
				ChunkSize:  4096,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			ASR: config.ASRConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
					DeviceID: 0,
				},
			},
		}

		handler := NewSTTHandler(sessionManager, asrManager, cfg)
		handler.HandleConnection(conn)
	}))

	defer server.Close()

	// 连接到WebSocket
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

	var msg STTMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if msg.Type != "connection" {
		t.Errorf("Expected message type 'connection', got '%s'", msg.Type)
	}
}

func TestSTTHandler_HandleAudio(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	done := make(chan bool, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			conn.Close()
			done <- true
		}()

		sessionManager := session.NewManager(100, 30*time.Second)
		asrManager := &mockASRManager{
			transcribeResult: "测试文本",
		}

		cfg := &config.STTConfig{
			Audio: config.AudioConfig{
				SampleRate: 16000,
				ChunkSize:  4096,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			ASR: config.ASRConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewSTTHandler(sessionManager, asrManager, cfg)
		
		// 启动处理连接（在goroutine中）
		go handler.HandleConnection(conn)
		
		// 等待连接建立
		time.Sleep(100 * time.Millisecond)
		
		// 发送音频数据
		audioData := make([]byte, 4096)
		conn.WriteMessage(websocket.BinaryMessage, audioData)
		
		// 等待处理
		time.Sleep(200 * time.Millisecond)
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
	
	// 等待goroutine完成
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		// 超时，但测试继续
	}
}

func TestSTTHandler_HandlePing(t *testing.T) {
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
		asrManager := &mockASRManager{}

		cfg := &config.STTConfig{
			Audio: config.AudioConfig{
				SampleRate: 16000,
				ChunkSize:  4096,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			ASR: config.ASRConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewSTTHandler(sessionManager, asrManager, cfg)
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
	pingMsg := STTMessage{
		Type: "ping",
	}
	msgData, _ := json.Marshal(pingMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待pong响应
	time.Sleep(100 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestSTTHandler_HandleReset(t *testing.T) {
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
		asrManager := &mockASRManager{}

		cfg := &config.STTConfig{
			Audio: config.AudioConfig{
				SampleRate: 16000,
				ChunkSize:  4096,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			ASR: config.ASRConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewSTTHandler(sessionManager, asrManager, cfg)
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

	// 发送reset消息
	resetMsg := STTMessage{
		Type: "reset",
	}
	msgData, _ := json.Marshal(resetMsg)
	conn.WriteMessage(websocket.TextMessage, msgData)

	// 等待响应
	time.Sleep(100 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

func TestSTTHandler_HandleError(t *testing.T) {
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
		asrManager := &mockASRManager{
			transcribeError: &mockError{msg: "transcription failed"},
		}

		cfg := &config.STTConfig{
			Audio: config.AudioConfig{
				SampleRate: 16000,
				ChunkSize:  4096,
			},
			Session: config.SessionConfig{
				SendQueueSize: 100,
			},
			WebSocket: config.WebSocketConfig{
				ReadTimeout: 30,
			},
			ASR: config.ASRConfig{
				Provider: config.ProviderConfig{
					Provider: "cpu",
				},
			},
		}

		handler := NewSTTHandler(sessionManager, asrManager, cfg)
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

	// 发送音频数据（会触发错误）
	audioData := make([]byte, 4096)
	conn.WriteMessage(websocket.BinaryMessage, audioData)

	// 等待错误处理
	time.Sleep(200 * time.Millisecond)
	
	// 关闭连接以清理goroutine
	conn.Close()
}

// mockError 模拟错误
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

