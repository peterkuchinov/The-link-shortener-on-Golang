package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/apperror"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
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
