// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package image

import (
	"context"
	"encoding/json"
	"oci"
	"oci/driver"
)

func imageService(conn driver.Conn) (any, error) {
	imgsrv, err := conn.Prepare("image")
	if err != nil {
		return nil, err
	}

	return imgsrv, nil
}

func puller(s any) (driver.Puller, error) {
	p, ok := s.(driver.Puller)
	if !ok {
		return nil, oci.ErrUnsupportedOperation
	}

	return p, nil
}

func Pull(ctx context.Context, conn driver.Conn, ref string) (string, error) {
	s, _ := imageService(conn)
	p, _ := puller(s)

	return oci.Pull(ctx, p, ref)
}

func Stat(ctx context.Context, conn driver.Conn, id string) (image *Image, err error) {
	s, _ := imageService(conn)
	p, _ := imagePusher(s)
	c, _ := oci.Stat(ctx, p, ref, id)

	_ = jsonUnmarshal(c, &image)

	return
}

func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
