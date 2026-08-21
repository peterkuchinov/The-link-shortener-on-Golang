package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/utils"
)

var validCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

//go:generate mockgen -source=link.go -destination=mocks/mock_store.go -package=mocks
type LinkStore interface {
	Save(ctx context.Context, code string, url string) error
	Get(ctx context.Context, code string) (string, error)
	IncrementClicks(ctx context.Context, code string) error
}

type LinkService struct {
	store      LinkStore
	clickQueue chan ClickJob
}

func NewLinkService(store LinkStore) *LinkService {
	workerCount := 5
	queueBuffer := 10000

	s := &LinkService{
		store:      store,
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
				return "", err
			}
		} else {
			return "", apperror.ErrCodeAlreadyExists
		}
	} else {
		var err error
		code, err = utils.GenerateRandomCode(6)
		if err != nil {
			return "", err
		}
	}

	if err := s.store.Save(ctx, code, url); err != nil {
		return "", err
	}

	return code, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	url, err := s.store.Get(ctx, code)
	if err != nil {
		return "", err
	}

	s.TrackClickAsync(ctx, code)

	return url, nil
}
