package runtime

import (
	"context"
	"oci/driver"
	"sync"
)

type Runtime struct {
	driver driver.Connector

	conns map[string]driver.Conn

	cancel func()

	mu sync.RWMutex
}

func New(ctx context.Context, driver driver.Driver) (*Runtime, error) {
	_, cancel := context.WithCancel(ctx)
	return &Runtime{
		conns:  make(map[string]driver.Conn),
		cancel: cancel,
	}, nil
}

func (r *Runtime) Open(ctx context.Context, uri string) (driver.Conn, error) {
	return r.driver.Open(ctx, uri)
}

func (r *Runtime) ServeOCI(ctx context.Context, h driver.Handler) error {
	return h.ServeOCI(ctx)
}

func (r *Runtime) Close() error {
	if r.cancel == nil {
		return nil
	}

	r.cancel()
	return nil
}

type ChainConf struct {
	Ref string
}

func (r *Runtime) NewChain(conf ChainConf) *Chain {
	return NewChain(nil, conf)
}
