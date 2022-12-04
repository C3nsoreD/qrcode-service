package service

import (
	"fmt"

	"github.com/c3nsored/qrcode-service/pkg/model"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/skip2/go-qrcode"
)

const (
	QrCodeSize = 256
	QrCodeLevel = qrcode.Medium
)

type store struct {
	db *badger.DB
}

func New(db *badger.DB) *store {
	return &store{
		db: db,
	}
}




func (s *store) CreateQrCode(ctx context.Context, data model.QRCode)


func (s *Server) GetQRCode(ctx context.Context, id string) (model.QrCode, error)


func encodeQrCode(content string) ([]byte, error) {
	fmt.Println("generating qrcode image...")
	png, err := qrcode.Encode(content, QrCodeLevel, QrCodeSize)
	if err != nil {
		return nil, err
	}
	return png, nil
}