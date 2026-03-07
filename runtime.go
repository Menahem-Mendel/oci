package oci

import (
	"context"
	"fmt"
	"oci/driver"
	"sort"
	"strings"
	"sync"
	"time"
)

// Runtime is the high-level handle for interacting with a container runtime,
// similar in spirit to database/sql's DB handle.
type Runtime struct {
	driverName string
	dsn        string
	drv        driver.Driver

	mu     sync.RWMutex
	conn   driver.Conn
	closed bool
}

func (r *Runtime) DriverName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.driverName
}

func (r *Runtime) DSN() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.dsn
}

// Connect opens a runtime connection. If an existing connection exists,
// it is replaced.
func (r *Runtime) Connect(dsn string) error {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return fmt.Errorf("%w: dsn is required", ErrInvalidArgument)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return ErrRuntimeClosed
	}

	conn, err := r.drv.Open(dsn)
	if err != nil {
		return err
	}

	if r.conn != nil {
		_ = r.conn.Close()
	}

	r.conn = conn
	r.dsn = dsn
	return nil
}

func (r *Runtime) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true

	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

func (r *Runtime) connOrError() (driver.Conn, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return nil, ErrRuntimeClosed
	}
	if r.conn == nil {
		return nil, ErrNotConnected
	}
	return r.conn, nil
}

func (r *Runtime) Ping(ctx context.Context) error {
	conn, err := r.connOrError()
	if err != nil {
		return err
	}
	pinger, ok := conn.(driver.Pinger)
	if !ok {
		return ErrUnsupported
	}
	return pinger.Ping(ctx)
}

func (r *Runtime) Capabilities() []Capability {
	conn, err := r.connOrError()
	if err != nil {
		return nil
	}

	caps := make([]Capability, 0, 10)
	if _, ok := conn.(driver.Pinger); ok {
		caps = append(caps, CapabilityPing)
	}
	if _, ok := conn.(driver.Puller); ok {
		caps = append(caps, CapabilityPull)
	}
	if _, ok := conn.(driver.Pusher); ok {
		caps = append(caps, CapabilityPush)
	}
	if _, ok := conn.(driver.Inspector); ok {
		caps = append(caps, CapabilityInspect)
	}
	if _, ok := conn.(driver.Lister); ok {
		caps = append(caps, CapabilityList)
	}
	if _, ok := conn.(driver.Remover); ok {
		caps = append(caps, CapabilityRemove)
	}
	if _, ok := conn.(driver.Creator); ok {
		caps = append(caps, CapabilityCreate)
	}
	if _, ok := conn.(driver.Starter); ok {
		caps = append(caps, CapabilityStart)
	}
	if _, ok := conn.(driver.Stopper); ok {
		caps = append(caps, CapabilityStop)
	}
	if _, ok := conn.(driver.Execer); ok {
		caps = append(caps, CapabilityExec)
	}

	sort.Slice(caps, func(i, j int) bool { return caps[i] < caps[j] })
	return caps
}

func (r *Runtime) Pull(ctx context.Context, reference string, options map[string]string) (string, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", ErrInvalidReference
	}

	conn, err := r.connOrError()
	if err != nil {
		return "", err
	}
	puller, ok := conn.(driver.Puller)
	if !ok {
		return "", ErrUnsupported
	}

	return puller.Pull(ctx, reference, cloneStringMap(options))
}

func (r *Runtime) Push(ctx context.Context, reference string, options map[string]string) (string, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return "", ErrInvalidReference
	}

	conn, err := r.connOrError()
	if err != nil {
		return "", err
	}
	pusher, ok := conn.(driver.Pusher)
	if !ok {
		return "", ErrUnsupported
	}

	return pusher.Push(ctx, reference, cloneStringMap(options))
}

