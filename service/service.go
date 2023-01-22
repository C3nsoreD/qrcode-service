package service

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type Server struct {
	Repo qrCodeStore
}

// QRCodeStore implements API calls to kv store
type qrCodeStore interface {
	GetQrCode(ctx context.Context, id string) (*Response, error)
	CreateQrCode(ctx context.Context, url string) (*Response, error)
}

func NewService(repo qrCodeStore) *Server {
	return &Server{
		Repo: repo,
	}
}

// path: api/qrcode
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		if err := s.handleGet(ctx, rw, req); err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
	case http.MethodPost:
		if err := s.handlePost(ctx, rw, req); err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
	}
}

func (s *Server) handleGet(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	id, err := extractId(req.URL.Path)
	if err != nil {
		return errors.Wrapf(err, "service failed to extract id")
	}
	resp, err := s.Repo.GetQrCode(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "service failed to get qrcode")
	}

	WriteResponse(rw, resp)
	return nil
}

func (s *Server) handlePost(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	payload, err := getPayload(req)
	if err != nil {
		return errors.Wrapf(err, "service failed")
	}

	resp, err := s.Repo.CreateQrCode(ctx, payload)
	if err != nil {
		return errors.Wrapf(err, "service failed create qrcode")
	}

	WriteResponse(rw, resp)
	return nil
}
