package model

import (
	"github.com/skip2/go-qrcode"
)

type QRCode struct {
	Id   string
	Data *qrcode.QRCode 
}
