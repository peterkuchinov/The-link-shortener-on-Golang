package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"APP_PORT" validate:"required,numeric"`
	Env         string `mapstructure:"APP_ENV" validate:"required,oneof=dev stage prod"`
	DatabaseURL string `mapstructure:"APP_DATABASE_URL" validate:"required"`
	BaseURL     string `mapstructure:"APP_BASE_URL" validate:"required"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./configs/.env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error read file config: %w", err)
		}
	}

	_ = viper.BindEnv("APP_PORT")
	_ = viper.BindEnv("APP_ENV")
	_ = viper.BindEnv("APP_KEY")
	_ = viper.BindEnv("APP_DATABASE_URL")
	_ = viper.BindEnv("APP_REDIS_URL")
	_ = viper.BindEnv("APP_BASE_URL")

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
