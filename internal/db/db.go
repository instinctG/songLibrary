package db

import (
	"context"
	"fmt"
	"github.com/instinctG/songLibrary/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Database struct {
	Client *pgxpool.Pool
	Log    *slog.Logger
}

func NewDatabase(dbSetting *config.Database, log *slog.Logger) (*Database, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbSetting.Host, dbSetting.Port, dbSetting.User, dbSetting.Name, dbSetting.Password, dbSetting.SSLMode)
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{Client: pool, Log: log}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx)
}
