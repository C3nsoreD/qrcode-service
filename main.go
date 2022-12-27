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

	// initialized database
	db, err := badger.Open(badger.DefaultOptions("/tmp/Badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := service.NewStore(db, qrCodesData)

	srv := service.NewService(store)

	actionCh := make(chan service.Action)      // actions channel
	qrCodes := make(map[string]service.QRCode) // used to mock data storage.

	go srv.ServiceManager(qrCodes, actionCh)
	qrCodeHandlers := MakeHandlers(service.SeviceHandler, "/api/qrcode/", actionCh)

	if err := initServer(cfg, qrCodeHandlers); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServer(cfg Config, handlers http.HandlerFunc) error {
	log.Printf("Starting qrcode-server on %s...", cfg.Addr)

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: handlers,
	}

	return server.ListenAndServe()
}

// Creates a generic handler for qr-code server
func MakeHandlers(
	fn func(http.ResponseWriter, *http.Request, string, string, chan<- service.Action),
	endpoint string,
	actionCh chan<- service.Action,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		log.Printf("Recieved request [%s] for path: [%s]", method, path)
		id := path[len(endpoint):]

		fn(w, r, id, method, actionCh)
	}

}
