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
	"errors"
	"io"
	"oci/driver"
	"sync"
	"time"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
	services  = make(map[driver.Driver]string)
)

func Register(name string, driver driver.Driver, services ...any) {
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

// func Handle(h driver.Handler) {
// 	h.ServeOCI(ctx)
// }

func NewRuntime(driver string) (*Runtime, error) {
	driversMu.RLock()
	drv, ok := drivers[driver]
	driversMu.RUnlock()
	if !ok {
		return nil, ErrUnregisteredDriver
	}

	// _, cancel := context.WithCancel(ctx)
	return &Runtime{
		// cancel: cancel,
		driver: drv,
	}, nil
}

// type daemonConn struct {
// 	c Conn
// }

// func (d *daemonConn) Begin(ctx context.Context) error {
// 	d.c.Begin(ctx)
// }

// type Configer interface {
// 	Set(key string, value any) error
// 	Get(key string) (any, error)
// }

type driverConn struct {
	db        *Runtime
	createdAt time.Time

	sync.Mutex  // guards following
	ci          driver.Conn
	needReset   bool // The connection session should be reset before use if true.
	closed      bool
	finalClosed bool // ci.Close has been called
	// openStmt    map[*driverStmt]bool

	// guarded by db.mu
	inUse      bool
	rtmuClosed bool      // same as closed, but guarded by rt.mu, for removeClosedStmtLocked
	returnedAt time.Time // Time the connection was created or returned.
	onPut      []func()  // code (with db.mu held) run when conn is next returned
}

// Open a new connection to the
func Open(runtime driver.Driver, uri string) (driver.Conn, error) {
	if runtime == nil {
		return nil, errors.New("oci: no driver is provided")
	}

	return runtime.Open(uri)
}

func Pull(ctx context.Context, p driver.Puller, dsn string) (string, error) {
	p.Pull(ctx, dsn)

	return p.Pull(ctx, dsn)
}

// func Pull(ctx context.Context, p driver.Puller, args ...any) error {
// 	if len(args) != 1 {
// 		return nil
// 	}

// 	dsn, ok := args[0].(string)
// 	if !ok {
// 		return nil
// 	}
// 	return p.Pull(ctx, dsn)
// }

// func Puller(drv driver.Driver, h driver.Handler) driver.Puller {

// }

// func Push(ctx context.Context, p driver.Pusher, dsn, id string) error {
// 	return p.Push(ctx, dsn, id)
// }

// func Stat(ctx context.Context, conf Configer, s driver.Inspector, id string) error {
// 	return s.Stat(ctx, conf, id)
// }

// func List(ctx context.Context, conf []Configer, l driver.Lister) error {
// 	return l.List(ctx, conf)
// }

// func Create(ctx context.Context, c driver.Creator, id string, args ...string) (string, error) {
// 	return c.Build(ctx)
// }

// func Start(ctx context.Context, s driver.Starter, id string) error {
// 	return s.Start(ctx, id)
// }

// func Stop(ctx context.Context, s driver.Stoper, id string) error {
// 	return s.Stop(ctx, id)
// }

// func Pause(ctx context.Context, p driver.Pauser, id string) error {
// 	return p.Pause(ctx, id)
// }

// func Kill(ctx context.Context, k driver.Killer, id string) error {
// 	return k.Kill(ctx, id)
// }
