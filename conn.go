// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package oci

import (
	"context"
)

type Conn struct {
	uri string
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Begin(ctx context.Context) error {
	return nil
}

func (c *Conn) Prepare(service string) (any, error) {
	return nil, nil
}
