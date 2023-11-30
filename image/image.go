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
	srv driver.Service
}

func (s *Service) Pull(ctx context.Context, ref string) (string, error) {
	//  s.service.ServeOCI(ctx, driver.Puller)
	// driver.Puller(service, ref) (string, error)
	return "", nil
}

func Pull(ctx context.Context, conn driver.Conn, ref string) (string, error) {
	var p driver.Puller

	// var service driver.Service
	//
	// oci.Pull(ctx, service.(driver.Puller), ref)

	conn.ExecPull(ctx, oci.Pull, ref)

	return "", nil
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
