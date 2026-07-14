package app

import (
	"context"
	"time"

	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	transHTTP "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/http"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/store"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/util/config"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load configuration: " + err.Error())
	}
	
	log := logger.Init(cfg.Env)

	log.Info("Configuration loaded and validated successfully", zap.String("env", cfg.Env))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()
	
	dbStore := store.NewMemoryStore()
	linkService := service.NewLinkService(dbStore)

	server := transHTTP.NewServer(":"+cfg.Port, log, linkService)

	if err := server.Start(); err != nil {
		log.Fatal("HTTP server failed to start", zap.Error(err))
	}
}