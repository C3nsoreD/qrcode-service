package store

import (
	badger "github.com/dgraph-io/badger/v3"
	
	"context"
	"github.com/c3nsored/qrcode-service/pkg/model"
)

type store struct {
	db *badger.DB
}

func New(db *badger.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) GetQRCode(ctx context.Context, id string) (*model.QRCode, error) {
	return nil, nil
}

func (s *store) CreateQRCode(ctx context.Context, data []byte) (*model.QRCode, error) {
	return nil, nil
}