package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Save(ctx context.Context, code string, url string) error {
	query := `insert into links (code, original_url, created_at) values ($1, $2, $3)`
	
	_, err := r.db.Exec(ctx, query, code, url, time.Now())
	if err != nil {
		return fmt.Errorf("postgres save link error: %w", err)
	}
	return nil
}

func (r *LinkRepository) GetByCode(ctx context.Context, code string) (string, error) {
	query := `select original_url from links where code = $1`
	
	var originalURL string
	err := r.db.QueryRow(ctx, query, code).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("error found: %w", err)
		}
		return "", fmt.Errorf("postgres get link error: %w", err)
	}
	return originalURL, nil
}
