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
		return fmt.Errorf("%w", err)
	}

	cfg := config.MustLoad()

	
	log, err := logging.Setup(cfg.Env)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Info("Starting time-manager", slog.String("env", cfg.Env))

	return nil
}
