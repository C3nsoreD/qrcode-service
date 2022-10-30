package main

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"net/http"

	"github.com/c3nsored/qrcode-service/config"
)

func main() {
	cfg := config.Config{
		Port: ":8080",
	}

	// initialized database
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initServer(cfg); err != nil {
		fmt.Printf("Failed to initialize server: %v", err)
	}
}

func initServer(cfg config.Config) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		return err
	}
	return nil
}

func mustInitDatabase(cfg config.Config) {

}
