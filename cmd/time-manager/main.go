package main

import (
	"fmt"
	"log/slog"
	"time"
	"time-manager/internal/adapter/sqlite-repo/event"
	"time-manager/internal/config"
	"time-manager/internal/entity"
	"time-manager/internal/logging"
	"time-manager/internal/service"
	"time-manager/internal/storage/sqlite"

	"github.com/joho/godotenv"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}


func run() error {
	
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	cfg := config.MustLoad()

	
	log, err := logging.Setup(cfg.Env)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Info("Starting time-manager", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	eventEntity := entity.Event{
		ID: 1,
		Title: "demo1",
		Time:  time.Now(),
		
	}


	eService := service.NewEventService(eventEntity, event.NewRepo(storage))

	e, err := eService.GetByName()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(e.Title)
	return nil
}
