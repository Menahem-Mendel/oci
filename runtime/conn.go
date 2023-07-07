// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package runtime

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

// func (c *Conn) Prepare(service string) (io.ReadWriteCloser, driver.ParserFunc, error) {
// 	return nil, nil, nil
// }
