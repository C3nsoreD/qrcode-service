package api

import (
	"fmt"
	"context"
	"net/http"

	"github.com/c3nsored/qrcode-service/pkg/model"
)

// QrCodeService defines the service API methods 
type QrCodeService interface {
	GetQRCode(ctx context.Context, id string) (*model.QRCode, error)
}

// App contains configuration details
type App interface {
	GetEnv() string
	GetVersion() string
}

// Server represents a complete server
type Server struct {
	srv QrCodeService
	app App
}

// Creates a new server exposing necessary Endpoints
func New(application App, srv QrCodeService) *Server {
	return &Server{
		srv: srv,
		app: application,
	}
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status: available\n")
	fmt.Fprintf(w, "envirnoment: %s\n", s.app.GetEnv())
	fmt.Fprintf(w, "version: %s\n", s.app.GetVersion())
}