package asr

import (
	"fmt"
	"os"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/pkg/utils"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// Provider ASR Provider接口
type Provider interface {
	Transcribe(audio []byte) (string, error)
	Warmup() error
	Reset() error
	Release() error
	GetSampleRate() int
}

// ASRProvider sherpa-onnx ASR Provider实现
type ASRProvider struct {
	recognizer *sherpa.OfflineRecognizer
	config     *config.ASRConfig
	sampleRate int
}

// NewASRProvider 创建ASR Provider
func NewASRProvider(cfg *config.ASRConfig) (*ASRProvider, error) {
	// 检查模型文件是否存在
	if cfg.ModelPath != "" {
		if _, err := os.Stat(cfg.ModelPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("model file not found: %s (please check if the model file exists or download models using scripts/download_models.sh)", cfg.ModelPath)
		}
	}
	
	// 检查tokens文件是否存在
	if cfg.TokensPath != "" {
		if _, err := os.Stat(cfg.TokensPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("tokens file not found: %s (please check if the tokens file exists or download models using scripts/download_models.sh)", cfg.TokensPath)
		}
	}
	
	// 构建sherpa-onnx配置
	sampleRate := 16000 // 默认采样率
	recognizerConfig := sherpa.OfflineRecognizerConfig{
		FeatConfig: sherpa.FeatureConfig{
			SampleRate: sampleRate,
			FeatureDim: 80,
		},
		ModelConfig: sherpa.OfflineModelConfig{
			Tokens: cfg.TokensPath,
			NumThreads: cfg.Provider.NumThreads,
			Debug: 0,
			Provider: config.GetProvider(&cfg.Provider),
		},
		DecodingMethod: "greedy_search",
		MaxActivePaths: 4,
	}

	// 根据模型类型设置配置
	// 这里使用SenseVoice作为默认模型
	recognizerConfig.ModelConfig.SenseVoice = sherpa.OfflineSenseVoiceModelConfig{
		Model: cfg.ModelPath,
		Language: cfg.Language,
		UseInverseTextNormalization: 1,
	}

	// 创建识别器
	recognizer := sherpa.NewOfflineRecognizer(&recognizerConfig)
	if recognizer == nil {
		return nil, fmt.Errorf("failed to create offline recognizer")
	}

	provider := &ASRProvider{
		recognizer: recognizer,
		config:     cfg,
		sampleRate: 16000, // 默认采样率
	}

	return provider, nil
}

// Transcribe 识别音频
func (p *ASRProvider) Transcribe(audio []byte) (string, error) {
	if len(audio) == 0 {
		return "", fmt.Errorf("audio data is empty")
	}

	// 转换音频数据
	samples := utils.SamplesInt16ToFloat(audio)
	if samples == nil {
		return "", fmt.Errorf("failed to convert audio data")
	}

	// 创建识别流
	stream := sherpa.NewOfflineStream(p.recognizer)
	if stream == nil {
		return "", fmt.Errorf("failed to create offline stream")
	}
	defer sherpa.DeleteOfflineStream(stream)

	// 接受音频数据
	stream.AcceptWaveform(p.sampleRate, samples)

	// 执行识别
	p.recognizer.Decode(stream)

	// 获取结果
	result := stream.GetResult()
	if result == nil {
		// result为nil可能是正常情况（例如空音频或静音），返回空字符串而不是错误
		return "", nil
	}
	
	// 检查result.Text是否为空
	if result.Text == "" {
		return "", nil // 空结果也是有效结果
	}
	
	return result.Text, nil
}

// Warmup 预热模型
func (p *ASRProvider) Warmup() error {
	// 使用空音频进行预热
	// 注意：空音频可能返回空结果，这是正常的
	// 只要recognizer能正常工作（不崩溃），就认为warmup成功
	dummyAudio := make([]byte, 1600) // 0.1秒的音频（16kHz, 16-bit）
	_, err := p.Transcribe(dummyAudio)
	// Transcribe现在对空结果返回nil错误，所以这里应该总是成功
	// 如果真的有错误（比如recognizer内部错误），会在实际使用时发现
	return err
}

// Reset 重置Provider
func (p *ASRProvider) Reset() error {
	// ASR Provider通常不需要重置
	return nil
}

// Release 释放资源
func (p *ASRProvider) Release() error {
	if p.recognizer != nil {
		sherpa.DeleteOfflineRecognizer(p.recognizer)
		p.recognizer = nil
	}
	return nil
}

// GetSampleRate 获取采样率
func (p *ASRProvider) GetSampleRate() int {
	return p.sampleRate
}

