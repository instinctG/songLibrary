package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`
	Database Database
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

type Database struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER"`
	Name     string `env:"DB_NAME"`
	Password string `env:"DB_PASSWORD"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

func MustLoad() *Config {
	var cfg Config

	if err := godotenv.Load("config/local.env"); err != nil {
		slog.Warn("config wasn't loaded")
		os.Exit(1)
	}

	if err := env.Parse(&cfg); err != nil {
		slog.Warn("cannot parse config: ", err)
		os.Exit(1)
	}

	return &cfg
}
