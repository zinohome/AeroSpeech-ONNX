package asr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
)

// Pool ASR资源池
type Pool struct {
	providers   chan Provider
	config      *config.ASRConfig
	size        int
	mu          sync.RWMutex
	stats       *PoolStats
	ctx         context.Context
	cancel      context.CancelFunc
}

// PoolStats 资源池统计信息
type PoolStats struct {
	TotalCreated   int
	TotalDestroyed int
	CurrentActive  int
	MaxWaitTime    time.Duration
	TotalWaits     int64
	mu             sync.RWMutex
}

// NewPool 创建ASR资源池
func NewPool(cfg *config.ASRConfig, size int) (*Pool, error) {
	if size <= 0 {
		size = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		providers: make(chan Provider, size),
		config:    cfg,
		size:      size,
		stats: &PoolStats{
			CurrentActive: 0,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 预创建Provider
	for i := 0; i < size; i++ {
		provider, err := NewASRProvider(cfg)
		if err != nil {
			logger.Warnf("Failed to create ASR provider %d: %v", i, err)
			continue
		}

		// 预热Provider
		if err := provider.Warmup(); err != nil {
			logger.Warnf("Failed to warmup ASR provider %d: %v", i, err)
			provider.Release()
			continue
		}

		pool.providers <- provider
		pool.stats.mu.Lock()
		pool.stats.TotalCreated++
		pool.stats.CurrentActive++
		pool.stats.mu.Unlock()
	}

	if len(pool.providers) == 0 {
		cancel()
		return nil, fmt.Errorf("failed to create any ASR provider")
	}

	logger.Infof("ASR pool initialized with %d providers", len(pool.providers))

	return pool, nil
}

// Get 从资源池获取Provider
func (p *Pool) Get(ctx context.Context) (Provider, error) {
	startTime := time.Now()

	select {
	case provider := <-p.providers:
		waitTime := time.Since(startTime)
		p.stats.mu.Lock()
		if waitTime > p.stats.MaxWaitTime {
			p.stats.MaxWaitTime = waitTime
		}
		p.stats.TotalWaits++
		p.stats.CurrentActive--
		p.stats.mu.Unlock()
		return provider, nil

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-p.ctx.Done():
		return nil, fmt.Errorf("pool is closed")
	}
}

// Put 归还Provider到资源池
func (p *Pool) Put(provider Provider) {
	if provider == nil {
		return
	}

	select {
	case p.providers <- provider:
		p.stats.mu.Lock()
		p.stats.CurrentActive++
		p.stats.mu.Unlock()
	default:
		// 池已满，释放Provider
		provider.Release()
		p.stats.mu.Lock()
		p.stats.TotalDestroyed++
		p.stats.mu.Unlock()
	}
}

// GetUsage 获取资源池使用率
func (p *Pool) GetUsage() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	active := len(p.providers)
	if p.size == 0 {
		return 0
	}

	return float64(active) / float64(p.size)
}

// GetStats 获取资源池统计信息
func (p *Pool) GetStats() map[string]interface{} {
	p.stats.mu.RLock()
	defer p.stats.mu.RUnlock()

	return map[string]interface{}{
		"size":            p.size,
		"available":       len(p.providers),
		"total_created":  p.stats.TotalCreated,
		"total_destroyed": p.stats.TotalDestroyed,
		"current_active": p.stats.CurrentActive,
		"max_wait_time":  p.stats.MaxWaitTime.String(),
		"total_waits":    p.stats.TotalWaits,
		"usage":          p.GetUsage(),
	}
}

// Close 关闭资源池
func (p *Pool) Close() error {
	p.cancel()

	// 关闭通道
	close(p.providers)

	// 释放所有Provider
	for provider := range p.providers {
		provider.Release()
		p.stats.mu.Lock()
		p.stats.TotalDestroyed++
		p.stats.mu.Unlock()
	}

	logger.Info("ASR pool closed")
	return nil
}

