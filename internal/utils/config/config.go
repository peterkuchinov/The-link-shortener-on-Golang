package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"app_port" validate:"required,numeric"`
	Env         string `mapstructure:"app_env" validate:"required,oneof=dev stage prod"`
	DatabaseURL string `mapstructure:"app_database_url" validate:"required"`
	BaseURL     string `mapstructure:"app_base_url" validate:"required"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./configs/.env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundErr) || os.IsNotExist(err) {
			log.Println("Local .env file not found. Loading configuration from system environment...")
		} else {
			return nil, fmt.Errorf("critical error reading config file: %w", err)
		}
	}

	_ = viper.BindEnv("app_port", "APP_PORT")
	_ = viper.BindEnv("app_env", "APP_ENV")
	_ = viper.BindEnv("app_database_url", "APP_DATABASE_URL")
	_ = viper.BindEnv("app_base_url", "APP_BASE_URL")

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
