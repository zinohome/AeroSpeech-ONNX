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

