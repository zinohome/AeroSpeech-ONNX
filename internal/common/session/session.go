package session

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SessionStatus 会话状态
type SessionStatus int

const (
	StatusActive SessionStatus = iota
	StatusIdle
	StatusTimeout
	StatusClosed
)

// String 返回会话状态的字符串表示
func (s SessionStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusIdle:
		return "idle"
	case StatusTimeout:
		return "timeout"
	case StatusClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// Session 会话结构
type Session struct {
	ID          string
	Conn        *websocket.Conn
	Status      SessionStatus
	CreatedAt   time.Time
	LastActive  time.Time
	SendQueue   chan interface{}
	mu          sync.RWMutex
	closeOnce   sync.Once
	closeChan   chan struct{}
}

// Manager 会话管理器
type Manager struct {
	sessions    map[string]*Session
	mu          sync.RWMutex
	maxSessions int
	timeout     time.Duration
}

// NewManager 创建会话管理器
func NewManager(maxSessions int, timeout time.Duration) *Manager {
	return &Manager{
		sessions:    make(map[string]*Session),
		maxSessions: maxSessions,
		timeout:     timeout,
	}
}

// CreateSession 创建新会话
func (m *Manager) CreateSession(conn *websocket.Conn, queueSize int) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.sessions) >= m.maxSessions {
		return nil, ErrMaxSessionsReached
	}

	sessionID := uuid.New().String()
	session := &Session{
		ID:         sessionID,
		Conn:       conn,
		Status:     StatusActive,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
		SendQueue:  make(chan interface{}, queueSize),
		closeChan:  make(chan struct{}),
	}

	m.sessions[sessionID] = session

	// 启动发送goroutine
	go session.sendLoop()

	return session, nil
}

// GetSession 获取会话
func (m *Manager) GetSession(sessionID string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// RemoveSession 移除会话
func (m *Manager) RemoveSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if ok {
		session.Close()
		delete(m.sessions, sessionID)
	}
}

// UpdateActivity 更新会话活动时间
func (m *Manager) UpdateActivity(sessionID string) error {
	m.mu.RLock()
	session, ok := m.sessions[sessionID]
	m.mu.RUnlock()

	if !ok {
		return ErrSessionNotFound
	}

	session.mu.Lock()
	session.LastActive = time.Now()
	if session.Status == StatusIdle {
		session.Status = StatusActive
	}
	session.mu.Unlock()

	return nil
}

// CleanupTimeoutSessions 清理超时会话
func (m *Manager) CleanupTimeoutSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		session.mu.RLock()
		lastActive := session.LastActive
		status := session.Status
		session.mu.RUnlock()

		if status == StatusClosed {
			delete(m.sessions, id)
			continue
		}

		if now.Sub(lastActive) > m.timeout {
			session.mu.Lock()
			if session.Status != StatusClosed {
				session.Status = StatusTimeout
			}
			session.mu.Unlock()
			session.Close()
			delete(m.sessions, id)
		}
	}
}

// GetStats 获取统计信息
func (m *Manager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := Stats{
		Total:   len(m.sessions),
		Active:  0,
		Idle:    0,
		Timeout: 0,
	}

	for _, session := range m.sessions {
		session.mu.RLock()
		status := session.Status
		session.mu.RUnlock()

		switch status {
		case StatusActive:
			stats.Active++
		case StatusIdle:
			stats.Idle++
		case StatusTimeout:
			stats.Timeout++
		}
	}

	return stats
}

// Stats 统计信息
type Stats struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	Idle    int `json:"idle"`
	Timeout int `json:"timeout"`
}

// sendLoop 发送消息循环
func (s *Session) sendLoop() {
	for {
		select {
		case message := <-s.SendQueue:
			s.mu.RLock()
			conn := s.Conn
			status := s.Status
			s.mu.RUnlock()

			if status == StatusClosed || conn == nil {
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				s.mu.Lock()
				s.Status = StatusClosed
				s.mu.Unlock()
				return
			}

		case <-s.closeChan:
			return
		}
	}
}

// Send 发送消息
func (s *Session) Send(message interface{}) error {
	s.mu.RLock()
	status := s.Status
	s.mu.RUnlock()

	if status == StatusClosed {
		return ErrSessionClosed
	}

	select {
	case s.SendQueue <- message:
		return nil
	default:
		return ErrSendQueueFull
	}
}

// Close 关闭会话
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.mu.Lock()
		s.Status = StatusClosed
		if s.Conn != nil {
			s.Conn.Close()
		}
		close(s.closeChan)
		s.mu.Unlock()
	})
}

// GetID 获取会话ID
func (s *Session) GetID() string {
	return s.ID
}

// GetStatus 获取会话状态
func (s *Session) GetStatus() SessionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status
}

// GetDuration 获取会话持续时间
func (s *Session) GetDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Since(s.CreatedAt)
}

// 错误定义
var (
	ErrSessionNotFound     = &SessionError{Message: "session not found"}
	ErrSessionClosed       = &SessionError{Message: "session is closed"}
	ErrMaxSessionsReached  = &SessionError{Message: "max sessions reached"}
	ErrSendQueueFull       = &SessionError{Message: "send queue is full"}
)

// SessionError 会话错误
type SessionError struct {
	Message string
}

func (e *SessionError) Error() string {
	return e.Message
}

