package cora

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type cora struct {
	req     *http.Request
	client  *http.Client
	host    string
	headers []Header
}

type Header struct {
	Name   string
	Values []string
}

type Cora interface {
	Get(string, ...Header) Response
	Head(string, ...Header) Response
	Post(string, []byte, ...Header) Response
	Put(string, []byte, ...Header) Response
	Delete(string, ...Header) Response
	Options(string, ...Header) Response
	Trace(string, ...Header)
	Patch(string, []byte, ...Header) Response
	SetHost(string) Cora
	SetHeaders(...Header) Cora
}

// New cora instance
func New() Cora {
	return &cora{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		req: &http.Request{},
	}
}

// Request custom configuration
func (c *cora) Request(req *http.Request) {
	c.req = req
}

// Client custom configuration
func (c *cora) Client(client *http.Client) {
	c.client = client
}

// SetHost custom for a single request, apply it to all requests using Request method
func (c *cora) SetHost(host string) Cora {
	c.host = host
	return c
}

// SetHeaders replace headers for a single request, apply it to all requests using Request method
func (c *cora) SetHeaders(headers ...Header) Cora {
	c.headers = headers
	return c
}

// Get request
func (c *cora) Get(u string, h ...Header) Response {
	return c.make(http.MethodGet, u, nil, h...)
}

// Head request
func (c *cora) Head(u string, h ...Header) Response {
	return c.make(http.MethodHead, u, nil, h...)
}

// Post request
func (c *cora) Post(u string, b []byte, h ...Header) Response {
	return c.make(http.MethodPost, u, b, h...)
}

// Put request
func (c *cora) Put(u string, b []byte, h ...Header) Response {
	return c.make(http.MethodPut, u, b, h...)
}

// Delete request
func (c *cora) Delete(u string, h ...Header) Response {
	return c.make(http.MethodDelete, u, nil, h...)
}

// Options request
func (c *cora) Options(u string, h ...Header) Response {
	return c.make(http.MethodOptions, u, nil, h...)
}

// Trace request
func (c *cora) Trace(u string, h ...Header) {}

// Patch request
func (c *cora) Patch(u string, b []byte, h ...Header) Response {
	return c.make(http.MethodPatch, u, b, h...)
}

// headers used by a single request
func (c *cora) setHeaders(headers []Header, req *http.Request) {
	for _, v := range headers {
		for _, h := range v.Values {
			req.Header.Add(v.Name, h)
		}
	}
}

// make manage the http request
func (c *cora) make(method string, url string, b []byte, headers ...Header) Response {
	// check body content
	var body io.Reader
	if len(b) > 0 {
		body = bytes.NewBuffer(b)
		contentType := http.DetectContentType(b)
		log.Println(contentType)
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
		c.setHeaders(c.headers, req)
	}

	// set given headers for the request
	if len(headers) > 0 {
		c.setHeaders(headers, req)
	}

	log.Println("Type", req.Header.Values("Content-Type"))

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
