package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// ProviderConfig Provider配置
type ProviderConfig struct {
	Provider   string `mapstructure:"provider" json:"provider"`     // "cpu", "cuda", "auto"
	DeviceID   int    `mapstructure:"device_id" json:"device_id"`   // GPU设备ID（默认0）
	NumThreads int    `mapstructure:"num_threads" json:"num_threads"` // 线程数
}

// ASRConfig ASR配置
type ASRConfig struct {
	ModelPath  string          `mapstructure:"model_path" json:"model_path"`
	TokensPath string          `mapstructure:"tokens_path" json:"tokens_path"`
	Language   string          `mapstructure:"language" json:"language"`
	Provider   ProviderConfig  `mapstructure:"provider" json:"provider"`
	Debug      bool            `mapstructure:"debug" json:"debug"`
}

// TTSModelConfig TTS模型配置
type TTSModelConfig struct {
	ModelPath  string         `mapstructure:"model_path" json:"model_path"`
	VoicesPath string         `mapstructure:"voices_path" json:"voices_path"`
	TokensPath string         `mapstructure:"tokens_path" json:"tokens_path"`
	DataDir    string         `mapstructure:"data_dir" json:"data_dir"`
	DictDir    string         `mapstructure:"dict_dir" json:"dict_dir"`
	Lexicon    string         `mapstructure:"lexicon" json:"lexicon"` // 逗号分隔的lexicon文件路径
	Provider   ProviderConfig `mapstructure:"provider" json:"provider"`
	Debug      bool           `mapstructure:"debug" json:"debug"`
}

// AudioConfig 音频配置
type AudioConfig struct {
	SampleRate     int     `mapstructure:"sample_rate" json:"sample_rate"`
	FeatureDim     int     `mapstructure:"feature_dim" json:"feature_dim"`
	ChunkSize      int     `mapstructure:"chunk_size" json:"chunk_size"`
	NormalizeFactor float64 `mapstructure:"normalize_factor" json:"normalize_factor"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	ReadTimeout int    `mapstructure:"read_timeout" json:"read_timeout"`
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	ReadTimeout      int  `mapstructure:"read_timeout" json:"read_timeout"`
	MaxMessageSize   int  `mapstructure:"max_message_size" json:"max_message_size"`
	ReadBufferSize   int  `mapstructure:"read_buffer_size" json:"read_buffer_size"`
	WriteBufferSize  int  `mapstructure:"write_buffer_size" json:"write_buffer_size"`
	EnableCompression bool `mapstructure:"enable_compression" json:"enable_compression"`
}

// SessionConfig 会话配置
type SessionConfig struct {
	SendQueueSize int `mapstructure:"send_queue_size" json:"send_queue_size"`
	MaxSendErrors int `mapstructure:"max_send_errors" json:"max_send_errors"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Format     string `mapstructure:"format" json:"format"`
	Output     string `mapstructure:"output" json:"output"`
	FilePath   string `mapstructure:"file_path" json:"file_path"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled" json:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second" json:"requests_per_second"`
	BurstSize         int  `mapstructure:"burst_size" json:"burst_size"`
	MaxConnections    int  `mapstructure:"max_connections" json:"max_connections"`
}

// VADConfig VAD配置
type VADConfig struct {
	Enabled   bool    `mapstructure:"enabled" json:"enabled"`
	Provider  string  `mapstructure:"provider" json:"provider"` // "silero", "ten", etc.
	PoolSize  int     `mapstructure:"pool_size" json:"pool_size"`
	Threshold float32 `mapstructure:"threshold" json:"threshold"`
}

