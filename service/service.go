package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type Server struct {
	repository qrCodeStore
}

// QRCodeStore implements API calls to kv store
type qrCodeStore interface {
	GetQrCode(api map[string]QRCode, act Action)
	CreateQrCode(api map[string]QRCode, act Action)
}

func NewService(repo qrCodeStore) *Server {
	return &Server{
		repository: repo,
	}
}

type reqPayload struct {
	id       string
	site_id  string
	resource string
}

type Response struct {
	StatusCode int
	QrCodes    string
}

type Action struct {
	Id      string
	Type    string
	Payload reqPayload
	RetChan chan<- Response
}

func (s *Server) StartServiceManager(api map[string]QRCode, action <-chan Action) {
	for {
		select {
		case act := <-action:
			switch act.Type {
			case "GET":
				s.repository.GetQrCode(api, act)
			case "POST":
				id, _ := uuid.NewUUID() // TODO:
				fmt.Println(id)
				s.repository.CreateQrCode(api, act)
			}
		}
	}
}

func SeviceHandler(w http.ResponseWriter, r *http.Request, id, method string, action chan<- Action) {
	respCh := make(chan Response)
	act := Action{
		Id:      id,
		Type:    method,
		RetChan: respCh,
	}

	if method == "POST" {
		var payload reqPayload
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err := json.Unmarshal(body, &payload); err != nil {
			fmt.Println("Error on POST request")
			return
		}

		act.Payload = payload
	}
	action <- act
	var resp Response
	if resp = <-respCh; resp.StatusCode > http.StatusCreated {
		writeError(w, resp.StatusCode)
		return
	}

	writeResponse(w, resp)
}

func writeResponse(w http.ResponseWriter, resp Response) {
	serializedPayload, err := json.Marshal(resp.QrCodes)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		fmt.Println("Error while serializing payload:", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(serializedPayload)
	}
}

func writeError(w http.ResponseWriter, statusCode int) {
	jsonMsg := struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	}{
		Code: statusCode,
		Msg:  http.StatusText(statusCode),
	}
	if serializedPayload, err := json.Marshal(jsonMsg); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error while serializing paylaod:", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(serializedPayload)
	}
}
