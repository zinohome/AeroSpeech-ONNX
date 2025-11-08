package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志记录器接口
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	SetLevel(level LogLevel)
	SetOutput(w io.Writer)
}

// Config 日志配置
type Config struct {
	Level      string
	Format     string // "text" or "json"
	Output     string // "console", "file", or "both"
	FilePath   string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// defaultLogger 默认日志记录器
type defaultLogger struct {
	level  LogLevel
	writer io.Writer
	format string
}

var (
	globalLogger Logger
	levelMap     = map[string]LogLevel{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
	}
)

// InitLogger 初始化日志系统
func InitLogger(config Config) error {
	level, ok := levelMap[config.Level]
	if !ok {
		level = InfoLevel
	}

	var writers []io.Writer

	// 控制台输出
	if config.Output == "console" || config.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if config.Output == "file" || config.Output == "both" {
		if config.FilePath != "" {
			// 确保目录存在
			dir := filepath.Dir(config.FilePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create log directory: %w", err)
			}

			fileWriter := &lumberjack.Logger{
				Filename:   config.FilePath,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Compress:   config.Compress,
			}
			writers = append(writers, fileWriter)
		}
	}

	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = io.MultiWriter(writers...)
	}

	globalLogger = &defaultLogger{
		level:  level,
		writer: writer,
		format: config.Format,
	}

	return nil
}

// GetLogger 获取全局日志记录器
func GetLogger() Logger {
	if globalLogger == nil {
		// 默认日志记录器
		globalLogger = &defaultLogger{
			level:  InfoLevel,
			writer: os.Stdout,
			format: "text",
		}
	}
	return globalLogger
}

// log 记录日志
func (l *defaultLogger) log(level LogLevel, message string) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := level.String()

	if l.format == "json" {
		// JSON格式日志
		logEntry := fmt.Sprintf(
			`{"timestamp":"%s","level":"%s","message":"%s"}`,
			timestamp, levelStr, escapeJSON(message),
		)
		fmt.Fprintln(l.writer, logEntry)
	} else {
		// 文本格式日志
		logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, levelStr, message)
		fmt.Fprintln(l.writer, logEntry)
	}
}

// escapeJSON 转义JSON字符串
func escapeJSON(s string) string {
	result := ""
	for _, r := range s {
		switch r {
		case '"':
			result += "\\\""
		case '\\':
			result += "\\\\"
		case '\n':
			result += "\\n"
		case '\r':
			result += "\\r"
		case '\t':
			result += "\\t"
		default:
			result += string(r)
		}
	}
	return result
}

// Debug 记录Debug级别日志
func (l *defaultLogger) Debug(args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprint(args...)
		l.log(DebugLevel, message)
	}
}

// Debugf 记录Debug级别日志（格式化）
func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(DebugLevel, message)
}

// Info 记录Info级别日志
func (l *defaultLogger) Info(args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprint(args...)
		l.log(InfoLevel, message)
	}
}

// Infof 记录Info级别日志（格式化）
func (l *defaultLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(InfoLevel, message)
}

// Warn 记录Warn级别日志
func (l *defaultLogger) Warn(args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprint(args...)
		l.log(WarnLevel, message)
	}
}

// Warnf 记录Warn级别日志（格式化）
func (l *defaultLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(WarnLevel, message)
}

// Error 记录Error级别日志
func (l *defaultLogger) Error(args ...interface{}) {
	if len(args) > 0 {
		message := fmt.Sprint(args...)
		l.log(ErrorLevel, message)
	}
}

// Errorf 记录Error级别日志（格式化）
func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(ErrorLevel, message)
}

// SetLevel 设置日志级别
func (l *defaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput 设置输出
func (l *defaultLogger) SetOutput(w io.Writer) {
	l.writer = w
}

// 便捷函数
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

