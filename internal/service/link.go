package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
)

const alph = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var (
	ErrCodeAlreadyExists = errors.New("custom code is already taken")
	ErrInvalidCustomCode = errors.New("custom code contains invalid characters")
)
var validCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type LinkStore interface {
	Save(ctx context.Context, url string, code string) error
	Get(ctx context.Context, code string) (string, error)
}

type LinkService struct {
	store LinkStore
}

func NewLinkService(store LinkStore) *LinkService {
	return &LinkService{store: store}
}

func generateRandomCode(length int) (string, error) {
	res := make([]byte, length)
	n := int64(len(alph))
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(n))
		if err != nil {
			return "", err
		}
		res[i] = alph[num.Int64()]
	}
	return string(res), nil
}

func (s *LinkService) Shorten(ctx context.Context, url string, customCode string) (string, error) {
	code := customCode

	if code != "" {
		if !validCodeRegex.MatchString(code) {
			return "", ErrInvalidCustomCode
		}
		if existing, err := s.store.Get(ctx, code); err == nil && existing != "" {
			return "", ErrCodeAlreadyExists
		}
	} else {
		var err error
		if code, err = generateRandomCode(6); err != nil {
			return "", err
		}
	}
	
	if err := s.store.Save(ctx, code, url); err != nil {
		return "", err
	}

	return code, nil
}
