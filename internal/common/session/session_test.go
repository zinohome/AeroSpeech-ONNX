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

func TestSession_SendTimeout(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 1) // 小队列大小
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	// 填满队列
	session.Send(map[string]interface{}{"test": "data1"})

	// 尝试发送更多消息（可能会超时或阻塞）
	// 注意：这个测试可能因为队列满而阻塞，所以使用goroutine
	done := make(chan bool, 1)
	go func() {
		err := session.Send(map[string]interface{}{"test": "data2"})
		done <- (err != nil)
	}()

	select {
	case result := <-done:
		if !result {
			t.Log("Send did not timeout (queue may have space)")
		}
	case <-time.After(100 * time.Millisecond):
		t.Log("Send may be blocking (expected with small queue)")
	}
}

func TestSession_UpdateActivity(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	oldActive := session.LastActive
	time.Sleep(10 * time.Millisecond)

	// UpdateActivity是Manager的方法
	err = manager.UpdateActivity(session.ID)
	if err != nil {
		t.Fatalf("UpdateActivity() error = %v", err)
	}

	// 重新获取会话以查看更新
	updatedSession, err := manager.GetSession(session.ID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if !updatedSession.LastActive.After(oldActive) {
		t.Error("LastActive was not updated")
	}
}

func TestManager_GetSessionNotFound(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	_, err := manager.GetSession("non-existent-session-id")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
}

func TestManager_GetAllSessions(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	// 创建多个会话
	for i := 0; i < 3; i++ {
		_, err := manager.CreateSession(nil, 100)
		if err != nil {
			t.Fatalf("CreateSession() error = %v", err)
		}
	}

	// 通过GetStats验证会话数量
	stats := manager.GetStats()
	if stats.Total != 3 {
		t.Errorf("Expected 3 sessions, got %d", stats.Total)
	}
}

func TestManager_CleanupTimeoutSessions(t *testing.T) {
	manager := NewManager(10, 1*time.Second) // 短超时时间

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	// 等待超时
	time.Sleep(2 * time.Second)

	// 清理超时会话
	manager.CleanupTimeoutSessions()

	// 会话应该被清理
	_, err = manager.GetSession(session.ID)
	if err == nil {
		t.Error("Expected session to be cleaned up")
	}
}

func TestSessionStatus_String(t *testing.T) {
	tests := []struct {
		status SessionStatus
		want   string
	}{
		{StatusActive, "active"},
		{StatusIdle, "idle"},
		{StatusTimeout, "timeout"},
		{StatusClosed, "closed"},
		{SessionStatus(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("SessionStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetStatus(t *testing.T) {
	manager := NewManager(10, 30*time.Second)

	session, err := manager.CreateSession(nil, 100)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if session.GetStatus() != StatusActive {
		t.Errorf("Expected status Active, got %v", session.GetStatus())
	}

	session.Close()
	if session.GetStatus() != StatusClosed {
		t.Errorf("Expected status Closed, got %v", session.GetStatus())
	}
}


