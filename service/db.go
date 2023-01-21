package service

import (
	"context"
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

func (s *Store) CreateQrCode(ctx context.Context, payload string) (*Response, error) {
	qr, err := GenerateQrCode(ctx, payload)
	if err != nil {
		return nil, err
	}

	png, err := qr.Code.PNG(250)
	if err != nil {
		return nil, err
	}

	if err = s.storeResource(qr.Id, png); err != nil {
		log.Println("failed to store new qr-code", err)
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}, nil
	}

	return &Response{
		StatusCode: http.StatusCreated,
		Data:       png,
	}, nil
}

func (s *Store) GetQrCode(ctx context.Context, id string) (*Response, error) {
	log.Printf("retrieving qr-code with id: %s\n", id)
	var res *Response
	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)

		res = &Response{
			StatusCode: http.StatusOK,
			Data:       val,
		}
		return nil
	}); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		}, nil
	}
	return res, nil
}

func (s *Store) storeResource(key string, value []byte) error {
	txn := s.db.NewTransaction(true) // read-write txn

	qrCodeEntry := badger.NewEntry([]byte(key), value)
	if err := txn.SetEntry(qrCodeEntry); err != nil {
		log.Println("error storing qr-code", err)
		return err
	}
	return txn.Commit()
}

func logErr(err error) *Response {
	log.Println("failed to genereate qr-code", err)
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Data:       nil,
	}
}
