package service

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"go.uber.org/zap"
)

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
	const alph = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	res := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alph))))
		if err != nil {
			return "", err
		}
		res[i] = alph[num.Int64()]
	}
	return string(res), nil
}

func (s *LinkService) Shorten(ctx context.Context, url string) (string, error) {
	log := logger.FromContext(ctx)
	
	log.Info("shortening url", zap.String("url", url))

	code, err := generateRandomCode(6)
	if err != nil {
		log.Error("failed to generate code", zap.Error(err))
		return "", err
	}

	err = s.store.Save(ctx, url, code)
	if err != nil {
		log.Error("failed to save link to store", zap.String("code", code), zap.Error(err))
		return "", err
	}

	log.Info("url successfully shortened", zap.String("code", code))
	return code, nil
}
