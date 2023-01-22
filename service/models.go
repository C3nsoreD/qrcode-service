package service

// reqPayload contains basic params expected in a request
type Payload struct {
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
	Payload Payload
	RetChan chan<- Response
}
