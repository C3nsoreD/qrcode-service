package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	Repo qrCodeStore
}

// QRCodeStore implements API calls to kv store
type qrCodeStore interface {
	GetQrCode(ctx context.Context, id string) (*Response, error)
	CreateQrCode(ctx context.Context, url string) (*Response, error)
}

func NewService(repo qrCodeStore) *Server {
	return &Server{
		Repo: repo,
	}
}

// reqPayload contains basic params expected in a request
type rPayload struct {
	Id   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
}

// Response
type Response struct {
	StatusCode int
	Message    string // detailed information about the reponse.
	Data       []byte
}

type Action struct {
	Id      string
	Type    string
	Payload rPayload
	RetChan chan<- Response
}

// path: api/qrcode
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	switch req.Method {
	case http.MethodGet:
		id, err := extractId(req.URL.Path)
		if err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
			break
		}
		resp, err := s.Repo.GetQrCode(ctx, id)
		if err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
		WriteResponse(rw, resp)
	case http.MethodPost:
		payload, err := getPayload(req)
		if err != nil {
			return
		}
		resp, err := s.Repo.CreateQrCode(ctx, payload)
		if err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
		WriteResponse(rw, resp)
	}
}

func getPayload(req *http.Request) (string, error) {
	var payload rPayload
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error on POST request")
		return "", err
	}
	return payload.Text, nil
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	jsonMsg := struct {
		Msg   string `json:"msg"`
		Error string `json:"error"`
	}{
		Msg:   http.StatusText(statusCode),
		Error: err.Error(),
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

func WriteResponse(w http.ResponseWriter, resp *Response) {
	_, err := json.Marshal(resp.Data)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		fmt.Println("Error while serializing payload:", err)
	} else {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(resp.Data)))
		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Data)
	}
}

// extract id from all links to api/qrcode/...
func extractId(path string) (string, error) {
	id := strings.Split(path[1:], "/")[2]
	if id == "" {
		return "", fmt.Errorf("no id provided")
	}
	return id, nil
}
