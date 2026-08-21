package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

func Init(env string) (*zap.Logger, error) {
	var cfg zap.Config

	if env == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	return cfg.Build()
}

func ToContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return log
	}
	return zap.NewNop()
}
