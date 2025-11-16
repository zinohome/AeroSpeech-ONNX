package vad

// VADPoolInterface VAD池接口
type VADPoolInterface interface {
	// Initialize 初始化池
	Initialize() error

	// Get 获取VAD实例
	Get() (VADInstanceInterface, error)

	// Put 归还VAD实例
	Put(instance VADInstanceInterface)

	// GetStats 获取统计信息
	GetStats() map[string]interface{}

	// Shutdown 关闭池
	Shutdown()
}

// VADInstanceInterface VAD实例接口
type VADInstanceInterface interface {
	// GetID 获取实例ID
	GetID() int

	// GetType 获取VAD类型
	GetType() string

	// IsInUse 检查是否在使用中
	IsInUse() bool

	// SetInUse 设置使用状态
	SetInUse(inUse bool)

	// GetLastUsed 获取最后使用时间
	GetLastUsed() int64

	// SetLastUsed 设置最后使用时间
	SetLastUsed(timestamp int64)

	// Reset 重置实例状态
	Reset() error

	// Destroy 销毁实例
	Destroy() error

	// Process 处理音频数据
	Process(audio []float32) (bool, error)
}

// VADFactory VAD工厂接口
type VADFactory interface {
	// CreatePool 创建VAD池
	CreatePool(config interface{}) (VADPoolInterface, error)

	// GetSupportedTypes 获取支持的VAD类型
	GetSupportedTypes() []string
}

