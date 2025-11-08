package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader WebSocket升级器
type Upgrader struct {
	*websocket.Upgrader
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	PingPeriod       time.Duration
	PongWait         time.Duration
	MaxMessageSize   int64
	EnableCompression bool
}

// NewUpgrader 创建WebSocket升级器
func NewUpgrader(readTimeout, writeTimeout, pingPeriod, pongWait time.Duration, maxMessageSize int64, enableCompression bool) *Upgrader {
	return &Upgrader{
		Upgrader: &websocket.Upgrader{
			CheckOrigin:       func(r *http.Request) bool { return true },
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: enableCompression,
		},
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		PingPeriod:       pingPeriod,
		PongWait:         pongWait,
		MaxMessageSize:   maxMessageSize,
		EnableCompression: enableCompression,
	}
}

// Upgrade 升级HTTP连接为WebSocket
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	conn, err := u.Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	// 设置读取限制
	if u.MaxMessageSize > 0 {
		conn.SetReadLimit(u.MaxMessageSize)
	}

	// 设置读取超时
	if u.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(u.ReadTimeout))
	}

	// 设置写入超时
	if u.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(u.WriteTimeout))
	}

	return conn, nil
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(conn *websocket.Conn, messageType int, message []byte) error
	HandleError(conn *websocket.Conn, err error)
	HandleClose(conn *websocket.Conn)
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	upgrader *Upgrader
	handler  MessageHandler
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(upgrader *Upgrader, handler MessageHandler) *ConnectionManager {
	return &ConnectionManager{
		upgrader: upgrader,
		handler:  handler,
	}
}

// HandleConnection 处理WebSocket连接
func (cm *ConnectionManager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		cm.handler.HandleError(nil, err)
		return
	}
	defer conn.Close()

	// 启动心跳
	done := make(chan struct{})
	go cm.pingLoop(conn, done)

	// 消息循环
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				cm.handler.HandleError(conn, err)
			}
			close(done)
			cm.handler.HandleClose(conn)
			break
		}

		// 更新读取超时
		if cm.upgrader.ReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(cm.upgrader.ReadTimeout))
		}

		// 处理消息
		if err := cm.handler.HandleMessage(conn, messageType, message); err != nil {
			cm.handler.HandleError(conn, err)
			break
		}
	}
}

// pingLoop 心跳循环
func (cm *ConnectionManager) pingLoop(conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(cm.upgrader.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(cm.upgrader.WriteTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SetPongHandler 设置Pong处理器
func SetPongHandler(conn *websocket.Conn, pongWait time.Duration) {
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

