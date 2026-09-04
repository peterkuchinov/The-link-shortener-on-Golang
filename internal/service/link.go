package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var validCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

//go:generate mockgen -source=link.go -destination=mocks/mock_store.go -package=mocks

type LinkStore interface {
	Save(ctx context.Context, code string, url string) error
	Get(ctx context.Context, code string) (string, error)
	IncrementClicks(ctx context.Context, code string) error
}

type LinkCache interface {
	Get(ctx context.Context, code string) (string, error)
	Set(ctx context.Context, code string, url string, ttl time.Duration) error
	Delete(ctx context.Context, code string) error
}

type LinkService struct {
	store      LinkStore
	cache      LinkCache
	cacheTTL   time.Duration
	clickQueue chan ClickJob
}

func NewLinkService(store LinkStore, cache LinkCache, cacheTTL time.Duration) *LinkService {
	workerCount := 5
	queueBuffer := 10000

	s := &LinkService{
		store:      store,
		cache:      cache,
		cacheTTL:   cacheTTL,
		clickQueue: make(chan ClickJob, queueBuffer),
	}

	for i := 0; i < workerCount; i++ {
		go s.clickWorker()
	}

	return s
}

func (s *LinkService) Shorten(ctx context.Context, url, customCode string) (string, error) {
	code := customCode

	if code != "" {
		if !validCodeRegex.MatchString(code) {
			return "", apperror.ErrInvalidCustomCode
		}

		_, err := s.store.Get(ctx, code)
		if err != nil {
			if !errors.Is(err, apperror.ErrNotFound) {
				return "", fmt.Errorf("service failed to check existing code: %w", err)
			}
		} else {
			return "", apperror.ErrCodeAlreadyExists
		}
	} else {
		var err error
		code, err = utils.GenerateRandomCode(6)
		if err != nil {
			return "", fmt.Errorf("failed to generate random code: %w", err)
		}
	}

	if err := s.store.Save(ctx, code, url); err != nil {
		return "", fmt.Errorf("failed to save link: %w", err)
	}

	_ = s.cache.Delete(ctx, code)

	return code, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	log := logger.FromContext(ctx)

	cachedURL, err := s.cache.Get(ctx, code)
	if err == nil && cachedURL != "" {
		s.TrackClickAsync(ctx, code)
		return cachedURL, nil
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		log.Warn("cache failure, bypassing to database", zap.Error(err))
	}

	url, err := s.store.Get(ctx, code)
	if err != nil {
		return "", err
	}

	if setErr := s.cache.Set(ctx, code, url, s.cacheTTL); setErr != nil {
		log.Warn("failed to write to cache", zap.Error(setErr))
	}

	s.TrackClickAsync(ctx, code)

	return url, nil
}
