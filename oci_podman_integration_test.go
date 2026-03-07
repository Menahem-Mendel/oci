//go:build integration

package oci_test

import (
	"context"
	"oci"
	_ "oci/pkg/podman"
	"os"
	"testing"
	"time"
)

func TestPodmanPingIntegration(t *testing.T) {
	dsn := os.Getenv("PODMAN_DSN")
	if dsn == "" {
		t.Skip("PODMAN_DSN is not set")
	}

	rt, err := oci.Open("podman", dsn)
	if err != nil {
		t.Fatalf("open runtime: %v", err)
	}
	defer rt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rt.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
}
