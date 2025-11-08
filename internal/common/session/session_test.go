package session

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(100, 30*time.Second)
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if manager.maxSessions != 100 {
		t.Errorf("Expected maxSessions 100, got %d", manager.maxSessions)
	}
}

func TestCreateSession(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	// 创建模拟会话（不实际连接）
	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if session == nil {
		t.Fatal("CreateSession() returned nil")
	}

	if session.ID == "" {
		t.Error("Session ID is empty")
	}

	if session.Status != StatusActive {
		t.Errorf("Expected status Active, got %v", session.Status)
	}
}

func TestGetSession(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := manager.GetSession(session.ID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if got.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, got.ID)
	}

	// 测试不存在的会话
	_, err = manager.GetSession("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
}

func TestRemoveSession(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	manager.RemoveSession(session.ID)

	_, err = manager.GetSession(session.ID)
	if err == nil {
		t.Error("Expected error after removing session")
	}
}

func TestUpdateActivity(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	oldActive := session.LastActive
	time.Sleep(10 * time.Millisecond)

	err = manager.UpdateActivity(session.ID)
	if err != nil {
		t.Fatalf("UpdateActivity() error = %v", err)
	}

	if !session.LastActive.After(oldActive) {
		t.Error("LastActive was not updated")
	}
}

func TestGetStats(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	// 创建多个会话
	for i := 0; i < 3; i++ {
		_, err := manager.CreateSession(nil, 100)
		if err != nil {
			t.Fatalf("CreateSession() error = %v", err)
		}
	}

	stats := manager.GetStats()
	if stats.Total != 3 {
		t.Errorf("Expected total 3, got %d", stats.Total)
	}

	if stats.Active != 3 {
		t.Errorf("Expected active 3, got %d", stats.Active)
	}
}

func TestSessionSend(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	message := map[string]interface{}{
		"type": "test",
		"data": "test message",
	}

	err = session.Send(message)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSessionClose(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	session.Close()

	if session.GetStatus() != StatusClosed {
		t.Error("Session status should be Closed")
	}

	// 关闭后发送应该失败
	err = session.Send(map[string]interface{}{"test": "data"})
	if err == nil {
		t.Error("Expected error when sending to closed session")
	}
}

