package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

func getPayload(req *http.Request) (string, error) {
	var payload Payload
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
		log.Println("Error while serializing paylaod:", err)
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
		log.Println("Error while serializing payload:", err)
	} else {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(resp.Data)))
		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Data)
	}
	return
}
