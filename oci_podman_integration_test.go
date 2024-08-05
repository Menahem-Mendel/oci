// oci_podman_integration_test.go

package oci_test

import (
	"context"
	"log"
	"oci"
	"testing"
	"time"

	_ "oci/pkg/podman"
)

var runtime *oci.Runtime

func init() {
	var err error
	runtime, err = oci.NewRuntime("podman")
	if err != nil {
		log.Fatalf("Error initializing runtime: %v", err)
	}

	if runtime == nil {
		log.Fatalf("No runtime is initialized")
	}
}

func TestOpenPodman(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "podman unix socket connection",
			uri:     "unix:///var/run/podman.sock",
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

			if runtime == nil {
				t.Errorf("runtime is nil")
			}
			conn, err := oci.Open(ctx, runtime, tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error = %v, got %v", tt.wantErr, err)
			} else if conn == nil && tt.wantErr == false {
				t.Errorf("no connection is established")
			}

			defer conn.Close()
		})
	}
}
