package service

import (
	// "fmt"
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
		logErr(act, err)
	}

	png, err := qr.Code.PNG(250)
	if err != nil {
		logErr(act, err)
	}

	if err = s.storeResource(qr.Id, png); err != nil {
		log.Println("failed to store new qr-code", err)
		act.RetChan <- Response{
			StatusCode: http.StatusInternalServerError,
			QrData:     nil,
		}
	}

	act.RetChan <- Response{
		StatusCode: http.StatusCreated,
	}
}

func (s *Store) GetQrCode(api map[string]QRCode, act Action) {
	log.Printf("retrieving qr-code with id: %s\n", act.Id)

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(act.Id))
		if err != nil {
			return err
		}
		val, _ := item.ValueCopy(nil)
		act.RetChan <- Response{
			StatusCode: http.StatusOK,
			QrData:     val,
		}
		return nil
	}); err != nil {
		act.RetChan <- Response{
			StatusCode: http.StatusInternalServerError,
			QrData:     nil,
		}
	}
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

func logErr(act Action, err error) {
	log.Println("failed to genereate qr-code", err)
	act.RetChan <- Response{
		StatusCode: http.StatusInternalServerError,
		QrData:     nil,
	}
}
