// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package oci

This file is the main entry point for the oci package. It contains definitions for core types such as Request and Response, as well as the main interfaces and functions used to interact with the package.

The oci package provides a driver interface for interacting with different OCI (Open Container Initiative) runtime engines, such as Docker, Podman, containerd, etc. It provides an abstract layer for handling container lifecycle operations in a generic way, allowing the end user to switch between different OCI runtime engines without changing the main application code.

Example usage:

// TODO: Add example usage.

// ... handle the response ...
*/
package oci

import (
	"context"
	"oci/driver"
	"oci/runtime"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("oci: Register nil driver")
	}

	if _, dup := drivers[name]; dup {
		panic("oci: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

func Runtime(driver string) (*runtime.Runtime, error) {
	driversMu.RLock()
	drv, ok := drivers[driver]
	driversMu.RUnlock()
	if !ok {
		return nil, ErrUnregisteredDriver
	}

	r, err := runtime.New(drv)
	if err != nil {
		return nil, err
	}

	return r, nil
}

type Configer interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}

func Open(ctx context.Context, runtime driver.Driver, uri string) (driver.Conn, error) {
	return runtime.Open(ctx, uri)
}

func Pull(ctx context.Context, p driver.Puller, ref string) (string, error) {
	return p.Pull(ctx, ref)
}

func Push(ctx context.Context, p driver.Pusher, ref, id string) error {
	return p.Push(ctx, ref, id)
}

func Stat(ctx context.Context, conf Configer, s driver.Inspector, id string) error {
	return s.Stat(ctx, conf, id)
}

func List(ctx context.Context, conf []Configer, l driver.Lister) error {
	return l.List(ctx, conf)
}

func Create(ctx context.Context, c driver.Creator, id string, args ...string) (string, error) {
	return c.Build(ctx)
}

func Start(ctx context.Context, s driver.Starter, id string) error {
	return s.Start(ctx, id)
}

func Stop(ctx context.Context, s driver.Stoper, id string) error {
	return s.Stop(ctx, id)
}

func Pause(ctx context.Context, p driver.Pauser, id string) error {
	return p.Pause(ctx, id)
}

func Kill(ctx context.Context, k driver.Killer, id string) error {
	return k.Kill(ctx, id)
}
