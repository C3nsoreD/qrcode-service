package main

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"net/http"

	"github.com/c3nsored/qrcode-service/service"
)

type Config struct {
	Port string
}



func main() {
	cfg := Config{
		Port: ":8080",
	}

	// initialized database
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := service.NewStore(db)

	srv := service.NewService(store)

	// actions channel
	actionCh := make(chan service.Action)
	qrCodes := make(map[string]service.QRCode)

	go srv.StartServiceManager(qrCodes, actionCh)
	api := MakeHandlers(service.SeviceHandler, "api/qrcode/", actionCh)

	if err := initServer(cfg, api); err != nil {
		fmt.Printf("Failed to initialize server: %v", err)
	}
}

func initServer(cfg Config, handlers http.HandlerFunc) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, handlers); err != nil {
		return err
	}
	return nil
}

func mustInitDatabase(cfg Config) {

}