// STTConfig STT服务配置
type STTConfig struct {
	Server    ServerConfig    `mapstructure:"server" json:"server"`
	ASR       ASRConfig       `mapstructure:"asr" json:"asr"`
	Audio     AudioConfig     `mapstructure:"audio" json:"audio"`
	WebSocket WebSocketConfig `mapstructure:"websocket" json:"websocket"`
	Session   SessionConfig   `mapstructure:"session" json:"session"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit" json:"rate_limit"`
	VAD       VADConfig       `mapstructure:"vad" json:"vad"`
	Logging   LoggingConfig   `mapstructure:"logging" json:"logging"`
}

// TTSConfig TTS服务配置
type TTSConfig struct {
	Server    ServerConfig    `mapstructure:"server" json:"server"`
	TTS       TTSModelConfig  `mapstructure:"tts" json:"tts"`
	Audio     AudioConfig     `mapstructure:"audio" json:"audio"`
	WebSocket WebSocketConfig `mapstructure:"websocket" json:"websocket"`
	Session   SessionConfig   `mapstructure:"session" json:"session"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit" json:"rate_limit"`
	Logging   LoggingConfig   `mapstructure:"logging" json:"logging"`
}

// UnifiedConfig 统一配置（同时支持STT和TTS）
type UnifiedConfig struct {
	Mode      string          `mapstructure:"mode" json:"mode"` // "unified" 或 "separated"
	Server    ServerConfig    `mapstructure:"server" json:"server"`
	STT       *ASRConfig      `mapstructure:"stt" json:"stt,omitempty"`
	TTS       *TTSModelConfig `mapstructure:"tts" json:"tts,omitempty"`
	Audio     AudioConfig     `mapstructure:"audio" json:"audio"`
	WebSocket WebSocketConfig `mapstructure:"websocket" json:"websocket"`
	Session   SessionConfig   `mapstructure:"session" json:"session"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit" json:"rate_limit"`
	VAD       VADConfig       `mapstructure:"vad" json:"vad"`
	Logging   LoggingConfig   `mapstructure:"logging" json:"logging"`
}

// GlobalConfig 全局配置（STT或TTS）
var GlobalConfig interface{}

