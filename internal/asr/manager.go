package asr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
)

// Manager ASR管理器
type Manager struct {
	pool      *Pool
	config    *config.ASRConfig
	stats     *Stats
	statsMu   sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// Stats 统计信息
type Stats struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	TotalLatency       time.Duration
	LatencyHistory     []time.Duration
	LastRequestTime    time.Time
	mu                 sync.RWMutex
}

// NewManager 创建ASR管理器
func NewManager(cfg *config.ASRConfig, poolSize int) (*Manager, error) {
	pool, err := NewPool(cfg, poolSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create ASR pool: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	manager := &Manager{
		pool:   pool,
		config: cfg,
		stats: &Stats{
			LatencyHistory: make([]time.Duration, 0, 1000),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 启动统计信息清理goroutine
	go manager.cleanupStats()

	return manager, nil
}

// Transcribe 识别音频
func (m *Manager) Transcribe(ctx interface{}, audio []byte) (string, error) {
	startTime := time.Now()

	// 从资源池获取Provider
	// 如果ctx是context.Context，使用它；否则使用context.Background()
	var poolCtx context.Context
	if ctxCtx, ok := ctx.(context.Context); ok {
		poolCtx = ctxCtx
	} else {
		poolCtx = context.Background()
	}
	provider, err := m.pool.Get(poolCtx)
	if err != nil {
		m.recordFailure()
		return "", fmt.Errorf("failed to get provider from pool: %w", err)
	}
	defer m.pool.Put(provider)

	// 执行识别
	result, err := provider.Transcribe(audio)
	latency := time.Since(startTime)

	if err != nil {
		m.recordFailure()
		logger.Errorf("ASR transcription failed: %v", err)
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	m.recordSuccess(latency)
	return result, nil
}

// recordSuccess 记录成功请求
func (m *Manager) recordSuccess(latency time.Duration) {
	m.statsMu.Lock()
	defer m.statsMu.Unlock()

	m.stats.TotalRequests++
	m.stats.SuccessfulRequests++
	m.stats.TotalLatency += latency
	m.stats.LastRequestTime = time.Now()

	// 记录延迟历史（保留最近1000条）
	if len(m.stats.LatencyHistory) >= 1000 {
		m.stats.LatencyHistory = m.stats.LatencyHistory[1:]
	}
	m.stats.LatencyHistory = append(m.stats.LatencyHistory, latency)
}

// recordFailure 记录失败请求
func (m *Manager) recordFailure() {
	m.statsMu.Lock()
	defer m.statsMu.Unlock()

	m.stats.TotalRequests++
	m.stats.FailedRequests++
	m.stats.LastRequestTime = time.Now()
}

// GetStats 获取统计信息
func (m *Manager) GetStats() interface{} {
	m.statsMu.RLock()
	defer m.statsMu.RUnlock()

	// 返回统计信息的副本
	stats := &Stats{
		TotalRequests:      m.stats.TotalRequests,
		SuccessfulRequests: m.stats.SuccessfulRequests,
		FailedRequests:     m.stats.FailedRequests,
		TotalLatency:       m.stats.TotalLatency,
		LastRequestTime:    m.stats.LastRequestTime,
	}

	// 复制延迟历史
	stats.LatencyHistory = make([]time.Duration, len(m.stats.LatencyHistory))
	copy(stats.LatencyHistory, m.stats.LatencyHistory)

	return stats
}

// GetAvgLatency 获取平均延迟
func (m *Manager) GetAvgLatency() interface{} {
	m.statsMu.RLock()
	defer m.statsMu.RUnlock()

	if m.stats.SuccessfulRequests == 0 {
		return time.Duration(0)
	}

	return m.stats.TotalLatency / time.Duration(m.stats.SuccessfulRequests)
}

// GetPoolUsage 获取资源池使用率
func (m *Manager) GetPoolUsage() float64 {
	return m.pool.GetUsage()
}

// GetPoolStats 获取资源池统计信息
func (m *Manager) GetPoolStats() map[string]interface{} {
	return m.pool.GetStats()
}

// cleanupStats 定期清理统计信息
func (m *Manager) cleanupStats() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.statsMu.Lock()
			// 保留最近1000条延迟历史
			validHistory := make([]time.Duration, 0)
			for _, latency := range m.stats.LatencyHistory {
				// 简化处理：只保留最近1000条
				if len(validHistory) < 1000 {
					validHistory = append(validHistory, latency)
				}
			}
			m.stats.LatencyHistory = validHistory
			m.statsMu.Unlock()
		}
	}
}

// Close 关闭管理器
func (m *Manager) Close() error {
	m.cancel()
	return m.pool.Close()
}

