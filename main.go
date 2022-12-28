package main

import (
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"net/http"

	"github.com/c3nsored/qrcode-service/service"
)

type Config struct {
	Addr string
}

const localStore = ".store"

func main() {
	cfg := Config{
		Addr: "127.0.0.1:8080",
	}
	qrCodesData := make(map[string][]byte)

	// initialized kv data store
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := service.NewStore(db, qrCodesData)

	srv := service.NewService(store)

	if err := initServer(cfg, srv); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServer(cfg Config, handlers http.Handler) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Addr)

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: handlers,
	}

	return server.ListenAndServe()
}
