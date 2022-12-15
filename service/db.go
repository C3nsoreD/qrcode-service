package service

import (
	// "context"
	"fmt"
	"net/http"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/skip2/go-qrcode"
)

type QRCode struct {
	Id     string
	SiteId string
	Data   []byte
}

const (
	QrCodeSize  = 256
	QrCodeLevel = qrcode.Medium
)

type Store struct {
	db *badger.DB
}

func NewStore(db *badger.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateQrCode(api map[string]QRCode, act Action) {
	// act.RetChan <- response{
	// 	StatusCode: http.StatusCreated,
	// }
	fmt.Println("Created QRCode")
}

func (s *Store) GetQrCode(api map[string]QRCode, act Action) {
	act.RetChan <- Response{
		StatusCode: http.StatusOK,
		QrCodes:    fmt.Sprintf("QrCodeUrl %s", act.Id),
	}
}

func encodeQrCode(content string) ([]byte, error) {
	fmt.Println("generating qrcode image...")
	png, err := qrcode.Encode(content, QrCodeLevel, QrCodeSize)
	if err != nil {
		return nil, err
	}
	return png, nil
}
