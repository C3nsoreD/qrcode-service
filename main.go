package main

import (
	"fmt"
	"log"
	"net/http"

	badger "github.com/dgraph-io/badger/v3"

	"github.com/c3nsored/qrcode-service/service"
)

type Config struct {
	Addr string
}

const localStore = ".store"

func main() {
	cfg := Config{
		Addr: "127.0.0.1:3000",
	}
	qrCodesData := make(map[string][]byte)

	// initialized kv data store
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := service.NewStore(db, qrCodesData)

	apiHandler := service.NewService(store)

	if err := initServer(cfg, apiHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServer(cfg Config, handlers http.Handler) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "Welcome to home page")
	})
	mux.Handle("/api/qrcode/", handlers)
	return http.ListenAndServe(cfg.Addr, mux)
}
