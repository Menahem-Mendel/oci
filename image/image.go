// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package image

import (
	"context"
	"oci/driver"
)

type Service struct {
	Conn driver.Conn
}

func NewService(conn driver.Conn) *Service {
	return &Service{Conn: conn}
}

func (s *Service) ServeOCI(ctx context.Context) error {
	return nil
}

func (p *Service) Pull(ctx context.Context, dsn string) (string, error) {
	return "", nil
}

func (p *Service) Push(ctx context.Context, dsn, id string) error {
	return nil
}

func (p *Service) Stat(ctx context.Context, id string) (map[string]any, error) {
	return nil, nil
}

func (p *Service) Remove(ctx context.Context, id string) error {
	return nil
}

type pullOption struct {
}

type puller struct {
}

func (p *puller) Options(drv driver.Driver) []pullOption {
	// var opt driver.Option

	return nil
}

// oci/image, oci, oci/driver, pkg/podman, pkg/docker,
// pkg/podman <- oci, oci/image
// oci.driver <-
// oci/image <- oci.driver
type PodmanDriver struct {
	conns []oci.Conn
}
type conn struct {
}

type Runtime struct {

}

rt = oci.NewRuntime("podman")
conn = rt.Open("socket.sock")


type Image struct {

}
