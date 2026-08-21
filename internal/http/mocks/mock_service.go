package mocks

import (
	"context"
)

type MockLinkServiceShortener struct {
	ShortenFunc        func(ctx context.Context, url, customCode string) (string, error)
	GetOriginalURLFunc func(ctx context.Context, code string) (string, error)
}

func (m *MockLinkServiceShortener) Shorten(ctx context.Context, url, customCode string) (string, error) {
	if m.ShortenFunc != nil {
		return m.ShortenFunc(ctx, url, customCode)
	}
	return "", nil
}

func (m *MockLinkServiceShortener) GetOriginalURL(ctx context.Context, code string) (string, error) {
	if m.GetOriginalURLFunc != nil {
		return m.GetOriginalURLFunc(ctx, code)
	}
	return "", nil
}
