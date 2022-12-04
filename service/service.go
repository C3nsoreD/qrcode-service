package service

type Server struct {
	db QRCodeStore
}

// QRCodeStore implements API calls to kv store
type QRCodeStore interface {
	func GetQRCode(ctx context.Context, id string) (model.QrCode, error)
	func CreateQrCode(ctx context.Context, data interface) (model.QRCode, error)
}
