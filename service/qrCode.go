package service

import (
	"crypto/md5"
	"fmt"

	"github.com/skip2/go-qrcode"
)

type QrCode struct {
	Id   string // unique id generated from url hashcode
	Code *qrcode.QRCode
	url  string
}

func GenerateQrCode(url string) (*QrCode, error) {
	id := md5.Sum([]byte(url))

	q, err := qrcode.New(url, qrcode.Medium)

	if err != nil {
		return nil, err
	}

	return &QrCode{
		Id:   fmt.Sprintf("%x", id),
		Code: q,
		url:  url,
	}, nil
}
