package service

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type Server struct {
	repository qrCodeStore
}
// QRCodeStore implements API calls to kv store
type qrCodeStore interface {
	GetQrCode(api map[string]QRCode, act <- chan Action)
	CreateQrCode(api map[string]QRCode, act <- chan Action)
}


func NewService(repo qrCodeStore) *Server {
	return &Server{
		repository: repo,
	}
}

type reqPayload struct {
	id string
	site_id string
	resource string
}

type response struct {
	StatusCode int
	QrCodes []byte
}

type Action struct {
	Id string
	Type string
	Payload reqPayload
	RetChan chan<- response 
}


func (s *Server) StartServiceManager(api map[string]QRCode, action <- chan Action) {
	for {
		select {
		case act := <-action:
			switch act.Type {
			case "GET":
				s.repository.GetQrCode(api, action)
			case "POST":
				id, _ := uuid.NewUUID() // TODO: 
				fmt.Println(id)
				s.repository.CreateQrCode(api, action)
			}
		}
	}
}


func SeviceHandler(w http.ResponseWriter, r *http.Request, method string, action chan<- Action) {
	resp := make(chan response)
	act := Action{
		Type: method,
		RetChan: resp,
	}

	if method == "POST" {
		var payload reqPayload
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err := json.Unmarshal(body, payload); err != nil {
			fmt.Println("Error on POST request")
			return
		}

		act.Payload = payload
	}
}
