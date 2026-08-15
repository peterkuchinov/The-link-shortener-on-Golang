package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

func Init(env string) *zap.Logger {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	log, err := config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	return log
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
