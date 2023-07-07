package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"oci"
	"oci/driver"
	"strconv"
	"time"

	"github.com/containers/image/v5/types"
	"github.com/containers/podman/v4/pkg/auth"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/opencontainers/go-digest"
)

func init() {
	oci.Register("podman", &Podman{})
}

type Podman struct {
	conns map[string]driver.Conn

	imageService *imageService
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

	return &Conn{}, nil
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
	if err := attach(c.conn.conn, c.id, pr, nil, nil); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pw, nil
}

func (c *Container) StdoutPipe() (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	if err := attach(c.conn.conn, c.id, nil, pw, nil); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pr, nil
}

func (c *Container) StderrPipe() (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	if err := attach(c.conn.conn, c.id, nil, nil, pw); err != nil {
		pr.Close()
		pw.Close()
		return nil, err
	}
	return pr, nil
}

type Conn struct {
	conn context.Context
	uri  string
	perm byte
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) Begin(ctx context.Context) error {
	return nil
}

// func (c *Conn) Prepare(service string) (io.ReadWriteCloser, driver.ParserFunc, error) {
// 	return nil, nil, nil
// }

type imageService struct {
	name string

	conn *bindings.Connection
}

func (i *imageService) Pull(ctx context.Context, ref string) (int, error) {
	opts := images.PullOptions{}

	if i.conn == nil {
		return 0, fmt.Errorf("podman: %w", "ErrNoConnection")
	}

	params, err := opts.ToParams()
	if err != nil {
		return 0, err
	}
	params.Set("reference", ref)

	// SkipTLSVerify is special.  It's not being serialized by ToParams()
	// because we need to flip the boolean.
	if opts.SkipTLSVerify != nil {
		params.Set("tlsVerify", strconv.FormatBool(!opts.GetSkipTLSVerify()))
	}

	header, err := auth.MakeXRegistryAuthHeader(&types.SystemContext{AuthFilePath: opts.GetAuthfile()}, opts.GetUsername(), opts.GetPassword())
	if err != nil {
		return 0, err
	}

	response, err := i.conn.DoRequest(ctx, nil, http.MethodPost, "/images/pull", params, header)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if !response.IsSuccess() {
		return 0, response.Process(err)
	}

	dec := json.NewDecoder(response.Body)
	var pullErrors []error

	var n int
LOOP:
	for {
		var report entities.ImagePullReport
		if err := dec.Decode(&report); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			report.Error = err.Error() + "\n"
		}

		select {
		case <-response.Request.Context().Done():
			break LOOP
		default:
			// non-blocking select
		}

		switch {
		case report.Stream != "":
			n = len(report.Stream)
		case report.Error != "":
			pullErrors = append(pullErrors, errors.New(report.Error))
		case report.ID != "":
		default:
			return 0, fmt.Errorf("failed to parse pull results stream, unexpected input: %v", report)
		}
	}

	return n, errorhandling.JoinErrors(pullErrors)
}

type Status struct {
	id     string
	ref    string
	size   int64
	digest digest.Digest
	os     string
	arch   string
}

func (s Status) ID() string {
	return s.id
}

func (s Status) Ref() string {
	return s.ref
}

func (s Status) Size() int64 {
	return s.size
}

func (s Status) Digest() digest.Digest {
	return s.digest
}

func (s Status) OS() string {
	return s.os
}

func (s Status) Arch() string {
	return s.arch
}
