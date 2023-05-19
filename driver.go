package oci

import "context"

type Driver interface {
	Connect(ctx context.Context, sock string) error

	Handler(string) Handler

	Close() error
}

type Conn interface {
	Close() error
}

type Handler interface {
	ServeOCI(r *Request) (*Response, error)
}

type HandlerFunc func(r *Request) (*Response, error)

func (h HandlerFunc) ServeOCI(r *Request) (*Response, error) {
	return h(r)
}
