package main

import (
	"context"
	"github.com/instinctG/songLibrary/internal/db"
	"github.com/instinctG/songLibrary/internal/http-server/handler"
	"github.com/instinctG/songLibrary/internal/service"
	"github.com/instinctG/songLibrary/pkg/config"
	sl "github.com/instinctG/songLibrary/pkg/logger"
	"log/slog"
)

func Run() error {

	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.LogLevel)

	log.Info("starting song library service", slog.String("LOG-LEVEL", cfg.LogLevel))
	log.Debug("debug messages are enabled")

	database, err := db.NewDatabase(&cfg.Database, log)
	if err != nil {
		log.Error("failed to connect to database")
		return err
	}

	if err = database.Ping(context.Background()); err != nil {
		log.Error("Failed to ping db connection")
		return err
	}

	if err = database.MigrateDB(); err != nil {
		log.Error("Failed to migrate database")
		return err
	}

	songService := service.NewService(database)

	httpHandler := handler.NewHandler(songService, cfg.HttpPort, log)
	if err = httpHandler.Serve(); err != nil {
		log.Error("failed to start server")
		return err
	}

	return nil
}

// @title Song Library API
// @version 1.0
// @description Swagger API for Song Library
// @host localhost:8080
// @BasePath /
func main() {
	if err := Run(); err != nil {
		slog.Error("could not run the application", sl.Err(err))
	}

}
