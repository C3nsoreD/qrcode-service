package service

import (
	"context"
	"github.com/c3nsored/qrcode-service/pkg/model"
)

type Service struct {
	db QRCodeStore
}


// QRCodeStore implements API calls to kv store
type QRCodeStore interface {
	GetQRCode(ctx context.Context, id string) (*model.QRCode, error)
	CreateQRCode(ctx context.Context, data []byte) (*model.QRCode, error)
}


func New(db QRCodeStore) *Service {
	return &Service{
		db: db,
	}
}


func (s *Service) GetQRCode(ctx context.Context, id string) (*model.QRCode, error) {
	return nil, nil
}

func (s *Service) CreateQRCode(ctx context.Context, data []byte) (*model.QRCode, error) {
	return nil, nil
}