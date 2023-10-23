package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		PG
	}

	// App -.
	App struct {
		Name            string `env-required:"true" env:"APP_NAME"`
		Version         string `env-required:"true" env:"APP_VERSION"`
		FileStoragePath string `env-required:"true" env:"FILE_STORAGE_PATH"`
	}

	// HTTP -.
	HTTP struct {
		Port          string `env-required:"true" env:"HTTP_PORT"`
		Debug         bool   `env:"HTTP_DEBUG" default:"false"`
		EnableSwagger bool   `env:"HTTP_ENABLE_SWAGGER" default:"false"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" env:"PG_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Print("Error loading .env file")
	}

	cfg := &Config{}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
