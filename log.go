package amqp

import (
	"context"
	"fmt"
)

var logger ILogger

type ILogger interface {
	DebugCtxf(ctx context.Context, format string, args ...interface{})
	InfoCtxf(ctx context.Context, format string, args ...interface{})
	WarnCtxf(ctx context.Context, format string, args ...interface{})
	ErrorCtxf(ctx context.Context, format string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) DebugCtxf(ctx context.Context, format string, args ...interface{}) {
	println(fmt.Sprintf(format, args...))
}
func (l *defaultLogger) InfoCtxf(ctx context.Context, format string, args ...interface{}) {
	println(fmt.Sprintf(format, args...))
}
func (l *defaultLogger) WarnCtxf(ctx context.Context, format string, args ...interface{}) {
	println(fmt.Sprintf(format, args...))
}
func (l *defaultLogger) ErrorCtxf(ctx context.Context, format string, args ...interface{}) {
	println(fmt.Sprintf(format, args...))
}

func SetLogger(log ILogger) {
	logger = log
}

func GetLogger() ILogger {
	if logger == nil {
		logger = &defaultLogger{}
	}
	return logger
}
