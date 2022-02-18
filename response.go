package cora

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Response struct {
	Code int
	Err  error
	Body []byte
	Raw  *http.Response
}

// JSON decode result
func (r Response) JSON(model interface{}) Response {
	if err := json.Unmarshal(r.Body, &model); err != nil {
		r.Err = err
	}
	return r
}

// XML decode result
func (r Response) XML(model interface{}) Response {
	if err := xml.Unmarshal(r.Body, &model); err != nil {
		r.Err = err
	}
	return r
}