func (r *Runtime) Inspect(ctx context.Context, kind string, idOrRef string) (map[string]any, error) {
	kind = strings.TrimSpace(kind)
	idOrRef = strings.TrimSpace(idOrRef)
	if kind == "" || idOrRef == "" {
		return nil, fmt.Errorf("%w: kind and id/reference are required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return nil, err
	}
	inspector, ok := conn.(driver.Inspector)
	if !ok {
		return nil, ErrUnsupported
	}

	return inspector.Inspect(ctx, kind, idOrRef)
}

func (r *Runtime) List(ctx context.Context, kind string, filters map[string]string) ([]map[string]any, error) {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		return nil, fmt.Errorf("%w: kind is required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return nil, err
	}
	lister, ok := conn.(driver.Lister)
	if !ok {
		return nil, ErrUnsupported
	}

	return lister.List(ctx, kind, cloneStringMap(filters))
}

func (r *Runtime) Remove(ctx context.Context, kind string, idOrRef string, force bool) error {
	kind = strings.TrimSpace(kind)
	idOrRef = strings.TrimSpace(idOrRef)
	if kind == "" || idOrRef == "" {
		return fmt.Errorf("%w: kind and id/reference are required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return err
	}
	remover, ok := conn.(driver.Remover)
	if !ok {
		return ErrUnsupported
	}

	return remover.Remove(ctx, kind, idOrRef, force)
}

func (r *Runtime) Create(ctx context.Context, kind string, spec map[string]any) (string, error) {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		return "", fmt.Errorf("%w: kind is required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return "", err
	}
	creator, ok := conn.(driver.Creator)
	if !ok {
		return "", ErrUnsupported
	}

	return creator.Create(ctx, kind, cloneAnyMap(spec))
}

func (r *Runtime) Start(ctx context.Context, kind string, id string) error {
	kind = strings.TrimSpace(kind)
	id = strings.TrimSpace(id)
	if kind == "" || id == "" {
		return fmt.Errorf("%w: kind and id are required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return err
	}
	starter, ok := conn.(driver.Starter)
	if !ok {
		return ErrUnsupported
	}

	return starter.Start(ctx, kind, id)
}

func (r *Runtime) Stop(ctx context.Context, kind string, id string, timeout time.Duration) error {
	kind = strings.TrimSpace(kind)
	id = strings.TrimSpace(id)
	if kind == "" || id == "" {
		return fmt.Errorf("%w: kind and id are required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return err
	}
	stopper, ok := conn.(driver.Stopper)
	if !ok {
		return ErrUnsupported
	}

	seconds := int(timeout / time.Second)
	if seconds < 0 {
		seconds = 0
	}
	return stopper.Stop(ctx, kind, id, seconds)
}

func (r *Runtime) Exec(ctx context.Context, id string, command []string, options map[string]any) (map[string]any, error) {
	id = strings.TrimSpace(id)
	if id == "" || len(command) == 0 {
		return nil, fmt.Errorf("%w: id and command are required", ErrInvalidArgument)
	}

	conn, err := r.connOrError()
	if err != nil {
		return nil, err
	}
	execer, ok := conn.(driver.Execer)
	if !ok {
		return nil, ErrUnsupported
	}

	return execer.Exec(ctx, id, append([]string(nil), command...), cloneAnyMap(options))
}

// Convenience wrappers for common kinds.

func (r *Runtime) PullImage(ctx context.Context, reference string, options map[string]string) (string, error) {
	return r.Pull(ctx, reference, options)
}

func (r *Runtime) InspectImage(ctx context.Context, idOrRef string) (map[string]any, error) {
	return r.Inspect(ctx, driver.KindImage, idOrRef)
}

func (r *Runtime) ListImages(ctx context.Context, filters map[string]string) ([]map[string]any, error) {
	return r.List(ctx, driver.KindImage, filters)
}

func (r *Runtime) RemoveImage(ctx context.Context, idOrRef string, force bool) error {
	return r.Remove(ctx, driver.KindImage, idOrRef, force)
}

func (r *Runtime) CreateContainer(ctx context.Context, spec map[string]any) (string, error) {
	image, _ := spec[SpecImage].(string)
	if strings.TrimSpace(image) == "" {
		return "", fmt.Errorf("%w: container image is required in spec[%q]", ErrInvalidArgument, SpecImage)
	}
	return r.Create(ctx, driver.KindContainer, spec)
}

func (r *Runtime) StartContainer(ctx context.Context, id string) error {
	return r.Start(ctx, driver.KindContainer, id)
}

func (r *Runtime) StopContainer(ctx context.Context, id string, timeout time.Duration) error {
	return r.Stop(ctx, driver.KindContainer, id, timeout)
}

func (r *Runtime) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	return r.Inspect(ctx, driver.KindContainer, id)
}

func (r *Runtime) ListContainers(ctx context.Context, filters map[string]string) ([]map[string]any, error) {
	return r.List(ctx, driver.KindContainer, filters)
}

func (r *Runtime) RemoveContainer(ctx context.Context, id string, force bool) error {
	return r.Remove(ctx, driver.KindContainer, id, force)
}

func (r *Runtime) ExecInContainer(ctx context.Context, id string, command []string, options map[string]any) (map[string]any, error) {
	return r.Exec(ctx, id, command, options)
}
