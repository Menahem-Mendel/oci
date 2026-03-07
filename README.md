# OCI Interface Library

`oci` is a runtime-agnostic container interface, modeled after the ergonomics of `database/sql`.

Core goal:
- Simple app API in `oci`
- Complex runtime differences hidden inside adapters
- Adapter interfaces in `oci/driver` use builtin data types only (`string`, `[]string`, `map[...]...`, `bool`, `int`, `[]byte`, `any`)

## Package Layout

- `oci`: high-level runtime handle (`Open`, `PullImage`, `CreateContainer`, `ExecInContainer`, ...)
- `oci/driver`: low-level adapter contracts
- `oci/image`: OCI image constants/types (stdlib only)
- `oci/pkg/podman`: Podman adapter implementation

## 60-Second Example

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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := rt.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	imageID, err := rt.PullImage(ctx, "docker.io/library/alpine:latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	containerID, err := rt.CreateContainer(ctx, map[string]any{
		oci.SpecImage:   imageID,
		oci.SpecCommand: []string{"sh", "-lc", "echo hello"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := rt.StartContainer(ctx, containerID); err != nil {
		log.Fatal(err)
	}
}
```

## Common Operations

```go
// Pull with options
imageID, _ := rt.PullImage(ctx, "docker.io/library/nginx:latest", map[string]string{
	oci.OptionPlatform: "linux/amd64",
})

// Inspect image/container
img, _ := rt.InspectImage(ctx, imageID)
ctr, _ := rt.InspectContainer(ctx, "my-container")

// List resources
images, _ := rt.ListImages(ctx, nil)
containers, _ := rt.ListContainers(ctx, map[string]string{"all": "1"})

// Exec
out, _ := rt.ExecInContainer(ctx, "my-container", []string{"/bin/sh", "-lc", "id"}, map[string]any{
	oci.ExecOptionTTY: false,
})
exitCode := out[oci.ExecResultExitCode]
stdout := out[oci.ExecResultStdout]
stderr := out[oci.ExecResultStderr]

_ = exitCode
_ = stdout
_ = stderr
```

## Data Conventions

Container spec keys used by `CreateContainer`:
- `oci.SpecName` (string)
- `oci.SpecImage` (string, required)
- `oci.SpecCommand` ([]string)
- `oci.SpecEnv` (map[string]string)
- `oci.SpecLabels` (map[string]string)
- `oci.SpecWorkingDir` (string)
- `oci.SpecNetworkMode` (string)

Pull option keys used by `PullImage`:
- `oci.OptionPlatform` (example: `linux/amd64`)
- `oci.OptionAllTags` (`true`/`false` as string)
- `oci.OptionInsecureSkipTLS` (`true`/`false` as string)
- `oci.OptionUsername`
- `oci.OptionPassword`

Exec option/result keys:
- options: `oci.ExecOptionEnv`, `oci.ExecOptionTTY`, `oci.ExecOptionStdin`
- result: `oci.ExecResultExitCode`, `oci.ExecResultStdout`, `oci.ExecResultStderr`

## Driver API (Builtin-Only)

```go
package driver

type Driver interface {
	Open(dsn string) (Conn, error)
}

type Conn interface {
	Close() error
}

type Puller interface {
	Pull(ctx context.Context, reference string, options map[string]string) (string, error)
}

type Inspector interface {
	Inspect(ctx context.Context, kind string, idOrRef string) (map[string]any, error)
}

type Creator interface {
	Create(ctx context.Context, kind string, spec map[string]any) (string, error)
}
```

Adapters can implement only the capabilities they support. `oci.Runtime` checks capability interfaces at runtime and returns `oci.ErrUnsupported` when missing.

## Implementing an Adapter

```go
package myruntime

import (
	"context"
	"oci"
	"oci/driver"
)

type Driver struct{}
type Conn struct{}

func (d *Driver) Open(dsn string) (driver.Conn, error) { return &Conn{}, nil }
func (c *Conn) Close() error { return nil }
func (c *Conn) Pull(ctx context.Context, reference string, options map[string]string) (string, error) {
	return "image-id", nil
}

func init() {
	oci.Register("myruntime", &Driver{})
}
```

## Standards Scope

This library aligns with OCI and related ecosystem standards but does not pretend there is one universal control-plane standard.

- OCI Image Spec: image format
- OCI Runtime Spec: runtime bundle/process model
- OCI Distribution Spec: registry push/pull API
- CRI/CNI/CDI: adjacent integration layers used by orchestrators and runtimes

## Status

Refactored architecture with builtin-only driver interfaces, adapter isolation, and passing tests.
