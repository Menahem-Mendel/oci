package oci

import "context"

// Chain is a lightweight workflow helper for sequencing runtime operations.
type Chain struct {
	ops []func(context.Context) error
}

func NewChain() *Chain {
	return &Chain{ops: make([]func(context.Context) error, 0)}
}

func (c *Chain) Then(op func(context.Context) error) *Chain {
	if op == nil {
		return c
	}
	c.ops = append(c.ops, op)
	return c
}

func (c *Chain) Commit(ctx context.Context) error {
	for _, op := range c.ops {
		if err := op(ctx); err != nil {
			return err
		}
	}
	return nil
}
