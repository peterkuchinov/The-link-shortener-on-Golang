package service

import (
	"context"
	"log"
)

type ClickJob struct {
	Code string
}

func (s *LinkService) clickWorker() {
	ctx := context.Background() 
	for job := range s.clickQueue {
		err := s.store.IncrementClicks(ctx, job.Code)
		if err != nil {
			log.Printf("click save error for %s: %v", job.Code, err)
		}
	}
}

func (s *LinkService) TrackClickAsync(code string) {
	select {
	case s.clickQueue <- ClickJob{Code: code}:
	default:
		log.Printf("queue out of range, click miss for %s", code)
	}
}
