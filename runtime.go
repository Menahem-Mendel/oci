package oci

import (
	"context"
	"errors"
	"oci/driver"
	"sync"
)

type Runtime struct {
	driver driver.Driver

	conns map[string]driver.Conn

	cancel func()

	mu sync.RWMutex
}

type runtime struct {
}

func (r runtime) Serve(p driver.Puller) {

}

func New(driver driver.Driver) (*Runtime, error) {
	// _, cancel := context.WithCancel(ctx)
	return &Runtime{
		// cancel: cancel,
	}, nil
}

func (r *Runtime) Open(uri string) (driver.Conn, error) {
	if r.driver == nil {
		return nil, errors.New("oci: runtime have nil driver")
	}
	return r.driver.Open(uri)
}

func (r *Runtime) Serve(ctx context.Context, h driver.Handler) error {
	return h.ServeOCI(ctx)
}

func (r *Runtime) Close() error {
	if r.cancel == nil {
		return nil
	}

	r.cancel()
	return nil
}
