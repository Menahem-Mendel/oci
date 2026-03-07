# OCI Interface Library

`oci` is a runtime-agnostic interface layer for containers, inspired by `database/sql`.

The core packages are **standard library only**:
- `oci`
- `oci/driver`
- `oci/image`

Runtime-specific adapters live in subpackages (for example `oci/pkg/podman`) and may use runtime-specific logic.

## Architecture

### 1) High-level API (`oci`)
- Driver registry (`Register`, `Drivers`)
- Runtime handle (`Open`, `NewRuntime`, `Runtime.Connect`, `Runtime.Close`)
- Capability discovery (`Runtime.Capabilities()`)
- Stable operations for images and containers

### 2) Adapter contract (`oci/driver`)
- `Driver` and `Conn` interfaces
- Optional capability interfaces (pull image, create container, exec, etc.)
- Typed request/response structs shared by all adapters

### 3) Runtime adapters (`oci/pkg/...`)
- Implement `driver.Driver` and one or more optional capability interfaces
- Register in `init()`

## Standards Mapping

This project intentionally separates concerns because container standards are split:

- OCI Image Spec: image format (`oci/image` data types)
- OCI Runtime Spec: low-level runtime execution contract (implemented by runtimes like runc/crun)
- OCI Distribution Spec: registry pull/push transport semantics
- CRI: Kubernetes runtime contract (outside OCI scope, kubelet-specific)
- CNI: networking plugin contract (outside OCI image/runtime scope)
- CDI: device injection contract across runtimes

There is no single ISO standard that unifies all runtime control APIs. The library uses OCI-aligned data models and adapter capabilities to bridge runtime differences.

## Quick Start

```go
package main

import (
	"context"
	"log"
	"oci"
	_ "oci/pkg/podman"
	"time"
)

func main() {
	rt, err := oci.Open("podman", "unix:///run/podman/podman.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rt.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	img, err := rt.PullImage(ctx, "docker.io/library/alpine:latest")
	if err != nil {
		log.Fatal(err)
	}

	_, err = rt.CreateContainer(ctx, oci.ContainerSpec{
		Name:    "demo",
		Image:   img.Reference,
		Command: []string{"sh", "-lc", "echo hello"},
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Implementing A Custom Adapter

```go
package myruntime

import (
	"oci"
	"oci/driver"
)

type Driver struct{}

type Conn struct{}

func (d *Driver) Open(dsn string) (driver.Conn, error) { return &Conn{}, nil }
func (c *Conn) Close() error { return nil }

func init() {
	oci.Register("myruntime", &Driver{})
}
```

Then incrementally add optional interfaces from `oci/driver` (such as `ImagePuller` or `ContainerCreator`) as your runtime supports them.

## Project Status

Release candidate architecture: stable core API, adapter boundary, and tests for core behavior.
