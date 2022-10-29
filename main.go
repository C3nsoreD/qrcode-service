package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Port string
}

func main() {
	cfg := Config{
		Port: ":8080",
	}
	if err := initServer(cfg); err != nil {
		fmt.Printf("Failed to initialize server: %v", err)
	}
}

func initServer(cfg Config) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		return err
	}
	return nil
}
