package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (a *App) Run() error {
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Fatal(
				"HTTP server failed",
				zap.Error(err),
			)
		}
	}()

	a.logger.Info("HTTP server started")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	a.logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	return a.Shutdown(shutdownCtx)
}
