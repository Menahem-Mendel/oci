package adapter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"oci"
	"oci/driver"
	"oci/pkg/podman/image"
	"time"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	// "github.com/containers/podman/v4/pkg/bindings/network"
	// "github.com/containers/podman/v4/pkg/specgen"
)

func init() {
	// var p *Podman
	oci.Register("podman", nil, nil)

	// var imageService *image.Service
	var imgPullerService *image.Puller
	// var containerService container.Service

	oci.Handle(oci.Pull, imgPullerService)

	// oci.HandlePuller(containerService)
}

type Podman struct {
	conns map[string]driver.Conn
}

func (c *Conn) Exec(ctx context.Context, f func(context.Context, any, ...any), params ...any) {
	f()
}

func (p *Podman) Open(ctx context.Context, uri string) (driver.Conn, error) {
	ctx, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("Podman.Connect: %w", err)
	}

	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Podman.Connect: %w", err)
	}
	_ = conn

	return nil, nil
}

type UnixRoundTripper struct {
	path string
}

func (rt *UnixRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	dial := net.Dialer{}
	conn, err := dial.DialContext(req.Context(), "unix", rt.path)
	if err != nil {
		return nil, err
	}

	err = req.Write(conn)
	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Close closes all connections to the podman container manager.
func (p *Podman) Close() error {
	for _, conn := range p.conns {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("Podman.Close: %w", err)
		}
	}

	return nil
}

func attach(conn context.Context, cid string, rci io.ReadCloser, wco io.WriteCloser, wce io.WriteCloser) error {
	attachOpts := containers.AttachOptions{}
	attachOpts.WithStream(true)
	attachOpts.WithLogs(true)

	attachReady := make(chan bool)

	err := containers.Attach(conn, cid, rci, wco, wce, attachReady, &attachOpts)
	if err != nil {
		return err
	}

	select {
	case <-attachReady:
		// Successfully attached
	case <-time.After(time.Second * 5):
		// Timeout after 5 seconds
		return fmt.Errorf("timeout waiting for container attach")
	}

	close(attachReady)
	return nil
}

type Container struct {
	conn Conn

	id string
}

func (c *Container) StdinPipe() (io.WriteCloser, error) {
	pr, pw := io.Pipe()
	if err := attach(nil, c.id, pr, nil, nil); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pw, nil
}

func (c *Container) StdoutPipe() (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	if err := attach(nil, c.id, nil, pw, nil); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pr, nil
}

func (c *Container) StderrPipe() (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	if err := attach(nil, c.id, nil, nil, pw); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pr, nil
}

type Conn struct {
	conn *bindings.Connection
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Begin(ctx context.Context) error {
	return nil
}

func (c *Conn) Prepare(service string) (any, error) {
	return nil, nil
}

func f(p driver.Puller) {

}

func g() {
	var s *image.Service
	f(s)
}
