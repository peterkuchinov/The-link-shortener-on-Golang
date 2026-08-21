package service

import (
	"context"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"go.uber.org/zap"
)

type ClickJob struct {
	Ctx  context.Context
	Code string
}

func (s *LinkService) clickWorker() {
	for job := range s.clickQueue {
		requestLogger := logger.FromContext(job.Ctx)
		bgCtx := context.Background()
		bgCtx = logger.ToContext(bgCtx, requestLogger)

		err := s.store.IncrementClicks(bgCtx, job.Code)
		if err != nil {
			requestLogger.Error("failed to increment link clicks in background",
				zap.String("code", job.Code),
				zap.Error(err),
			)
			continue
		}

		requestLogger.Info("link clicks successfully incremented asynchronously",
			zap.String("code", job.Code),
		)
	}
}

func (s *LinkService) TrackClickAsync(ctx context.Context, code string) {
	s.clickQueue <- ClickJob{
		Ctx:  ctx,
		Code: code,
	}
}
