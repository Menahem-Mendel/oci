// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the Client type in the oci package.
// A Client maintains a connection with an OCI and allows for sending requests and receiving responses.

// Example usage:

// drv := oci.NewDriver()
// ctx := context.Background()
// client, err := oci.NewClient(ctx, drv, "localhost:5000")
//
//	if err != nil {
//	    log.Fatalf("Failed to create client: %v", err)
//	}
//
// defer client.Close()
// // ... use the client ...
package oci

import (
	"context"
	"errors"
	"sync"
)

// Client is a struct that encapsulates an OCI driver and a connection.
// It also includes a RWMutex for managing concurrent access to the connection.
type Client struct {
	// driver is the Driver interface for managing communication with the OCI runtime
	driver Driver

	// conn is the Conn interface for managing the connection to the OCI runtime
	conn Conn

	// RWMutex is used for managing concurrent access to the conn field
	sync.RWMutex
}

// NewClient constructs a new OCI client.
// It takes a context for managing the connection lifecycle, a Driver for managing communication with the OCI runtime,
// and a URI for specifying the target runtime.
// It returns a pointer to the Client and any error encountered.
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

// Close is responsible for closing the connection to the OCI runtime.
// It is safe for concurrent use.
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

// Do sends a Request to the OCI runtime and returns the Response or an error.
// It checks whether a connection to the OCI runtime has been established before sending the request.
// This method is safe for concurrent use.
func (c *Client) Do(req *Request) (*Response, error) {
	c.RLock()
	defer c.RUnlock()

	if c.conn == nil {
		return nil, errors.New("No connection established")
	}

	return c.do(req)
}

// do is an unexported method that sends a Request to the OCI runtime and returns the Response or an error.
// Unlike Do, it does not check whether a connection has been established and is not safe for concurrent use.
func (c *Client) do(req *Request) (*Response, error) {
	return c.driver.Handler(req.Method).ServeOCI(req)
}

// Pull initiates a PULL request to the OCI runtime.
// It takes a context for managing the request lifecycle and a reference string for specifying the target image.
// It returns the Response from the OCI runtime or an error.
func (c *Client) Pull(ctx context.Context, ref string) (*Response, error) {
	method := PULL

	req := &Request{
		Method: string(method),
		Ref:    ref,
		Kind:   "IMAGE",
	}

	return c.Do(req)
}
