// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package image

import (
	"context"
	"oci"
	"oci/driver"
)

type Service struct {
	Puller driver.Puller
}

type service struct {
}

// func imageService(conn driver.Conn) (any, error) {
// 	imgsrv, err := conn.Prepare("image")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return imgsrv, nil
// }

// func puller(s any) (driver.Puller, error) {
// 	p, ok := s.(driver.Puller)
// 	if !ok {
// 		return nil, oci.ErrUnsupportedOperation
// 	}

// 	return p, nil
// }

// func Pull(ctx context.Context, conn driver.Conn, ref string) (string, error) {
// 	s, _ := imageService(conn)
// 	p, _ := puller(s)

// 	return oci.Pull(ctx, p, ref)
// }

func Pull(ctx context.Context, conn driver.Conn, ref string) (string, error) {
	var p driver.Puller

	conn.Prepare()

	return oci.Pull(ctx, p, ref)
}

func pull(ctx context.Context, p driver.Puller, ref string) {

}

type image struct {
	drv driver.Driver

	conn driver.Conn
}

func (p *image) Pull(ctx context.Context, ref string) (id string, err error) {
	p.conn.Driver()

	return
}
