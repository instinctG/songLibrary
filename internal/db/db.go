package db

import (
	"context"
	"fmt"
	"github.com/instinctG/songLibrary/pkg/setting"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Client *pgxpool.Pool
}

func NewDatabase(dbSetting *setting.Database) (*Database, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbSetting.Host, dbSetting.Port, dbSetting.User, dbSetting.Name, dbSetting.Password, dbSetting.SSLMode)
	fmt.Println(connString)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{Client: pool}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx)
}
