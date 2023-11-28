package main

import (
	"context"
	"oci"
	"oci/container"
	"oci/image"
	"oci/namespace"
	"oci/net"
	"os"

	_ "oci/pkg/podman"
)

func main() {
	rt, _ := oci.NewRuntime("podman")

	conn, _ := rt.Open(context.Background(), "unix:///var/run/podman.sock")
	defer conn.Close()

	imageID, _ := image.Pull(context.Background(), conn, "docker.io/library/nginx:latest")
	img, _ := image.Inspect(context.Background(), conn, imageID)
	networkID, _ := net.Build(context.Background(), conn)
	namespaceID, _ := namespace.Build(context.Background(), conn)

	// create container
	cconf := container.Conf{
		CDIDevice:     "path/to/gpu",
		CPUPROCSLimit: 4,
		RAMMBLimit:    2048,
		Image:         imageID,
		Network:       networkID,
		Namespace:     namespaceID,
	}

	cntr, _ := container.Build(context.Background(), conn, cconf)
	cstat, _ := cntr.Stat(context.Background())
	_ = cstat
	wc, _ := cntr.StdinPipe()
	_ = wc

	// create image
	dockerfile, err := os.Open("path/to/Dockerfile")
	if err != nil {
		panic(err)
	}
	defer dockerfile.Close()

	iconf := image.Conf{
		Arch: "amd64",
		OS:   "linux",
		Base: "alpine:latest",
	}

	img, _ = image.Build(context.Background(), conn, dockerfile, iconf)
}
