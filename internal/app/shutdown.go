package app

import (
	"context"

	"go.uber.org/zap"
)

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown HTTP server", zap.Error(err))
		return err
	}

	a.logger.Info("Closing database")
	a.db.Close()

	_ = a.logger.Sync()

	return nil
}
