package oci

import (
	"context"
	"errors"
	"sync"
)

type Client struct {
	driver Driver

	conn Conn

	sync.RWMutex
}

func NewClient(ctx context.Context, drv Driver, uri string) (*Client, error) {
	conn, err := drv.Connect(ctx, uri)
	if err != nil {
		return nil, err
	}

	return &Client{
		driver: drv,
		conn:   conn,
	}, nil

}

func (c *Client) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *Client) Do(req *Request) (*Response, error) {
	c.RLock()
	defer c.RUnlock()

	if c.conn == nil {
		return nil, errors.New("No connection established")
	}

	return c.do(req)
}

func (c *Client) do(req *Request) (*Response, error) {
	return c.driver.Handler(req.Method).ServeOCI(req)
}

func (c *Client) Pull(ctx context.Context, ref string) (*Response, error) {
	method := PULL

	req := &Request{
		Method: string(method),
		Ref:    ref,
		Kind:   "IMAGE",
	}

	return c.Do(req)
}
