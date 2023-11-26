// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package image

import (
	"context"
	"encoding/json"
	"io"
	"oci"
	"oci/driver"
	"time"

	"github.com/opencontainers/go-digest"
)

type Conf struct {
	ID            string
	Architecture  string
	Author        string
	DockerVersion string
	Name          string `json:",omitempty"`
	Os            string
	Tag           string `json:",omitempty"`
	Variant       string
	Env           []string
	Layers        []string
	RepoTags      []string
	Created       *time.Time
	Digest        digest.Digest
	Labels        map[string]string
	LayersData    []Layer
}

type Layer struct {
	Annotations map[string]string
	Digest      digest.Digest
	MIMEType    string // "" if unknown.
	Size        int64  // -1 if unknown.
}

type Spec struct {
	ID            string
	Architecture  string
	Author        string
	DockerVersion string
	Name          string
	Os            string
	Tag           string
	Variant       string
	Env           []string
	Layers        []string
	RepoTags      []string
	Created       *time.Time
	Digest        digest.Digest
	Labels        map[string]string
	LayersData    []Layer
}

func imageBuilder(imgsrv any) (driver.Builder, error) {
	p, ok := imgsrv.(driver.Builder)
	if !ok {
		return nil, oci.ErrUnsupportedOperation
	}

	return p, nil
}

func Build(ctx context.Context, conn driver.Conn, dockerfile io.Reader, conf driver.Configer) (*Spec, error) {
	srv, _ := imageService(conn)
	b, _ := imageBuilder(srv)
	_ = oci.Build(ctx, b)

	return &Spec{}, nil
}

func NewConf(drv driver.Driver, opts ...driver.Option) (driver.Configer, error) {
	conf, _ := drv.Configer("images")

	for range opts {
		// if err := opt.Apply(conf); err == oci.ErrUnsupportedOption {
		// 	// TODO: should return the first unsupported error.
		// 	return nil, err
		// } else if err != nil {
		// 	return nil, err
		// }
	}

	return conf, nil
}

func (c *Conf) Apply(v any) error {
	c, ok := v.(*Conf)
	if !ok {
		// return oci.ErrUnknownConf
		return nil
	}

	return nil
}

func WithArch(arch string) driver.ArchOption {
	return driver.OptionFunc(func(conf driver.Configer) error {
		// conf.Set(arch)
		// return oci.ErrUnsupportedOption
		return nil
	})
}

func imageService(conn driver.Conn) (any, error) {
	imgsrv, err := conn.Prepare("images")
	if err != nil {
		return nil, err
	}

	return imgsrv, nil
}

func imagePuller(imgsrv any) (driver.Puller, error) {
	p, ok := imgsrv.(driver.Puller)
	if !ok {
		return nil, oci.ErrUnsupportedOperation
	}

	return p, nil
}

func Pull(ctx context.Context, conn driver.Conn, ref string) (*Spec, error) {
	s, err := imageService(conn)
	if err != nil {
		return nil, err
	}

	p, err := imagePuller(s)
	if err != nil {
		return nil, err
	}

	_, err = oci.Pull(ctx, p, ref)
	if err != nil {
		return nil, err
	}

	return &Spec{}, nil
}

func imagePusher(imgsrv any) (driver.Pusher, error) {
	p, ok := imgsrv.(driver.Pusher)
	if !ok {
		return nil, oci.ErrUnsupportedOperation
	}

	return p, nil
}

func Push(ctx context.Context, conn driver.Conn, ref, id string) error {
	s, err := imageService(conn)
	if err != nil {
		return err
	}

	p, err := imagePusher(s)
	if err != nil {
		return err
	}

	return oci.Push(ctx, p, ref, id)
}

func Stat(ctx context.Context, conn driver.Conn, id string) (image *Conf, err error) {
	s, err := imageService(conn)
	if err != nil {
		return nil, err
	}

	p, err := imagePusher(s)
	if err != nil {
		return nil, err
	}

	c, err := oci.Stat(ctx, p, ref, id)
	if err != nil {
		return nil, err
	}

	if err := jsonUnmarshal(c, &image); err != nil {
		return nil, err
	}

	return image, nil
}

func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
