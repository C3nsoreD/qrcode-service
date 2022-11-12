package main

import (
	"flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c3nsored/qrcode-service/config"
	"github.com/c3nsored/qrcode-service/pkg/api"
	"github.com/c3nsored/qrcode-service/pkg/service"
	"github.com/c3nsored/qrcode-service/pkg/store"
	badger "github.com/dgraph-io/badger/v3"
)

const Version = "1.0.0"

type application struct {
	Config config.Config
	logger *log.Logger
}

func main() {
	var cfg config.Config

	// initialized database
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|stage|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	
	app := application{
		Config: cfg,
		logger: logger,
	}
	dbStore := store.New(db)
	svc := service.New(dbStore)
	server := api.New(&app, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", server.HealthCheckHandler)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
	
}

func (a *application) GetEnv() string {
	return a.Config.Env
}

func (a *application) GetVersion() string {
	return Version
}
