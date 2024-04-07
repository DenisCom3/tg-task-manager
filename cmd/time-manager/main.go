package main

import (
	"fmt"
	"log/slog"
	"time-manager/internal/config"
	"time-manager/internal/logging"

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
		fmt.Printf("Error loading .env file")
		return err
	}

	cfg := config.MustLoad()

	
	log, err := logging.Setup(cfg.Env)

	if err != nil {
		fmt.Printf("failed to setup logging: %v", err)
		return err
	}

	log.Info("Starting time-manager", slog.String("env", cfg.Env))

	return nil
}