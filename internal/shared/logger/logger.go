package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, message string, fields map[string]interface{})
	Info(ctx context.Context, message string, fields map[string]interface{})
	Warn(ctx context.Context, message string, fields map[string]interface{})
	Error(ctx context.Context, message string, err error, fields map[string]interface{})
}
