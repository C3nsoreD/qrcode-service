package service

import (
	"fmt"
	"log"
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
	db   *badger.DB
	data map[string][]byte
}

func NewStore(db *badger.DB, data map[string][]byte) *Store {
	return &Store{
		db:   db,
		data: data,
	}
}

func (s *Store) CreateQrCode(api map[string]QRCode, act Action) {

	qr, err := GenerateQrCode(act.Payload.Resource)
	if err != nil {
		log.Println(fmt.Println("Error creating Qrcode", err))
		act.RetChan <- Response{
			StatusCode: http.StatusInternalServerError,
			QrData:     nil,
		}
	}
	png, err := qr.Code.PNG(250)
	if err != nil {
		log.Println(fmt.Println("Error creating Qrcode", err))
		act.RetChan <- Response{
			StatusCode: http.StatusInternalServerError,
			QrData:     nil,
		}
	}
	log.Println(fmt.Println("ID:", qr.Id))
	s.data[qr.Id] = png
	act.RetChan <- Response{
		StatusCode: http.StatusCreated,
	}
}

func (s *Store) GetQrCode(api map[string]QRCode, act Action) {
	qr, ok := s.data[act.Id]
	log.Println(fmt.Printf("getting %s, %v\n", act.Id, ok))
	if !ok {
		act.RetChan <- Response{
			StatusCode: http.StatusInternalServerError,
			QrData:     nil,
		}
	}

	act.RetChan <- Response{
		StatusCode: http.StatusOK,
		QrData:     qr,
	}
}
