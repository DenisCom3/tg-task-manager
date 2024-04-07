package main

import (
	"fmt"
	"log"
	"time-manager/internal/config"

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
		log.Fatal("Error loading .env file")
		return err
	}

	cfg := config.MustLoad()

	fmt.Println(cfg.Env)

	return nil
}