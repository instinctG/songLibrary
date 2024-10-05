package main

import (
	"fmt"
	"github.com/instinctG/songLibrary/internal/db"
	"github.com/instinctG/songLibrary/internal/handler"
	"github.com/instinctG/songLibrary/internal/service"
	"github.com/instinctG/songLibrary/pkg/setting"
	"log"
)

func init() {
	setting.Setup()
}

func Run() error {
	fmt.Println("starting our application")
	//todo : implement load configs
	dbSetting := setting.DatabaseSetting
	serverSetting := setting.ServerSetting

	database, err := db.NewDatabase(dbSetting)
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}

	if err = database.MigrateDB(); err != nil {
		fmt.Println("Failed to migrate database")
		return err
	}

	songService := service.NewService(database)
	httpHandler := handler.NewHandler(songService, serverSetting.Port)
	if err = httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}
func main() {
	log.Printf("[info] start https server listening ... ")
	if err := Run(); err != nil {
		log.Println(err)
	}
}
