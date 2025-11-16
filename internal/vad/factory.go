package vad

import (
	"fmt"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/logger"
)

// DefaultVADFactory 默认VAD工厂实现
type DefaultVADFactory struct {
	factories map[string]VADFactory
}

// NewVADFactory 创建新的VAD工厂
func NewVADFactory() *DefaultVADFactory {
	return &DefaultVADFactory{
		factories: make(map[string]VADFactory),
	}
}

// RegisterFactory 注册VAD池工厂
func (f *DefaultVADFactory) RegisterFactory(vadType string, factory VADFactory) {
	f.factories[vadType] = factory
	logger.Infof("Registered VAD factory for type: %s", vadType)
}

// CreateVADPool 根据配置创建VAD池
func (f *DefaultVADFactory) CreateVADPool(vadType string, config interface{}) (VADPoolInterface, error) {
	logger.Infof("Creating VAD pool with type: %s", vadType)

	factory, exists := f.factories[vadType]
	if !exists {
		return nil, fmt.Errorf("unsupported VAD type: %s", vadType)
	}

	// 使用工厂创建池
	pool, err := factory.CreatePool(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s VAD pool: %v", vadType, err)
	}

	return pool, nil
}

// GetVADType 获取当前VAD类型
func (f *DefaultVADFactory) GetVADType(vadType string) (VADFactory, error) {
	factory, exists := f.factories[vadType]
	if !exists {
		return nil, fmt.Errorf("unsupported VAD type: %s", vadType)
	}
	return factory, nil
}

// GetSupportedTypes 获取支持的VAD类型
func (f *DefaultVADFactory) GetSupportedTypes() []string {
	types := make([]string, 0, len(f.factories))
	for vadType := range f.factories {
		types = append(types, vadType)
	}
	return types
}

