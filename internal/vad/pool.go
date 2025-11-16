package vad

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
)

// Pool VAD资源池基础实现
type Pool struct {
	instances []*Instance
	available chan VADInstanceInterface
	config    *PoolConfig
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc

	// 统计信息
	totalCreated int64
	totalReused  int64
	totalActive  int64
}

// PoolConfig VAD池配置
type PoolConfig struct {
	PoolSize  int
	MaxIdle   int
	Threshold float32
	VADType   string
}

// Instance VAD实例基础实现
type Instance struct {
	ID       int
	Type     string
	LastUsed int64
	InUse    int32
	mu       sync.RWMutex
}

// GetID 获取实例ID
func (i *Instance) GetID() int {
	return i.ID
}

// GetType 获取VAD类型
func (i *Instance) GetType() string {
	return i.Type
}

// IsInUse 检查是否在使用中
func (i *Instance) IsInUse() bool {
	return atomic.LoadInt32(&i.InUse) == 1
}

// SetInUse 设置使用状态
func (i *Instance) SetInUse(inUse bool) {
	if inUse {
		atomic.StoreInt32(&i.InUse, 1)
	} else {
		atomic.StoreInt32(&i.InUse, 0)
	}
}

// GetLastUsed 获取最后使用时间
func (i *Instance) GetLastUsed() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.LastUsed
}

// SetLastUsed 设置最后使用时间
func (i *Instance) SetLastUsed(timestamp int64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.LastUsed = timestamp
}

// Reset 重置实例状态
func (i *Instance) Reset() error {
	// 基础实现，子类可以重写
	return nil
}

// Destroy 销毁实例
func (i *Instance) Destroy() error {
	// 基础实现，子类可以重写
	logger.Infof("VAD instance %d destroyed", i.ID)
	return nil
}

// Process 处理音频数据
func (i *Instance) Process(audio []float32) (bool, error) {
	// 基础实现，子类必须重写
	return false, fmt.Errorf("not implemented")
}

// NewPool 创建新的VAD资源池
func NewPool(config *PoolConfig) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		instances: make([]*Instance, 0, config.PoolSize),
		available: make(chan VADInstanceInterface, config.PoolSize),
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
	}

	return pool
}

// Initialize 并行初始化VAD池
func (p *Pool) Initialize() error {
	logger.Infof("Initializing VAD pool with %d instances...", p.config.PoolSize)

	// 由于这是基础实现，实际的VAD实例创建应该由具体的VAD类型实现
	// 这里只提供框架

	logger.Info("VAD pool initialized (base implementation)")
	return nil
}

// Get 获取VAD实例
func (p *Pool) Get() (VADInstanceInterface, error) {
	select {
	case instance := <-p.available:
		if atomic.CompareAndSwapInt32(&instance.(*Instance).InUse, 0, 1) {
			instance.SetLastUsed(time.Now().UnixNano())
			atomic.AddInt64(&p.totalReused, 1)
			atomic.AddInt64(&p.totalActive, 1)
			return instance, nil
		}
		// 实例已被使用，重新放回队列
		select {
		case p.available <- instance:
		default:
		}
		return p.Get() // 递归重试

	case <-time.After(100 * time.Millisecond):
		// 超时，返回错误
		return nil, fmt.Errorf("VAD pool timeout")

	case <-p.ctx.Done():
		return nil, fmt.Errorf("VAD pool is shutting down")
	}
}

// Put 归还VAD实例
func (p *Pool) Put(instance VADInstanceInterface) {
	if instance == nil {
		logger.Warn("Attempted to put nil VAD instance")
		return
	}

	if atomic.CompareAndSwapInt32(&instance.(*Instance).InUse, 1, 0) {
		instance.SetLastUsed(time.Now().UnixNano())
		atomic.AddInt64(&p.totalActive, -1)

		// 重置VAD状态
		if err := instance.Reset(); err != nil {
			logger.Warnf("Failed to reset VAD instance %d: %v", instance.GetID(), err)
		}

		select {
		case p.available <- instance:
			// 成功归还
		default:
			// 队列满，销毁实例
			logger.Warnf("VAD pool queue full, destroying instance %d", instance.GetID())
			instance.Destroy()
		}
	}
}

// GetStats 获取统计信息
func (p *Pool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"vad_type":        p.config.VADType,
		"pool_size":       p.config.PoolSize,
		"max_idle":        p.config.MaxIdle,
		"total_instances": len(p.instances),
		"available_count": len(p.available),
		"active_count":    atomic.LoadInt64(&p.totalActive),
		"total_created":   atomic.LoadInt64(&p.totalCreated),
		"total_reused":    atomic.LoadInt64(&p.totalReused),
	}
}

// Shutdown 关闭VAD池
func (p *Pool) Shutdown() {
	logger.Info("Shutting down VAD pool...")

	// 取消上下文
	p.cancel()

	// 销毁所有实例
	p.mu.Lock()
	defer p.mu.Unlock()

	// 清空可用队列
	for {
		select {
		case instance := <-p.available:
			instance.Destroy()
		default:
			goto cleanup_instances
		}
	}

cleanup_instances:
	// 销毁所有实例
	for _, instance := range p.instances {
		instance.Destroy()
	}

	p.instances = nil
	close(p.available)

	logger.Info("VAD pool shutdown complete")
}

