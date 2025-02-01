package setting

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`
	Database Database
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
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	return &cfg
}
