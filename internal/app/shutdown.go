package app

import (
	"context"

	"go.uber.org/zap"
)

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down application...")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown failed", zap.Error(err))
	}

	a.db.Close()

	if err := a.redisStore.Close(); err != nil {
		a.logger.Error("redis close failed", zap.Error(err))
	}

	return nil
}
