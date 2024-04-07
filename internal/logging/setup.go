package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	envDev = "dev"
	envProd = "prod"

	logPath = "./storage/logs/go.log"
)

func Setup(env string) (*slog.Logger, error) {
	
	var log *slog.Logger

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		fmt.Printf("failed to open log file: %v", err)
		return nil, err
	}

	w := io.Writer(logFile)


	switch env {
	case envDev:
		log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log, nil
}