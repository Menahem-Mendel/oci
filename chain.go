package oci

import (
	"context"
	"oci/driver"
)

type Chain struct {
	driver driver.Driver
	ref    string

	imageID     string
	networkID   string
	namespaceID string
	containerID string

	pool []driver.Handler
}

func NewChain(conn driver.Conn, ref string) *Chain {
	return &Chain{
		// conn: conn,
		pool: make([]driver.Handler, 0),
		ref:  ref,
	}
}

func (c Chain) ImageID() string {
	return c.imageID
}

func (c Chain) NetworkID() string {
	return c.networkID
}

func (c Chain) NamespaceID() string {
	return c.namespaceID
}

func (c Chain) ContainerID() string {
	return c.containerID
}

func (c *Chain) PullImage(opts ...driver.Option) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) StartContainer(opts ...driver.Option) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) NewNetwork(opts ...driver.Option) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) NewContainer(opts ...driver.Option) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) NewNamespace(opts ...driver.Option) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) Exec(cmd string, args ...string) *Chain {
	handler := func(ctx context.Context) error { return nil }

	c.pool = append(c.pool, driver.HandlerFunc(handler))

	return c
}

func (c *Chain) Commit(ctx context.Context) error {
	for _, handler := range c.pool {
		if err := handler.ServeOCI(ctx); err != nil {
			return err
		}
	}

	return nil
}
