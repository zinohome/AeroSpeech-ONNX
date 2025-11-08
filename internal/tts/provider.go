package tts

import (
	"fmt"

	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/config"
	"github.com/zhangjun/AeroSpeech-ONNX/pkg/utils"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// Provider TTS Provider接口
type Provider interface {
	Synthesize(text string, speakerID int, speed float32) ([]byte, error)
	Warmup() error
	Reset() error
	Release() error
	GetSampleRate() int
}

// TTSProvider sherpa-onnx TTS Provider实现
type TTSProvider struct {
	tts        *sherpa.OfflineTts
	config     *config.TTSModelConfig
	sampleRate int
}

// NewTTSProvider 创建TTS Provider
func NewTTSProvider(cfg *config.TTSModelConfig) (*TTSProvider, error) {
	// 构建sherpa-onnx配置
	sampleRate := 24000 // 默认采样率
	ttsConfig := sherpa.OfflineTtsConfig{
		Model: sherpa.OfflineTtsModelConfig{
			Kokoro: sherpa.OfflineTtsKokoroModelConfig{
				Model: cfg.ModelPath,
			},
			Provider:   config.GetProvider(&cfg.Provider),
			NumThreads: cfg.Provider.NumThreads,
			Debug:      0,
		},
		MaxNumSentences: 1,
	}

	// 创建TTS合成器
	tts := sherpa.NewOfflineTts(&ttsConfig)
	if tts == nil {
		return nil, fmt.Errorf("failed to create offline TTS")
	}

	provider := &TTSProvider{
		tts:        tts,
		config:     cfg,
		sampleRate: sampleRate,
	}

	return provider, nil
}

// Synthesize 合成语音
func (p *TTSProvider) Synthesize(text string, speakerID int, speed float32) ([]byte, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}

	// 生成音频
	audio := p.tts.Generate(text, speakerID, speed)
	if audio == nil {
		return nil, fmt.Errorf("failed to generate audio")
	}

	// 转换音频数据
	samples := audio.Samples
	audioData := utils.SamplesFloatToInt16(samples)

	return audioData, nil
}

// Warmup 预热模型
func (p *TTSProvider) Warmup() error {
	// 使用简单文本进行预热
	_, err := p.Synthesize("测试", 0, 1.0)
	return err
}

// Reset 重置Provider
func (p *TTSProvider) Reset() error {
	// TTS Provider通常不需要重置
	return nil
}

// Release 释放资源
func (p *TTSProvider) Release() error {
	if p.tts != nil {
		sherpa.DeleteOfflineTts(p.tts)
		p.tts = nil
	}
	return nil
}

// GetSampleRate 获取采样率
func (p *TTSProvider) GetSampleRate() int {
	return p.sampleRate
}

