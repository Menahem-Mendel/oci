// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package image

import (
	"context"
	"oci/driver"
)

// func Pull(ctx context.Context, conn driver.Conn, ref string) (string, error) {
// 	// oci.Pull(ctx, conn, ref)
// }

type Service struct {
}

func (s *Service) ServeOCI() {

}

// func (s *Service) Pull(ctx context.Context, ref string) (string, error) {

// }

type image struct {
}

func (p *image) Pull(ctx context.Context, ref string) (string, error) {
	return "", nil
}

func NewPuller(drv driver.Driver) *image {
	return &image{}
}

type Puller struct {
}
