package logger

import "go.uber.org/zap"

// NoOpLogger — заглушка для логгера
type NoOpLogger struct {
	logger *zap.Logger
}

func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{logger: zap.NewNop()}
}

func (n *NoOpLogger) Info(msg string, fields ...zap.Field)  {}
func (n *NoOpLogger) Error(msg string, fields ...zap.Field) {}
