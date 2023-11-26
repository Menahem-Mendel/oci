package oci

import (
	"context"
	"oci/driver"
	"sync"
)

type Runtime struct {
	driver driver.Driver

	conns map[string]driver.Conn

	cancel func()

	mu sync.RWMutex
}

func New(driver driver.Driver) (*Runtime, error) {
	// _, cancel := context.WithCancel(ctx)
	return &Runtime{
		// cancel: cancel,
	}, nil
}

func (r *Runtime) Open(ctx context.Context, uri string) (driver.Conn, error) {
	return r.driver.Open(ctx, uri)
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
