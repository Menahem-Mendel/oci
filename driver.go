/*
Package oci - Driver

Author: Menahem-Mendel Gelfand

Copyright: Copyright 2023, Menahem-Mendel Gelfand

License: This source code is licensed under the BSD 3-Clause License. You may obtain a copy of the License at:
https://opensource.org/licenses/BSD-3-Clause

This file contains the definitions for the Driver interface and related types in the oci package.
A Driver is the core interface in this package and is designed to enable interaction with different OCI implementations.

Example usage:

drv := oci.NewDriver()
ctx := context.Background()
conn, err := drv.Connect(ctx, "localhost:5000")
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

// ... use the conn ...
*/

package oci

import "context"

// Driver is an interface for handling communication with an OCI compatible runtime.
// It embeds the Connector interface and provides a Handler method.
type Driver interface {
	// Connector provides a method for establishing a connection with an OCI compatible runtime.
	Connector

	// Handler returns a Handler based on the provided string (usually a method name).
	Handler(string) Handler
}

// Conn is an interface for a connection to an OCI compatible runtime.
// It provides a method for closing the connection.
type Conn interface {
	// Close closes the connection to an OCI compatible runtime.
	// It should return an error if the connection cannot be closed.
	Close() error
}

// Connector is an interface for establishing a connection with an OCI compatible runtime.
type Connector interface {
	// Connect establishes a connection to an OCI compatible runtime based on the provided URI.
	// The context is used to control cancellation of the connection process.
	// It should return a Conn for the established connection or an error if the connection process fails.
	Connect(ctx context.Context, uri string) (Conn, error)
}

// Handler is an interface for handling a request to an OCI compatible runtime.
// It provides a method for serving an OCI request.
type Handler interface {
	// ServeOCI serves a Request to an OCI compatible runtime.
	// It should return a Response or an error if the request fails.
	ServeOCI(r *Request) (*Response, error)
}

// HandlerFunc is a function type that implements the Handler interface.
type HandlerFunc func(r *Request) (*Response, error)

// ServeOCI serves a Request to an OCI compatible runtime.
// It simply calls the function h with the request r.
// It returns a Response or an error if the request fails.
func (h HandlerFunc) ServeOCI(r *Request) (*Response, error) {
	return h(r)
}
