package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string        `mapstructure:"APP_PORT" validate:"required,numeric"`
	Env         string        `mapstructure:"APP_ENV" validate:"required,oneof=dev stage prod"`
	DatabaseURL string        `mapstructure:"APP_DATABASE_URL" validate:"required"`
	BaseURL     string        `mapstructure:"APP_BASE_URL" validate:"required"`
	RedisURL    string        `mapstructure:"APP_REDIS_URL" validate:"required"`
	CacheTTL    time.Duration `mapstructure:"APP_CACHE_TTL" validate:"required"` // Добавили TTL
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./configs/.env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("APP_CACHE_TTL", 10*time.Minute)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundErr) || os.IsNotExist(err) {
			log.Println("Local .env file not found. Loading configuration from system environment...")
		} else {
			return nil, fmt.Errorf("critical error reading config file: %w", err)
		}
	}

	_ = viper.BindEnv("APP_PORT")
	_ = viper.BindEnv("APP_ENV")
	_ = viper.BindEnv("APP_DATABASE_URL")
	_ = viper.BindEnv("APP_BASE_URL")
	_ = viper.BindEnv("APP_REDIS_URL")
	_ = viper.BindEnv("APP_CACHE_TTL")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("error validate config: %w", err)
	}

	return &cfg, nil
}
