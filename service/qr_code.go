package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"

	"github.com/skip2/go-qrcode"
)

type QrCode struct {
	Id   string // unique id generated from url hashcode
	Code *qrcode.QRCode
	url  string
}

func GenerateQrCode(ctx context.Context, url string) (*QrCode, error) {
	id := md5.Sum([]byte(url))

	q, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate qrcode")
	}

	return &QrCode{
		Id:   fmt.Sprintf("%x", id),
		Code: q,
		url:  url,
	}, nil
}
