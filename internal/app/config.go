package app

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"APP_PORT" validate:"required,numeric"`
	Env    string `mapstructure:"APP_ENV" validate:"required,oneof=dev stage prod"`
	AppKey string `mapstructure:"APP_KEY" validate:"required,min=32"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs") 
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения файла конфига: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка unmarshal конфига: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфига: %w", err)
	}

	return &cfg, nil
}
