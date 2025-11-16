package vad

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	config := &PoolConfig{
		PoolSize:  10,
		MaxIdle:   5,
		Threshold: 0.5,
		VADType:   "test",
	}

	pool := NewPool(config)

	if pool == nil {
		t.Fatal("Expected pool to be created")
	}

	if pool.config.PoolSize != 10 {
		t.Errorf("Expected PoolSize to be 10, got %d", pool.config.PoolSize)
	}

	if pool.config.VADType != "test" {
		t.Errorf("Expected VADType to be 'test', got %s", pool.config.VADType)
	}
}

func TestPoolGetStats(t *testing.T) {
	config := &PoolConfig{
		PoolSize:  10,
		MaxIdle:   5,
		Threshold: 0.5,
		VADType:   "test",
	}

	pool := NewPool(config)
	pool.Initialize()

	stats := pool.GetStats()

	if stats["vad_type"] != "test" {
		t.Error("Expected vad_type to be 'test' in stats")
	}

	if stats["pool_size"] != 10 {
		t.Error("Expected pool_size to be 10 in stats")
	}

	// 清理
	pool.Shutdown()
}

func TestInstance(t *testing.T) {
	instance := &Instance{
		ID:   1,
		Type: "test",
	}

	if instance.GetID() != 1 {
		t.Errorf("Expected ID to be 1, got %d", instance.GetID())
	}

	if instance.GetType() != "test" {
		t.Errorf("Expected Type to be 'test', got %s", instance.GetType())
	}

	// 测试使用状态
	instance.SetInUse(true)
	if !instance.IsInUse() {
		t.Error("Expected instance to be in use")
	}

	instance.SetInUse(false)
	if instance.IsInUse() {
		t.Error("Expected instance to not be in use")
	}
}

func TestNewVADFactory(t *testing.T) {
	factory := NewVADFactory()

	if factory == nil {
		t.Fatal("Expected factory to be created")
	}

	types := factory.GetSupportedTypes()
	if types == nil {
		t.Error("Expected supported types to be returned")
	}
}

