package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	transHTTP "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/http"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/store"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/utils/config"

	"go.uber.org/zap"
)

type App struct {
	server *transHTTP.Server
	db     *pgxpool.Pool
	rdb    *redis.Client
	logger *zap.Logger
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

	appLogger.Info(
		"Configuration loaded successfully",
		zap.String("env", cfg.Env),
	)

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

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		dbPool.Close()
		log.Fatalf("failed to parse redis url: %v", err)
	}
	rdb := redis.NewClient(redisOpts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		dbPool.Close()
		log.Fatalf("failed to ping redis: %v", err)
	}

	linkRepo := store.NewLinkRepository(dbPool, rdb)

	linkService := service.NewLinkService(linkRepo)

	server := transHTTP.NewServer(
		":"+cfg.Port,
		cfg.BaseURL,
		appLogger,
		linkService,
	)

	return &App{
		server: server,
		db:     dbPool,
		rdb:    rdb,
		logger: appLogger,
	}, nil
}
