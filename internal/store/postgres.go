package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
	"github.com/redis/go-redis/v9"
)

type LinkRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewLinkRepository(db *pgxpool.Pool, rdb *redis.Client) *LinkRepository {
	return &LinkRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *LinkRepository) Save(ctx context.Context, code string, url string) error {
	query := `insert into public.links (code, original_url, created_at) values ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, code, url, time.Now())
	if err != nil {
		return fmt.Errorf("postgres save link error: %w", err)
	}
	return nil
}

func (r *LinkRepository) Get(ctx context.Context, code string) (string, error) {
	query := `select original_url from public.links where code = $1`

	var originalURL string
	err := r.db.QueryRow(ctx, query, code).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperror.ErrNotFound
		}
		return "", err
	}
	return originalURL, nil
}

func (r *LinkRepository) IncrementClicks(ctx context.Context, code string) error {
	query := `UPDATE public.links SET clicks = clicks + 1 WHERE code = $1`
	_, err := r.db.Exec(ctx, query, code)
	if err != nil {
		return fmt.Errorf("postgres increment clicks error: %w", err)
	}
	return nil
}

func (r *LinkRepository) GetCache(ctx context.Context, code string) (string, error) {
	val, err := r.rdb.Get(ctx, code).Result()
	if err != nil {
		if err == redis.Nil {
			return "", apperror.ErrNotFound
		}
		return "", fmt.Errorf("redis get cache error: %w", err)
	}
	return val, nil
}

func (r *LinkRepository) SetCache(ctx context.Context, code string, url string, ttl time.Duration) error {
	err := r.rdb.Set(ctx, code, url, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set cache error: %w", err)
	}
	return nil
}

func (r *LinkRepository) DeleteCache(ctx context.Context, code string) error {
	err := r.rdb.Del(ctx, code).Err()
	if err != nil {
		return fmt.Errorf("redis del cache error: %w", err)
	}
	return nil
}
