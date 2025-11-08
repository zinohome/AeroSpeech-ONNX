package utils

import (
	"fmt"
	"os"
)

// ErrorCode 错误码
type ErrorCode string

const (
	ErrCodeInvalidParams      ErrorCode = "INVALID_PARAMS"
	ErrCodeAuthFailed         ErrorCode = "AUTH_FAILED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeAudioFormatError   ErrorCode = "AUDIO_FORMAT_ERROR"
	ErrCodeModelLoadError      ErrorCode = "MODEL_LOAD_ERROR"
	ErrCodeRecognitionError    ErrorCode = "RECOGNITION_ERROR"
	ErrCodeSynthesisError      ErrorCode = "SYNTHESIS_ERROR"
	ErrCodeResourceExhausted  ErrorCode = "RESOURCE_EXHAUSTED"
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode
	Message string
	Details string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s (%v)", e.Code, e.Message, e.Details, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建应用错误
func NewAppError(code ErrorCode, message, details string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// WrapError 包装错误
func WrapError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: err.Error(),
		Err:     err,
	}
}

// ReadFile 读取文件（工具函数）
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

