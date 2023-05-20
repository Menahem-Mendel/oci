// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package oci

This file is the main entry point for the oci package. It contains definitions for core types such as Request and Response, as well as the main interfaces and functions used to interact with the package.

The oci package provides a driver interface for interacting with different OCI (Open Container Initiative) runtime engines, such as Docker, Podman, containerd, etc. It provides an abstract layer for handling container lifecycle operations in a generic way, allowing the end user to switch between different OCI runtime engines without changing the main application code.

Example usage:

drv := oci.NewDriver()
ctx := context.Background()
client, err := oci.NewClient(ctx, drv, "unix:///var/run/docker.sock")

	if err != nil {
	    log.Fatalf("Failed to create client: %v", err)
	}

defer client.Close()

req := oci.NewRequest(ctx, oci.PULL, "nginx:latest", "", "IMAGE", nil)
res, err := client.Do(req)

	if err != nil {
	    log.Fatalf("Failed to pull image: %v", err)
	}

defer res.Body.Close()

// ... handle the response ...
*/
package oci

import (
	"context"
	"io"
)

// Register registers an OCI driver by its name. Currently not implemented.
func Register(drv Driver) {
	// drivers[drv.Name()] = drv
}

// Response represents an OCI response with a body that can be read and closed.
type Response struct {
	// Body is the body of the response. It can be read and should be closed after reading.
	Body io.ReadCloser
}

// NewResponse creates a new Response with the given body.
// If the body is not already an io.ReadCloser, it wraps it with io.NopCloser.
func NewResponse(body io.Reader) *Response {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}

	return &Response{
		Body: rc,
	}
}

// Request represents an OCI request with a method, a reference, an ID, a kind, a body, and a context.
type Request struct {
	Method string          // Method is the OCI method to use for this request.
	Ref    string          // Ref is the reference to the OCI object.
	ID     string          // ID is an optional ID for the OCI object.
	Kind   string          // Kind is the kind of the OCI object.
	Body   io.ReadCloser   // Body is the body of the request. It can be read and should be closed after reading.
	ctx    context.Context // ctx is the context of this request. It can be used to cancel the request.
}

// NewRequest creates a new Request with the given parameters.
// If the body is not already an io.ReadCloser, it wraps it with io.NopCloser.
func NewRequest(ctx context.Context, method, ref, id, kind string, body io.Reader) *Request {
	// TODO: if method is valid

	// TODO: if reference is OCI valid

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}

	return &Request{
		ctx:    ctx,
		Method: method,
		Ref:    ref,
		Kind:   kind,
		Body:   rc,
	}
}

// Context returns the context of this request. Use this to cancel the request.
func (r *Request) Context() context.Context {
	return r.ctx
}
