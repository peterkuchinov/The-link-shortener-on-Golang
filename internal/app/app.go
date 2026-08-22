package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	transHTTP "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/http"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/store"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/utils/config"

	"go.uber.org/zap"
)

type App struct {
	server     *transHTTP.Server
	db         *pgxpool.Pool
	redisStore *store.RedisStore
	logger     *zap.Logger
}

func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	appLogger, err := logger.Init(cfg.Env)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	redisStore, err := store.NewRedisStore(cfg.RedisURL)
	if err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	linkRepo := store.NewLinkRepository(dbPool)
	linkService := service.NewLinkService(linkRepo, redisStore, cfg.CacheTTL)

	server := transHTTP.NewServer(
		":"+cfg.Port,
		cfg.BaseURL,
		appLogger,
		linkService,
	)

	return &App{
		server:     server,
		db:         dbPool,
		redisStore: redisStore,
		logger:     appLogger,
	}, nil
}
