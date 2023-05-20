package oci

import "context"

type Driver interface {
	Connector

	Handler(string) Handler
}

type Conn interface {
	Close() error
}

type Connector interface {
	Connect(ctx context.Context, uri string) (Conn, error)
}

type Handler interface {
	ServeOCI(r *Request) (*Response, error)
}

type HandlerFunc func(r *Request) (*Response, error)

func (h HandlerFunc) ServeOCI(r *Request) (*Response, error) {
	return h(r)
}
