package main

import (
	"context"
	"oci"
	"oci/containers"
	"oci/images"
	"oci/namespaces"
	"oci/networks"

	_ "oci/pkg/podman"
)

func main() {
	runtime, err := oci.Open(context.Background(), "podman", "unix:///var/run/podman.sock")
	if err != nil {
		panic(err)
	}
	defer runtime.Close()

	imageID, err := images.Pull(context.Background(), runtime, "docker.io/library/nginx:latest")
	if err != nil {
		panic(err)
	}

	imageConf, err := images.Stat(context.Background(), runtime, imageID)
	if err != nil {
		panic(err)
	}

	networkID, err := networks.New(context.Background(), runtime)
	if err != nil {
		panic(err)
	}

	namespaceID, err := namespaces.New(context.Background(), runtime)
	if err != nil {
		panic(err)
	}

	containerID, err := containers.New(context.Background(), runtime, imageID, networkID, namespaceID)
	if err != nil {
		panic(err)
	}

	container, err := containers.Container(context.Background(), runtime, containerID)
	if err != nil {
		panic(err)
	}

	wc, err := container.StdinPipe(context.Background())
	if err != nil {
		panic(err)
	}
	_ = wc
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
