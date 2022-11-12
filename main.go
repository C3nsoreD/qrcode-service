package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"flag"
	"time"
)

const version = "1.0.0"

type Config struct {
	Env string
	Port int32
}


type application struct {
	config Config
	logger *log.Logger
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|stage|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	
	app := application{
		config: cfg,
		logger: logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)


	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
	
}

