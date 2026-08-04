package app

import (
	"context"
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
	server *transHTTP.Server
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := logger.Init(cfg.Env)

	log.Info(
		"Configuration loaded successfully",
		zap.String("env", cfg.Env),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, err
	}

	repository := store.NewLinkRepository(dbPool)
	linkService := service.NewLinkService(repository)

	server := transHTTP.NewServer(
		":"+cfg.Port,
		cfg.BaseURL,
		log,
		linkService,
	)

	return &App{
		server: server,
		db:     dbPool,
		logger: log,
	}, nil
}