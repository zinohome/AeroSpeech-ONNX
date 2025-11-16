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

// isVitsModel 检测是否为VITS/Piper模型
func isVitsModel(modelPath string) bool {
	// 通过文件名判断模型类型
	modelName := strings.ToLower(filepath.Base(modelPath))
	return strings.Contains(modelName, "vits") || 
	       strings.Contains(modelName, "piper") ||
	       strings.Contains(modelName, "huayan")
}

// NewTTSProvider 创建TTS Provider
func NewTTSProvider(cfg *config.TTSModelConfig) (*TTSProvider, error) {
	// 检查模型文件是否存在
	if cfg.ModelPath != "" {
		if _, err := os.Stat(cfg.ModelPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("TTS model file not found: %s (please check if the model file exists or download models using scripts/download_models.sh)", cfg.ModelPath)
		}
	}
	
	// 检查tokens文件是否存在
	if cfg.TokensPath != "" {
		if _, err := os.Stat(cfg.TokensPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("TTS tokens file not found: %s (please check if the tokens file exists or download models using scripts/download_models.sh)", cfg.TokensPath)
		}
	}
	
	// 检查voices文件是否存在（如果配置了）
	if cfg.VoicesPath != "" {
		if _, err := os.Stat(cfg.VoicesPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("TTS voices file not found: %s (please check if the voices file exists or download models using scripts/download_models.sh)", cfg.VoicesPath)
		}
	}
	
	// 构建sherpa-onnx配置
	sampleRate := 24000 // 默认采样率（Kokoro）
	
	// 判断模型类型
	useVits := isVitsModel(cfg.ModelPath)
	
	var kokoroConfig sherpa.OfflineTtsKokoroModelConfig
	var vitsConfig sherpa.OfflineTtsVitsModelConfig
	
	if useVits {
		// VITS/Piper 模型配置
		sampleRate = 22050 // Piper 模型默认采样率
		
		vitsConfig = sherpa.OfflineTtsVitsModelConfig{
			Model:   cfg.ModelPath,
			Tokens:  cfg.TokensPath,
			DataDir: cfg.DataDir,
		}
		
		// 设置 Lexicon（可选）
		if cfg.Lexicon != "" {
			vitsConfig.Lexicon = cfg.Lexicon
		}
		
		// 设置 DictDir（可选）
		if cfg.DictDir != "" {
			vitsConfig.DictDir = cfg.DictDir
		}
	} else {
		// Kokoro 模型配置
		kokoroConfig = sherpa.OfflineTtsKokoroModelConfig{
			Model: cfg.ModelPath,
		}
	}
	
	// 共同配置：DataDir 处理
	dataDir := cfg.DataDir
	if dataDir == "" && cfg.ModelPath != "" {
		// 尝试从model路径推断dataDir路径
		modelDir := filepath.Dir(cfg.ModelPath)
		potentialDataDir := filepath.Join(modelDir, "espeak-ng-data")
		if _, err := os.Stat(potentialDataDir); err == nil {
			dataDir = potentialDataDir
		}
	}
	
	// 检查dataDir是否存在以及必需的phontab文件
	if dataDir != "" {
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("TTS data directory not found: %s (please check if the espeak-ng-data directory exists or download models using scripts/download_models.sh)", dataDir)
		}
		// 检查必需的phontab文件
		phontabPath := filepath.Join(dataDir, "phontab")
		if _, err := os.Stat(phontabPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("TTS phontab file not found: %s (required file missing in espeak-ng-data directory, please download models using scripts/download_models.sh)", phontabPath)
		}
	} else {
		return nil, fmt.Errorf("TTS data_dir is required but not specified in config and cannot be inferred from model path (please set data_dir in config or ensure espeak-ng-data directory exists in model directory)")
	}
	
	// 根据模型类型完成配置
	if !useVits {
		// Kokoro 模型的额外配置
		if cfg.TokensPath != "" {
			kokoroConfig.Tokens = cfg.TokensPath
		}
		
		// 设置Voices（可选，用于多说话人）
		if cfg.VoicesPath != "" {
			kokoroConfig.Voices = cfg.VoicesPath
		}
		
		kokoroConfig.DataDir = dataDir
		
		// 设置DictDir（可选）
		if cfg.DictDir != "" {
			kokoroConfig.DictDir = cfg.DictDir
		}
		
		// 设置Lexicon（可选）
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
	} else {
		// VITS 模型已在上面配置
		vitsConfig.DataDir = dataDir
	}
	
	// 创建 TTS 配置
	var ttsConfig sherpa.OfflineTtsConfig
	if useVits {
		ttsConfig = sherpa.OfflineTtsConfig{
			Model: sherpa.OfflineTtsModelConfig{
				Vits:       vitsConfig,
				Provider:   config.GetProvider(&cfg.Provider),
				NumThreads: cfg.Provider.NumThreads,
				Debug:      0,
			},
			MaxNumSentences: 1,
		}
	} else {
		ttsConfig = sherpa.OfflineTtsConfig{
			Model: sherpa.OfflineTtsModelConfig{
				Kokoro:     kokoroConfig,
				Provider:   config.GetProvider(&cfg.Provider),
				NumThreads: cfg.Provider.NumThreads,
				Debug:      0,
			},
			MaxNumSentences: 1,
		}
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

