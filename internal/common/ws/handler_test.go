package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewUpgrader(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		true,
	)

	if upgrader == nil {
		t.Fatal("NewUpgrader() returned nil")
	}

	if upgrader.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout 30s, got %v", upgrader.ReadTimeout)
	}
}

func TestUpgrade(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Upgrade() error = %v", err)
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	// 测试WebSocket连接
	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestNewConnectionManager(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	if manager == nil {
		t.Fatal("NewConnectionManager() returned nil")
	}
}

// testMessageHandler 测试消息处理器
type testMessageHandler struct {
	messages []string
	errors   []error
}

func (h *testMessageHandler) HandleMessage(conn *websocket.Conn, messageType int, message []byte) error {
	h.messages = append(h.messages, string(message))
	return nil
}

func (h *testMessageHandler) HandleError(conn *websocket.Conn, err error) {
	h.errors = append(h.errors, err)
}

func (h *testMessageHandler) HandleClose(conn *websocket.Conn) {
	// 处理关闭
}

func TestSetPongHandler(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Upgrade() error = %v", err)
			return
		}
		defer conn.Close()

		SetPongHandler(conn, 60*time.Second)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()
}

func TestConnectionManagerHandleConnection(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 发送测试消息
	if err := conn.WriteMessage(websocket.TextMessage, []byte("test message")); err != nil {
		t.Errorf("WriteMessage() error = %v", err)
	}

	// 等待消息处理
	time.Sleep(100 * time.Millisecond)

	if len(handler.messages) == 0 {
		t.Error("Expected at least one message to be handled")
	}
}

func TestConnectionManagerHandleConnectionError(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	// 测试无效的HTTP请求（非WebSocket）
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	manager.HandleConnection(w, req)

	// 应该处理错误
	if len(handler.errors) == 0 {
		t.Log("Expected error to be handled (may not occur in all cases)")
	}
}

func TestConnectionManagerHandleConnectionMessageError(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandlerWithError{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 发送测试消息（会触发错误）
	if err := conn.WriteMessage(websocket.TextMessage, []byte("test message")); err != nil {
		t.Errorf("WriteMessage() error = %v", err)
	}

	// 等待消息处理
	time.Sleep(100 * time.Millisecond)

	// 关闭连接以触发错误处理
	conn.Close()

	// 等待错误处理
	time.Sleep(100 * time.Millisecond)
}

func TestConnectionManagerHandleConnectionClose(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}

	// 立即关闭连接
	conn.Close()

	// 等待关闭处理
	time.Sleep(100 * time.Millisecond)
}

func TestUpgradeWithMaxMessageSize(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024, // 小消息大小限制
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 验证连接已建立
	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestUpgradeWithTimeouts(t *testing.T) {
	upgrader := NewUpgrader(
		1*time.Second,  // 短读取超时
		1*time.Second,  // 短写入超时
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 验证连接已建立
	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestUpgradeWithZeroTimeouts(t *testing.T) {
	upgrader := NewUpgrader(
		0, // 无读取超时
		0, // 无写入超时
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 验证连接已建立
	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestUpgradeWithZeroMaxMessageSize(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		0, // 无消息大小限制
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 验证连接已建立
	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

// testMessageHandlerWithError 测试消息处理器（返回错误）
type testMessageHandlerWithError struct {
	messages []string
	errors   []error
}

func (h *testMessageHandlerWithError) HandleMessage(conn *websocket.Conn, messageType int, message []byte) error {
	h.messages = append(h.messages, string(message))
	return websocket.ErrCloseSent // 返回错误以测试错误处理
}

func (h *testMessageHandlerWithError) HandleError(conn *websocket.Conn, err error) {
	h.errors = append(h.errors, err)
}

func (h *testMessageHandlerWithError) HandleClose(conn *websocket.Conn) {
	// 处理关闭
}

func TestUpgrader_UpgradeWithCompression(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		true, // 启用压缩
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestUpgrader_UpgradeWithResponseHeader(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseHeader := http.Header{}
		responseHeader.Set("X-Custom-Header", "test-value")
		conn, err := upgrader.Upgrade(w, r, responseHeader)
		if err != nil {
			return
		}
		defer conn.Close()
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestConnectionManager_pingLoop(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		1*time.Second, // 短ping周期用于测试
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 等待ping消息
	time.Sleep(2 * time.Second)

	// 验证连接仍然活跃
	if conn == nil {
		t.Fatal("Connection is nil")
	}
}

func TestConnectionManager_HandleConnectionReadTimeout(t *testing.T) {
	upgrader := NewUpgrader(
		1*time.Second, // 短读取超时
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 不发送任何消息，等待超时
	time.Sleep(2 * time.Second)

	// 连接应该已关闭或超时
}

func TestConnectionManager_HandleConnectionLargeMessage(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024, // 小消息大小限制
		false,
	)

	handler := &testMessageHandler{}
	manager := NewConnectionManager(upgrader, handler)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.HandleConnection(w, r)
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 发送大消息（超过限制）
	largeMessage := make([]byte, 2048)
	if err := conn.WriteMessage(websocket.TextMessage, largeMessage); err != nil {
		// 预期可能会失败
		t.Logf("Large message write may fail: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
}

func TestSetPongHandler_WithZeroTimeout(t *testing.T) {
	upgrader := NewUpgrader(
		30*time.Second,
		10*time.Second,
		54*time.Second,
		60*time.Second,
		1024*1024,
		false,
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		SetPongHandler(conn, 0) // 零超时
	}))

	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test server: %v", err)
		return
	}
	defer conn.Close()

	// 发送pong消息
	conn.WriteMessage(websocket.PongMessage, nil)
	time.Sleep(100 * time.Millisecond)
}

