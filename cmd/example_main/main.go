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
	imageConf, _ := image.Stat(context.Background(), conn, imageID)
	_ = imageConf

	networkID, _ := net.Create(context.Background(), conn)
	namespaceID, _ := namespace.Create(context.Background(), conn)

	// create container
	conf = container.NewConf(
		rt,
		container.WithCDIDevice("path/to/gpu"),
		container.WithCPUPROCSLimit(4),
		container.WithRAMMBLimit(2048),
		container.WithImage(imageID),
		container.WithNetwork(networkID),
		container.WithNamespace(namespaceID),
	)

	containerID, _ := container.Create(context.Background(), conn, cconf)
	cntr, _ := container.Stat(context.Background(), conn, containerID)

	wc, _ := cntr.StdinPipe(context.Background())
	_ = wc

	// create image
	dockerfile, err := os.Open("path/to/Dockerfile")
	if err != nil {
		panic(err)
	}
	defer dockerfile.Close()

	conf = image.NewConf(
		image.WithArch("amd64"),
		image.WithOS("linux"),
		image.WithBase("alpine:latest"),
	)

	imageID, err = image.Create(context.Background(), conn, dockerfile, iconf)
	if err != nil {
		panic(err)
	}
}
