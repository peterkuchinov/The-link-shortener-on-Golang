package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"APP_PORT" validate:"required,numeric"`
	Env    string `mapstructure:"APP_ENV" validate:"required,oneof=dev stage prod"`
	AppKey string `mapstructure:"APP_KEY" validate:"required,min=32"`
	DatabaseURL string `mapstructure:"APP_DATABASE_URL" validate:"required"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./configs/.env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error read file config: %w", err)
	}

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