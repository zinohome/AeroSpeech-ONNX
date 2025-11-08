package tts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	
	// 构建Kokoro模型配置
	kokoroConfig := sherpa.OfflineTtsKokoroModelConfig{
		Model: cfg.ModelPath,
	}
	
	// 设置Tokens（必需）
	if cfg.TokensPath != "" {
		kokoroConfig.Tokens = cfg.TokensPath
	}
	
	// 设置Voices（可选，用于多说话人）
	// 注意：只有当VoicesPath非空时才设置，否则sherpa-onnx会尝试加载空字符串导致错误
	if cfg.VoicesPath != "" {
		kokoroConfig.Voices = cfg.VoicesPath
	}
	// 如果VoicesPath为空，不设置Voices字段，让sherpa-onnx使用默认行为
	
	// 设置DataDir（必需，用于文本处理）
	// 如果配置中没有指定，尝试使用默认路径
	dataDir := cfg.DataDir
	if dataDir == "" && cfg.ModelPath != "" {
		// 尝试从model路径推断dataDir路径
		// 例如：如果model在 models/tts/kokoro-multi-lang-v1_0/model.onnx
		// 则dataDir应该在 models/tts/kokoro-multi-lang-v1_0/espeak-ng-data
		modelDir := filepath.Dir(cfg.ModelPath)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}
	if dataDir != "" {
		kokoroConfig.DataDir = dataDir
	}
	
	// 设置DictDir（可选，用于字典文件）
	if cfg.DictDir != "" {
		kokoroConfig.DictDir = cfg.DictDir
	}
	
	// 设置Lexicon（可选，用于多语言支持）
	// 如果配置中有多个lexicon文件，用逗号分隔
	if cfg.Lexicon != "" {
		kokoroConfig.Lexicon = cfg.Lexicon
	} else if cfg.ModelPath != "" {
		// 尝试从model路径推断lexicon路径
		modelDir := filepath.Dir(cfg.ModelPath)
		var lexiconPaths []string
		for _, name := range []string{"lexicon-us-en.txt", "lexicon-zh.txt"} {
			potentialLexicon := filepath.Join(modelDir, name)
			if _, err := os.Stat(potentialLexicon); err == nil {
				lexiconPaths = append(lexiconPaths, potentialLexicon)
			}
		}
		if len(lexiconPaths) > 0 {
			kokoroConfig.Lexicon = strings.Join(lexiconPaths, ",")
		}
	}
	
	ttsConfig := sherpa.OfflineTtsConfig{
		Model: sherpa.OfflineTtsModelConfig{
			Kokoro:     kokoroConfig,
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

