// Package driver defines the low-level contracts implemented by runtime
// adapters (Podman, containerd, Docker, etc.).
package driver

import (
	"context"
	"errors"
)

var ErrNotSupported = errors.New("oci/driver: operation not supported")

const (
	KindImage     = "image"
	KindContainer = "container"
	KindNetwork   = "network"
	KindVolume    = "volume"
)

// Driver opens a connection to a runtime endpoint.
// The dsn format is driver-specific.
type Driver interface {
	Open(dsn string) (Conn, error)
}

// Conn is the minimal connection contract shared by all runtime adapters.
type Conn interface {
	Close() error
}

// Optional capability interfaces.

type Pinger interface {
	Ping(ctx context.Context) error
}

// Puller pulls a resource reference (typically an image) and returns its local ID.
type Puller interface {
	Pull(ctx context.Context, reference string, options map[string]string) (string, error)
}

// Pusher pushes a resource reference and returns a runtime-specific result string.
type Pusher interface {
	Push(ctx context.Context, reference string, options map[string]string) (string, error)
}

// Inspector returns a dynamic resource document.
type Inspector interface {
	Inspect(ctx context.Context, kind string, idOrRef string) (map[string]any, error)
}

// Lister returns dynamic resource documents.
type Lister interface {
	List(ctx context.Context, kind string, filters map[string]string) ([]map[string]any, error)
}

// Remover removes a resource.
type Remover interface {
	Remove(ctx context.Context, kind string, idOrRef string, force bool) error
}

// Creator creates a resource and returns its ID.
type Creator interface {
	Create(ctx context.Context, kind string, spec map[string]any) (string, error)
}

// Starter starts a resource.
type Starter interface {
	Start(ctx context.Context, kind string, id string) error
}

// Stopper stops a resource.
type Stopper interface {
	Stop(ctx context.Context, kind string, id string, timeoutSeconds int) error
}

// Execer executes a command in a container-like resource.
// Returned map is expected to include keys like "exit_code", "stdout", "stderr".
type Execer interface {
	Exec(ctx context.Context, id string, command []string, options map[string]any) (map[string]any, error)
}
