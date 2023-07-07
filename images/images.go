// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package images

import (
	"context"
	"oci"
	"oci/driver"
	"sync"
	"time"

	"github.com/opencontainers/go-digest"
)

type Conf struct {
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

type service struct {
	runtime driver.Driver

	mu sync.Mutex
}

func Handle(ctx context.Context, h driver.Handler) error {
	return h.ServeOCI(ctx)
}

func HandleFunc(ctx context.Context, f func(context.Context) error) error {
	return Handle(ctx, driver.HandlerFunc(f))
}

// func Stat(ctx context.Context, s driver.Server, id string) (image *Conf, err error) {
// 	return nil, err
// }

func Pull(ctx context.Context, r driver.Driver, ref string) (id string, err error) {
	p := puller{
		r:   r,
		ref: ref,
	}

	is, err := r.Service("image")
	if err != nil {
		return "", err
	}

	r.ServeOCI(ctx, p)
	return "", nil
}

type puller struct {
	r oci.Runtime

	ref string
}

func (p puller) ServeOCI(ctx context.Context) error {
	return p.p.Pull(ctx, p.Name)
}
