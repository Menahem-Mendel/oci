package oci_test

import (
	"context"
	"oci"
	"oci/driver"
	"strconv"
	"testing"
	"time"
)

type fakeDriver struct {
	conn driver.Conn
	dsn  string
}

func (d *fakeDriver) Open(dsn string) (driver.Conn, error) {
	d.dsn = dsn
	return d.conn, nil
}

type fakeConn struct {
	started string
	stopped string
	closed  bool
}

func (c *fakeConn) Close() error {
	c.closed = true
	return nil
}

func (c *fakeConn) Ping(ctx context.Context) error {
	return nil
}

func (c *fakeConn) Pull(ctx context.Context, reference string, options map[string]string) (string, error) {
	return "img-1", nil
}

func (c *fakeConn) Create(ctx context.Context, kind string, spec map[string]any) (string, error) {
	if kind != driver.KindContainer {
		return "", driver.ErrNotSupported
	}
	return "ctr-1", nil
}

func (c *fakeConn) Start(ctx context.Context, kind string, id string) error {
	c.started = kind + ":" + id
	return nil
}

func (c *fakeConn) Stop(ctx context.Context, kind string, id string, timeoutSeconds int) error {
	c.stopped = kind + ":" + id + ":" + strconv.Itoa(timeoutSeconds)
	return nil
}

func uniqueDriverName(t *testing.T) string {
	t.Helper()
	return "fake-" + strconv.FormatInt(time.Now().UnixNano(), 36)
}

func TestOpenAndBasicFlow(t *testing.T) {
	conn := &fakeConn{}
	drv := &fakeDriver{conn: conn}
	name := uniqueDriverName(t)
	oci.Register(name, drv)

	rt, err := oci.Open(name, "unix:///tmp/fake.sock")
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer rt.Close()

	if err := rt.Ping(context.Background()); err != nil {
		t.Fatalf("ping: %v", err)
	}

	imageID, err := rt.PullImage(context.Background(), "docker.io/library/alpine:latest", nil)
	if err != nil {
		t.Fatalf("pull image: %v", err)
	}
	if imageID == "" {
		t.Fatalf("expected image id")
	}

	ctrID, err := rt.CreateContainer(context.Background(), map[string]any{
		oci.SpecName:  "demo",
		oci.SpecImage: imageID,
	})
	if err != nil {
		t.Fatalf("create container: %v", err)
	}

	if err := rt.StartContainer(context.Background(), ctrID); err != nil {
		t.Fatalf("start container: %v", err)
	}
	if err := rt.StopContainer(context.Background(), ctrID, 2*time.Second); err != nil {
		t.Fatalf("stop container: %v", err)
	}

	if conn.started != driver.KindContainer+":ctr-1" {
		t.Fatalf("unexpected started marker: %s", conn.started)
	}
	if conn.stopped != driver.KindContainer+":ctr-1:2" {
		t.Fatalf("unexpected stopped marker: %s", conn.stopped)
	}
}

type pullOnlyConn struct{}

func (c *pullOnlyConn) Close() error { return nil }
func (c *pullOnlyConn) Pull(ctx context.Context, reference string, options map[string]string) (string, error) {
	return "img", nil
}

func TestUnsupportedOperation(t *testing.T) {
	name := uniqueDriverName(t)
	oci.Register(name, &fakeDriver{conn: &pullOnlyConn{}})

	rt, err := oci.Open(name, "unix:///tmp/fake.sock")
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer rt.Close()

	if err := rt.Ping(context.Background()); err == nil {
		t.Fatalf("expected unsupported ping error")
	}
}
