package podman

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"oci/driver"
	"testing"
)

func TestPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_ping" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer ts.Close()

	d := &Driver{}
	conn, err := d.Open(ts.URL)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer conn.Close()

	if err := conn.(*Conn).Ping(context.Background()); err != nil {
		t.Fatalf("ping: %v", err)
	}
}

func TestPull(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/images/create":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}\n"))
		case "/images/library%2Falpine:latest/json":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"Id":          "sha256:abc",
				"RepoTags":    []string{"library/alpine:latest"},
				"RepoDigests": []string{"library/alpine@sha256:def"},
				"Size":        123,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	d := &Driver{}
	conn, err := d.Open(ts.URL)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer conn.Close()

	id, err := conn.(*Conn).Pull(context.Background(), "library/alpine:latest", map[string]string{})
	if err != nil {
		t.Fatalf("pull: %v", err)
	}
	if id != "sha256:abc" {
		t.Fatalf("unexpected id: %s", id)
	}

	doc, err := conn.(*Conn).Inspect(context.Background(), driver.KindImage, "library/alpine:latest")
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	if doc["digest"] != "sha256:def" {
		t.Fatalf("unexpected digest: %#v", doc)
	}
}
