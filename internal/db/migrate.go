package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (d *Database) MigrateDB() error {

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(d.Client)

	if err := goose.Up(db, "./migrations"); err != nil {
		if errors.Is(err, goose.ErrAlreadyApplied) {
			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	d.Log.Info("Successfully migrated DB")
	return nil
}
