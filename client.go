package oci

import (
	"context"
)

type Client struct {
	driver Driver
}

func NewClient(drv Driver) (*Client, error) {
	return &Client{
		driver: drv,
	}, nil
}

func (c *Client) Connect(ctx context.Context, sock string) error {
	return c.driver.Connect(ctx, sock)
}

func (c *Client) Close() error {
	return c.driver.Close()
}

func (c *Client) Do(req *Request) (*Response, error) {
	return c.driver.Handler(req.Method).ServeOCI(req)
}

func (c *Client) Pull(ref string) (*Response, error) {
	method := PULL

	req := &Request{
		ctx:    context.Background(),
		Method: string(method),
		Ref:    ref,
		Kind:   "IMAGE",
	}

	return c.driver.Handler(string(method)).ServeOCI(req)
}
