package cora

import (
	"crypto/tls"
	"golang.org/x/net/http2"
	"net"
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
	Post(string, interface{}, ...Header) Response
	Put(string, interface{}, ...Header) Response
	Delete(string, ...Header) Response
	Options(string, ...Header) Response
	Patch(string, interface{}, ...Header) Response
	SetHost(string) Cora
	SetHeaders(...Header) Cora
}

// Http cora instance
func Http() Cora {
	return &cora{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		req: &http.Request{},
	}
}

// Http2 cora instance
func Http2() Cora {
	return &cora{
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http2.Transport{
				// So http2.Transport doesn't complain the URL scheme isn't 'https'
				AllowHTTP: true,
				// Pretend we are dialing a TLS endpoint. (Note, we ignore the passed tls.Config)
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
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
	return c.request(http.MethodGet, u, nil, h...)
}

// Head request
func (c *cora) Head(u string, h ...Header) Response {
	return c.request(http.MethodHead, u, nil, h...)
}

// Post request
func (c *cora) Post(u string, i interface{}, h ...Header) Response {
	return c.request(http.MethodPost, u, i, h...)
}

// Put request
func (c *cora) Put(u string, i interface{}, h ...Header) Response {
	return c.request(http.MethodPut, u, i, h...)
}

// Delete request
func (c *cora) Delete(u string, h ...Header) Response {
	return c.request(http.MethodDelete, u, nil, h...)
}

// Options request
func (c *cora) Options(u string, h ...Header) Response {
	return c.request(http.MethodOptions, u, nil, h...)
}

// Patch request
func (c *cora) Patch(u string, i interface{}, h ...Header) Response {
	return c.request(http.MethodPatch, u, i, h...)
}

// setHeaders used by a single request
func (c *cora) setHeaders(req *http.Request, headers ...Header) {
	for _, v := range headers {
		for _, h := range v.Values {
			req.Header.Add(v.Name, h)
		}
	}
}
