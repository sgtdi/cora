package cora

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func encode(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal("encode error: ", err)
	}
	return buf.Bytes()
}

// makeRequest manage the http request
func (c *cora) request(method string, url string, i interface{}, headers ...Header) Response {
	var b []byte
	var body io.Reader
	var contentType string

	// encode and compress body
	if i != nil {
		switch v := i.(type) {
		case []byte:
			b = encode(v)
			break
		default:
			res, err := json.Marshal(v)
			if err != nil {
				b = encode(v)
				break
			}
			b = res
			break
		}
	}

	// check body content type
	if len(b) > 0 && len(headers) <= 0 {
		body = bytes.NewBuffer(b)
		// Check content type
		var js json.RawMessage
		if json.Unmarshal(b, &js) == nil {
			contentType = "application/json"
		} else {
			contentType = http.DetectContentType(b)
		}
		headers = append(headers, Header{
			Name:   "Content-Type",
			Values: []string{contentType},
		})
	}

	// create new request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return Response{Err: err}
	}

	// set custom host
	if len(c.host) > 0 {
		req.Host = c.host
	}

	// set custom header
	if len(c.headers) > 0 {
		c.setHeaders(req, c.headers...)
	}

	// set given headers for the request
	if len(headers) > 0 {
		c.setHeaders(req, headers...)
	}

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{Err: err}
	}
	defer resp.Body.Close()

	// read body
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		res := Response{Body: body, Code: resp.StatusCode, Raw: resp, Err: err}
		// check status code
		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			res.Err = errors.New(http.StatusText(resp.StatusCode))
		}
		// response
		return res
	}

	// response
	return Response{Raw: resp, Code: resp.StatusCode}
}
