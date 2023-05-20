// Package oci provides a generic driver interface for OCI based runtime engines.
// It defines abstractions for drivers, connections, and requests and responses to the OCI.
// Use this package to interact with different OCI implementations in a uniform way.
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
