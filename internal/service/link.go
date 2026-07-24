package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/utils"
)

var validCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

//go:generate mockgen -source=link.go -destination=mocks/mock_store.go -package=mocks
type LinkStore interface {
	Save(ctx context.Context, code string, url string) error
	Get(ctx context.Context, code string) (string, error)
}

type LinkService struct {
	store LinkStore
}

func NewLinkService(store LinkStore) *LinkService {
	return &LinkService{store: store}
}

func (s *LinkService) Shorten(ctx context.Context, url, customCode string) (string, error) {
    code := customCode

    if code != "" {
        if !validCodeRegex.MatchString(code) {
            return "", apperror.ErrInvalidCustomCode
        }
		
        existing, err := s.store.Get(ctx, code)
        if err != nil {
            return "", fmt.Errorf("service failed to check existing code: %w", err)
        }

        if existing != "" {
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