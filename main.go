package main

import (
	// "flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"net/http"
	// "os"
	// "time"

	"github.com/c3nsored/qrcode-service/service"
)

type Config struct {
	Port string
}

func main() {
	cfg := Config{
		Port: ":8080",
	}
	qrCodesData := make(map[string][]byte)

	// initialized database
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := service.NewStore(db, qrCodesData)

	srv := service.NewService(store)

	// actions channel
	actionCh := make(chan service.Action)
	qrCodes := make(map[string]service.QRCode)

	go srv.StartServiceManager(qrCodes, actionCh)
	api := MakeHandlers(service.SeviceHandler, "/api/qrcode/", actionCh)

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

func MakeHandlers(
	fn func(http.ResponseWriter, *http.Request, string, string, chan<- service.Action),
	endpoint string,
	actionCh chan<- service.Action,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		log.Println(fmt.Sprintf("Recieved request [%s] for path: [%s]", method, path))
		id := path[len(endpoint):]

		fn(w, r, id, method, actionCh)
	}

}
