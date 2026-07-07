package app

import (
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/config"
	transHTTP "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/http"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/store"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load configuration: " + err.Error())
	}
	
	log := logger.Init(cfg.Env)
	// defer log.Sync()	// ДОДЕЛАТЬ!!!
	defer func() { _ = log.Sync() }()

	log.Info("Configuration loaded and validated successfully", zap.String("env", cfg.Env))

	dbStore := store.NewMemoryStore()
	linkService := service.NewLinkService(dbStore)

	server := transHTTP.NewServer(":"+cfg.Port, log, linkService)

	if err := server.Start(); err != nil {
		log.Fatal("HTTP server failed to start", zap.Error(err))
	}
}