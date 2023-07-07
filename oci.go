// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package oci

This file is the main entry point for the oci package. It contains definitions for core types such as Request and Response, as well as the main interfaces and functions used to interact with the package.

The oci package provides a driver interface for interacting with different OCI (Open Container Initiative) runtime engines, such as Docker, Podman, containerd, etc. It provides an abstract layer for handling container lifecycle operations in a generic way, allowing the end user to switch between different OCI runtime engines without changing the main application code.

Example usage:

drv := oci.NewDriver()
ctx := context.Background()
client, err := oci.NewClient(ctx, drv, "unix:///var/run/docker.sock")

	if err != nil {
	    log.Fatalf("Failed to create client: %v", err)
	}

defer client.Close()

req := oci.NewRequest(ctx, oci.PULL, "nginx:latest", "", "IMAGE", nil)
res, err := client.Do(req)

	if err != nil {
	    log.Fatalf("Failed to pull image: %v", err)
	}

defer res.Body.Close()

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
		panic("oci: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("oci: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

func Open(ctx context.Context, driver, uri string) (*runtime.Runtime, error) {
	driversMu.RLock()
	drv, ok := drivers[driver]
	delete(drivers, driver)
	driversMu.RUnlock()
	if !ok {
		return nil, ErrUnregisteredDriver
	}

	r, err := runtime.New(ctx, drv)
	if err != nil {
		return nil, err
	}

	_, err = r.Open(ctx, uri)
	if err != nil {
		return nil, err
	}

	return r, nil
}
