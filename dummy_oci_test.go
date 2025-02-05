package oci_test

import (
	"context"
	"fmt"
	"net/url"
	"oci"
	"oci/driver"
	"oci/image"
	"sync"
	"testing"
	"time"
)

func init() {
	oci.Register("fake", &fakeDriver{})

	rt, err := oci.NewRuntime("fake")
	if err != nil {
	}
	_, _ = rt, err

}

type fakeDriver struct {
	mu         sync.Mutex
	openCount  int
	closeCount int

	imgSrv fakeIMGServer
	// ctrSrv   fakeCTRServer
	// netSrv   fakeNETServer
	// nmspcSrv fakeNMSPCServer
}

func (fd *fakeDriver) Open(uri string) (driver.Conn, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	if !isValidScheme(u.Scheme) {
		return nil, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}
	return &fakeConn{}, nil
}

type fakeIMGServer struct{}

func (fis *fakeIMGServer) Pull(ctx context.Context, dsn string) (string, error) {
	return "", nil
}

func (fis *fakeIMGServer) Push(ctx context.Context, dsn, id string) error {
	return nil
}

func (fis *fakeIMGServer) Stat(ctx context.Context, id string) (map[string]any, error) {
	return nil, nil
}

// type fakeCTRServer struct{}

// func (fcs fakeCTRServer)

type fakeConn struct{}

func (fc *fakeConn) Close() error {
	return nil
}

func (fc *fakeConn) Begin(ctx context.Context) error {
	return nil
}

type imageService struct{}

type fakeImage struct{}

type fakeContainer struct{}

func (i *imageService) Pull(ctx context.Context, dsn string) (string, error) {
	return "", nil
}

func isValidScheme(scheme string) bool {
	validSchemes := map[string]struct{}{
		"unix": {},
		"tcp":  {},
		"ssh":  {},
	}

	_, valid := validSchemes[scheme]
	return valid
}

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "podman unix socket connection",
			uri:     "unix:///run/user/1000/podman/podman.sock",
			wantErr: false,
		},
		{
			name:    "podman tcp socket connection",
			uri:     "tcp://localhost:2375",
			wantErr: false,
		},
		{
			name:    "podman invalid uri",
			uri:     "invalid://localhost",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			fd := &fakeDriver{}
			conn, err := oci.Open(ctx, fd, tt.uri)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error = %v, got %v", tt.wantErr, err)
			}

			if err == nil && conn == nil {
				t.Fatalf("no connection is established")
			}

			if conn != nil {
				if err := conn.Close(); err != nil {
					t.Errorf("failed to close connection: %v", err)
				}
			}
		})
	}
}

func TestPull(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := oci.Open(ctx, fd, tt.dsn)

	puller := image.NewPuller(conn)

	tests := []struct {
		name    string
		dsn     string
		puller  driver.Puller
		wantErr bool
	}{
		{
			name:    "nginx latest",
			dsn:     "docker.io/library/nginx:latest",
			puller:  puller,
			wantErr: false,
		},
		{
			name:    "nginx latest with no scheme",
			dsn:     "nginx:latest",
			puller:  puller,
			wantErr: false,
		},
		{
			name:    "debian default tag",
			dsn:     "debian",
			puller:  puller,
			wantErr: false,
		},
		{
			name:    "debian bookworm",
			dsn:     "debian:bookworm",
			puller:  puller,
			wantErr: false,
		},
		{
			name:    "ubuntu digest",
			dsn:     "ubuntu@sha256:26c68657ccce2cb0a31b330cb0be2b5e108d467f641c62e13ab40cbec258c68d",
			puller:  puller,
			wantErr: false,
		},
		{
			name:    "pull from custom registry",
			dsn:     "myregistry.local:5000/testing/test-image",
			puller:  puller,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			id, err := oci.Pull(ctx, fDriver.imgSrv, tt.dsn)

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error = %v, got %v", tt.wantErr, err)
			}

			if err == nil && conn == nil {
				t.Fatalf("no connection is established")
			}

			if conn != nil {
				if err := conn.Close(); err != nil {
					t.Errorf("failed to close connection: %v", err)
				}
			}
		})
	}
}
