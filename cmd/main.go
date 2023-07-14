package main

import (
	"context"
	"oci"
	"oci/containers"
	"oci/images"
	"oci/namespaces"
	"oci/networks"
	"os"

	_ "oci/pkg/podman"
)

func main() {
	rt, err := oci.Runtime("podman")
	if err != nil {
		panic(err)
	}

	conn, err := rt.Open(context.Background(), "unix:///var/run/podman.sock")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	imageID, err := images.Pull(context.Background(), conn, "docker.io/library/nginx:latest")
	if err != nil {
		panic(err)
	}

	imageConf, err := images.Stat(context.Background(), conn, imageID)
	if err != nil {
		panic(err)
	}
	_ = imageConf

	networkID, err := networks.New(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	namespaceID, err := namespaces.New(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	options, err := containers.NewConf(
		rt,
		containers.WithCDIDevice("path/to/gpu"),
		containers.WithCPUPROCSLimit(4),
		containers.WithRAMMBLimit(2048),
		containers.WithImage(imageID),
		containers.WithNetwork(networkID),
		containers.WithNamespace(namespaceID),
	)
	if err != oci.ErrUnsupportedConf && err != nil {
		panic(err)
	}

	container, err := containers.New(context.Background(), conn, options)
	if err != nil {
		panic(err)
	}

	container, err := containers.Stat(context.Background(), conn, containerID)
	if err != nil {
		panic(err)
	}

	wc, err := container.StdinPipe(context.Background())
	if err != nil {
		panic(err)
	}
	_ = wc

	dockerfile, err := os.Open("path/to/Dockerfile")
	if err != nil {
		panic(err)
	}
	defer dockerfile.Close()

	iopts, err := images.NewConf(
		rt,
		dockerfile,
		images.WithArch("amd64"),
		images.WithOS("linux"),
		images.WithBase("alpine:latest"),
	)

	imageID, err = images.New(context.Background(), conn, iopts)
}

// 	chain := client.NewChain().
// 		PullImage("docker.io/library/nginx:latest").
// 		NewNetwork().
// 		NewNamespace().
// 		NewContainer()
// 	if err := chain.Commit(context.Background()); err != nil {
// 		panic(err)
// 	}
// }
