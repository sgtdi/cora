package cora

import (
	"io/ioutil"
	"net/http"
	"time"
)

type cora struct {
	req    *http.Request
	client *http.Client
}

type Cora interface {
	Get(string) Response
	Head(string) Response
	Post(string, []byte) Response
	Put(string, []byte) Response
	Delete(string) Response
	Options(string) Response
	Trace(string)
	Patch(string, []byte) Response
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
	c.req.Host = host
	return c
}

// SetHeaders replace headers for a single request, apply it to all requests using Request method
func (c *cora) SetHeaders(headers map[string]string) Cora {
	for k, v := range headers {
		c.req.Header.Add(k, v)
	}
	return c
}

// AddHeaders add additional headers for a single request, apply it to all requests using Request method
func (c *cora) AddHeaders(headers map[string]string) Cora {
	for k, v := range headers {
		c.req.Header.Set(k, v)
	}
	return c
}

// Get request
func (c *cora) Get(u string) Response {
	return c.make(http.MethodGet, u, nil)
}

// Head request
func (c *cora) Head(u string) Response {
	return c.make(http.MethodHead, u, nil)
}

// Post request
func (c *cora) Post(u string, b []byte) Response {
	return c.make(http.MethodPost, u, b)
}

// Put request
func (c *cora) Put(u string, b []byte) Response {
	return c.make(http.MethodPut, u, b)
}

// Delete request
func (c *cora) Delete(u string) Response {
	return c.make(http.MethodDelete, u, nil)
}

// Options request
func (c *cora) Options(u string) Response {
	return c.make(http.MethodOptions, u, nil)
}

// Trace request
func (c *cora) Trace(u string) {}

// Patch request
func (c *cora) Patch(u string, b []byte) Response {
	return c.make(http.MethodPatch, u, b)
}

// make manage the http request
func (c *cora) make(method string, url string, body []byte) Response {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return Response{Err: err}
	}

	// Set custom header
	if c.req.Header != nil {
		req.Header = c.req.Header
	}

	// Set custom host
	if len(c.req.Host) > 0 {
		req.Host = c.req.Host
	}

	// Do request
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{Err: err}
	}
	defer resp.Body.Close()

	// Check status code
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return Response{Err: err, Code: resp.StatusCode}
	}

	// Read body
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return Response{Body: body, Code: resp.StatusCode, Raw: resp, Err: err}
	}

	return Response{Raw: resp, Code: resp.StatusCode}
}
