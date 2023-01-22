package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
	Id        string `json:"id,omitempty"`
	UrlString string `json:"url,omitempty"`
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
		if err := s.handleGet(ctx, rw, req); err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
	case http.MethodPost:
		if err := s.handlePost(ctx, rw, req); err != nil {
			WriteError(rw, http.StatusInternalServerError, err)
		}
	}
}

func (s *Server) handleGet(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	id, err := extractId(req.URL.Path)
	if err != nil {
		return errors.Wrapf(err, "service failed to extract id")
	}
	resp, err := s.Repo.GetQrCode(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "service failed to get qrcode")
	}

	WriteResponse(rw, resp)
	return nil
}

func (s *Server) handlePost(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	payload, err := getPayload(req)
	if err != nil {
		return errors.Wrapf(err, "service failed")
	}

	resp, err := s.Repo.CreateQrCode(ctx, payload)
	if err != nil {
		return errors.Wrapf(err, "service failed create qrcode")
	}

	WriteResponse(rw, resp)
	return nil
}

func getPayload(req *http.Request) (string, error) {
	var payload rPayload
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	if ok, err := ValidateUrlString(payload.UrlString); !ok {
		return "", err
	}
	return payload.UrlString, nil
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	errMsg := Response{
		StatusCode: statusCode,
		Message:    err.Error(),
	}

	if serializedError, err := json.Marshal(errMsg); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error while serializing paylaod:", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(serializedError)
	}
	return
}

func WriteResponse(w http.ResponseWriter, resp *Response) {
	if resp == nil {
		return
	}

	_, err := json.Marshal(resp)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		fmt.Println("Error while serializing payload:", err)
	} else {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(resp.Data)))
		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Data)
	}
	return
}

// extract id from all links to api/qrcode/...
func extractId(path string) (string, error) {
	id := strings.Split(path[1:], "/")[2]
	if id == "" {
		return "", errors.New("no id provided")
	}
	return id, nil
}

// verifies that a url is valid
func ValidateUrlString(urlText string) (bool, error) {
	res, err := url.ParseRequestURI(urlText)
	if err != nil || res == nil {
		return false, errors.Wrapf(err, "failed validate url")
	}
	return true, nil
}