// LoadSTTConfig 加载STT服务配置
func LoadSTTConfig(configPath string) (*STTConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	// 支持环境变量
	viper.SetEnvPrefix("STT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config STTConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	setSTTDefaults(&config)

	// 验证配置
	if err := validateSTTConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Provider自动选择
	if err := resolveProvider(&config.ASR.Provider); err != nil {
		return nil, fmt.Errorf("failed to resolve provider: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// LoadTTSConfig 加载TTS服务配置
func LoadTTSConfig(configPath string) (*TTSConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	// 支持环境变量
	viper.SetEnvPrefix("TTS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TTSConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	setTTSDefaults(&config)

	// 验证配置
	if err := validateTTSConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Provider自动选择
	if err := resolveProvider(&config.TTS.Provider); err != nil {
		return nil, fmt.Errorf("failed to resolve provider: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// setSTTDefaults 设置STT配置默认值
func setSTTDefaults(config *STTConfig) {
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 20
	}

	if config.Audio.SampleRate == 0 {
		config.Audio.SampleRate = 16000
	}
	if config.Audio.FeatureDim == 0 {
		config.Audio.FeatureDim = 80
	}
	if config.Audio.ChunkSize == 0 {
		config.Audio.ChunkSize = 4096
	}
	if config.Audio.NormalizeFactor == 0 {
		config.Audio.NormalizeFactor = 32768.0
	}

	if config.ASR.Provider.Provider == "" {
		config.ASR.Provider.Provider = "cpu"
	}
	if config.ASR.Provider.NumThreads == 0 {
		if config.ASR.Provider.Provider == "cuda" {
			config.ASR.Provider.NumThreads = 1
		} else {
			config.ASR.Provider.NumThreads = runtime.NumCPU()
		}
	}

	if config.WebSocket.ReadTimeout == 0 {
		config.WebSocket.ReadTimeout = 20
	}
	if config.WebSocket.MaxMessageSize == 0 {
		config.WebSocket.MaxMessageSize = 2097152 // 2MB
	}
	if config.WebSocket.ReadBufferSize == 0 {
		config.WebSocket.ReadBufferSize = 1024
	}
	if config.WebSocket.WriteBufferSize == 0 {
		config.WebSocket.WriteBufferSize = 1024
	}

	if config.Session.SendQueueSize == 0 {
		config.Session.SendQueueSize = 500
	}
	if config.Session.MaxSendErrors == 0 {
		config.Session.MaxSendErrors = 10
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "text"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "console" // 默认只输出到控制台，避免文件路径问题
	}
	if config.Logging.FilePath == "" && (config.Logging.Output == "file" || config.Logging.Output == "both") {
		config.Logging.FilePath = "logs/stt.log" // 默认相对路径
	}
}

// setTTSDefaults 设置TTS配置默认值
func setTTSDefaults(config *TTSConfig) {
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8081
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 20
	}

	if config.Audio.SampleRate == 0 {
		config.Audio.SampleRate = 24000
	}
	if config.Audio.ChunkSize == 0 {
		config.Audio.ChunkSize = 4096
	}

	if config.TTS.Provider.Provider == "" {
		config.TTS.Provider.Provider = "cpu"
	}
	if config.TTS.Provider.NumThreads == 0 {
		if config.TTS.Provider.Provider == "cuda" {
			config.TTS.Provider.NumThreads = 1
		} else {
			config.TTS.Provider.NumThreads = 4
		}
	}

	if config.WebSocket.ReadTimeout == 0 {
		config.WebSocket.ReadTimeout = 20
	}
	if config.WebSocket.MaxMessageSize == 0 {
		config.WebSocket.MaxMessageSize = 2097152
	}

	if config.Session.SendQueueSize == 0 {
		config.Session.SendQueueSize = 500
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "text"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "console" // 默认只输出到控制台，避免文件路径问题
	}
	if config.Logging.FilePath == "" && (config.Logging.Output == "file" || config.Logging.Output == "both") {
		config.Logging.FilePath = "logs/tts.log" // 默认相对路径
	}
}

// validateSTTConfig 验证STT配置
func validateSTTConfig(config *STTConfig) error {
	if config.ASR.ModelPath == "" {
		return fmt.Errorf("asr.model_path is required")
	}
	if config.ASR.TokensPath == "" {
		return fmt.Errorf("asr.tokens_path is required")
	}

	// 检查模型文件是否存在
	if _, err := os.Stat(config.ASR.ModelPath); os.IsNotExist(err) {
		return fmt.Errorf("asr model file not found: %s", config.ASR.ModelPath)
	}
	if _, err := os.Stat(config.ASR.TokensPath); os.IsNotExist(err) {
		return fmt.Errorf("asr tokens file not found: %s", config.ASR.TokensPath)
	}

	// 验证Provider
	if config.ASR.Provider.Provider != "cpu" && 
		config.ASR.Provider.Provider != "cuda" && 
		config.ASR.Provider.Provider != "auto" {
		return fmt.Errorf("invalid provider: %s, must be cpu, cuda, or auto", config.ASR.Provider.Provider)
	}

	return nil
}

// validateTTSConfig 验证TTS配置
func validateTTSConfig(config *TTSConfig) error {
	if config.TTS.ModelPath == "" {
		return fmt.Errorf("tts.model_path is required")
	}

	// 检查模型文件是否存在
	if _, err := os.Stat(config.TTS.ModelPath); os.IsNotExist(err) {
		return fmt.Errorf("tts model file not found: %s", config.TTS.ModelPath)
	}

	// 验证Provider
	if config.TTS.Provider.Provider != "cpu" && 
		config.TTS.Provider.Provider != "cuda" && 
		config.TTS.Provider.Provider != "auto" {
		return fmt.Errorf("invalid provider: %s, must be cpu, cuda, or auto", config.TTS.Provider.Provider)
	}

	return nil
}

// resolveProvider 解析Provider配置（自动选择或回退）
func resolveProvider(provider *ProviderConfig) error {
	switch provider.Provider {
	case "auto":
		if isGPUAvailable() {
			provider.Provider = "cuda"
		} else {
			provider.Provider = "cpu"
		}
	case "cuda":
		if !isGPUAvailable() {
			// GPU不可用时回退到CPU
			provider.Provider = "cpu"
		}
	}
	return nil
}

// isGPUAvailable 检测GPU是否可用
func isGPUAvailable() bool {
	// 检查CUDA库是否存在
	// 这里可以调用nvidia-smi或检查CUDA库
	// 简化实现：检查环境变量或库文件
	if os.Getenv("CUDA_VISIBLE_DEVICES") != "" {
		return true
	}
	
	// 检查常见的CUDA库路径
	cudaPaths := []string{
		"/usr/local/cuda/lib64/libcudart.so",
		"/usr/lib/x86_64-linux-gnu/libcudart.so",
	}
	
	for _, path := range cudaPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	
	return false
}

// GetProvider 获取Provider配置
func GetProvider(provider *ProviderConfig) string {
	return provider.Provider
}

// GetDeviceID 获取GPU设备ID
func GetDeviceID(provider *ProviderConfig) int {
	return provider.DeviceID
}

// GetNumThreads 获取线程数
func GetNumThreads(provider *ProviderConfig) int {
	return provider.NumThreads
}

// LoadUnifiedConfig 加载统一配置
func LoadUnifiedConfig(configPath string) (*UnifiedConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	// 支持环境变量
	viper.SetEnvPrefix("SPEECH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config UnifiedConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	setUnifiedDefaults(&config)

	// 验证配置
	if err := validateUnifiedConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Provider自动选择
	if config.STT != nil {
		if err := resolveProvider(&config.STT.Provider); err != nil {
			return nil, fmt.Errorf("failed to resolve STT provider: %w", err)
		}
	}
	if config.TTS != nil {
		if err := resolveProvider(&config.TTS.Provider); err != nil {
			return nil, fmt.Errorf("failed to resolve TTS provider: %w", err)
		}
	}

	GlobalConfig = &config
	return &config, nil
}

// setUnifiedDefaults 设置统一配置默认值
func setUnifiedDefaults(config *UnifiedConfig) {
	if config.Mode == "" {
		config.Mode = "unified"
	}

	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 20
	}

	// STT默认值
	if config.STT != nil {
		if config.STT.Provider.Provider == "" {
			config.STT.Provider.Provider = "cpu"
		}
		if config.STT.Provider.NumThreads == 0 {
			if config.STT.Provider.Provider == "cuda" {
				config.STT.Provider.NumThreads = 1
			} else {
				config.STT.Provider.NumThreads = runtime.NumCPU()
			}
		}
	}

	// TTS默认值
	if config.TTS != nil {
		if config.TTS.Provider.Provider == "" {
			config.TTS.Provider.Provider = "cpu"
		}
		if config.TTS.Provider.NumThreads == 0 {
			if config.TTS.Provider.Provider == "cuda" {
				config.TTS.Provider.NumThreads = 1
			} else {
				config.TTS.Provider.NumThreads = 4
			}
		}
	}

	// 音频配置默认值
	if config.Audio.SampleRate == 0 {
		config.Audio.SampleRate = 16000 // STT默认采样率
	}
	if config.Audio.ChunkSize == 0 {
		config.Audio.ChunkSize = 4096
	}
	if config.Audio.FeatureDim == 0 {
		config.Audio.FeatureDim = 80
	}
	if config.Audio.NormalizeFactor == 0 {
		config.Audio.NormalizeFactor = 32768.0
	}

	// WebSocket配置默认值
	if config.WebSocket.ReadTimeout == 0 {
		config.WebSocket.ReadTimeout = 20
	}
	if config.WebSocket.MaxMessageSize == 0 {
		config.WebSocket.MaxMessageSize = 2097152 // 2MB
	}
	if config.WebSocket.ReadBufferSize == 0 {
		config.WebSocket.ReadBufferSize = 1024
	}
	if config.WebSocket.WriteBufferSize == 0 {
		config.WebSocket.WriteBufferSize = 1024
	}

	// 会话配置默认值
	if config.Session.SendQueueSize == 0 {
		config.Session.SendQueueSize = 500
	}
	if config.Session.MaxSendErrors == 0 {
		config.Session.MaxSendErrors = 10
	}

	// 限流配置默认值
	if config.RateLimit.RequestsPerSecond == 0 {
		config.RateLimit.RequestsPerSecond = 1000
	}
	if config.RateLimit.BurstSize == 0 {
		config.RateLimit.BurstSize = 2000
	}
	if config.RateLimit.MaxConnections == 0 {
		config.RateLimit.MaxConnections = 2000
	}

	// VAD配置默认值
	if config.VAD.PoolSize == 0 {
		config.VAD.PoolSize = 200
	}
	if config.VAD.Threshold == 0 {
		config.VAD.Threshold = 0.5
	}

	// 日志配置默认值
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "text"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "console" // 默认只输出到控制台，避免文件路径问题
	}
	if config.Logging.FilePath == "" && (config.Logging.Output == "file" || config.Logging.Output == "both") {
		config.Logging.FilePath = "logs/speech.log" // 默认相对路径
	}
}

// validateUnifiedConfig 验证统一配置
func validateUnifiedConfig(config *UnifiedConfig) error {
	// 验证模式
	if config.Mode != "unified" && config.Mode != "separated" {
		return fmt.Errorf("invalid mode: %s, must be 'unified' or 'separated'", config.Mode)
	}

	// 验证STT配置
	if config.STT != nil {
		if config.STT.ModelPath == "" {
			return fmt.Errorf("stt.model_path is required")
		}
		if config.STT.TokensPath == "" {
			return fmt.Errorf("stt.tokens_path is required")
		}

		// 检查模型文件是否存在
		if _, err := os.Stat(config.STT.ModelPath); os.IsNotExist(err) {
			return fmt.Errorf("stt model file not found: %s", config.STT.ModelPath)
		}
		if _, err := os.Stat(config.STT.TokensPath); os.IsNotExist(err) {
			return fmt.Errorf("stt tokens file not found: %s", config.STT.TokensPath)
		}

		// 验证Provider
		if config.STT.Provider.Provider != "cpu" &&
			config.STT.Provider.Provider != "cuda" &&
			config.STT.Provider.Provider != "auto" {
			return fmt.Errorf("invalid stt provider: %s, must be cpu, cuda, or auto", config.STT.Provider.Provider)
		}
	}

	// 验证TTS配置
	if config.TTS != nil {
		if config.TTS.ModelPath == "" {
			return fmt.Errorf("tts.model_path is required")
		}

		// 检查模型文件是否存在
		if _, err := os.Stat(config.TTS.ModelPath); os.IsNotExist(err) {
			return fmt.Errorf("tts model file not found: %s", config.TTS.ModelPath)
		}

		// 验证Provider
		if config.TTS.Provider.Provider != "cpu" &&
			config.TTS.Provider.Provider != "cuda" &&
			config.TTS.Provider.Provider != "auto" {
			return fmt.Errorf("invalid tts provider: %s, must be cpu, cuda, or auto", config.TTS.Provider.Provider)
		}
	}

	// 统一模式必须同时配置STT和TTS
	if config.Mode == "unified" {
		if config.STT == nil {
			return fmt.Errorf("stt config is required in unified mode")
		}
		if config.TTS == nil {
			return fmt.Errorf("tts config is required in unified mode")
		}
	}

	return nil
}

