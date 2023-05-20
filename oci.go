package oci

import (
	"context"
	"io"
)

func Register(drv Driver) {
	// drivers[drv.Name()] = drv
}

type Response struct {
	Body io.ReadCloser
}

func NewResponse(body io.Reader) *Response {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}

	return &Response{
		Body: rc,
	}
}

type Request struct {
	Method string
	Ref    string
	ID     string
	Kind   string
	Body   io.ReadCloser
	ctx    context.Context
}

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

func (r *Request) Context() context.Context {
	return r.ctx
}
