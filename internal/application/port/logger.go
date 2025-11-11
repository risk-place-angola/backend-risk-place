package port

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, fields ...map[string]interface{})
	Error(ctx context.Context, msg string, fields ...map[string]interface{})
	Debug(ctx context.Context, msg string, fields ...map[string]interface{})
	Warn(ctx context.Context, msg string, fields ...map[string]interface{})
}
