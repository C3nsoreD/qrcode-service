package store

import (
	badger "github.com/dgraph-io/badger/v3"
	"log"
)

type store struct {
	db *badger.DB
}

func New(db *badger.DB) *store {
	return &store{
		db: db,
	}
}
