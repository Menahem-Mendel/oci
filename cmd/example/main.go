package main

import (
	"context"
	"log"
	"oci"

	_ "oci/pkg/podman" // Import the Podman adapter
)

func main() {
	ctx := context.Background()

	// Initialize runtime
	runtime, err := oci.NewRuntime("podman")
	if err != nil {
		log.Fatalf("Error initializing runtime: %v", err)
	}

	// Open connection to the runtime
	conn, err := oci.Open(ctx, runtime, "unix:///var/run/podman.sock")
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer conn.Close()

	// if err := runOperations(ctx, conn); err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }
}

// // runOperations orchestrates image, network, namespace, and container operations.
// func runOperations(ctx context.Context, conn oci.Conn) error {
// 	imgID, err := handleImageOperations(ctx, conn)
// 	if err != nil {
// 		return fmt.Errorf("image operations failed: %w", err)
// 	}

// 	networkID, namespaceID, volumeID, err := setupNetworkNamespaceVolume(ctx, conn)
// 	if err != nil {
// 		return fmt.Errorf("network/namespace/volume setup failed: %w", err)
// 	}

// 	if err := manageContainer(ctx, conn, imgID, networkID, namespaceID, volumeID); err != nil {
// 		return fmt.Errorf("container operations failed: %w", err)
// 	}

// 	return nil
// }

// // handleImageOperations pulls and inspects the specified image, returning its ID.
// func handleImageOperations(ctx context.Context, conn oci.Conn) (string, error) {
// 	puller := image.NewPuller(conn)
// 	puller.Passwd("passwd")
// 	puller.AllTags(true)

// 	img, err := image.Pull(ctx, puller, "docker.io/library/nginx:latest")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to pull image: %w", err)
// 	}

// 	imgID, err := img.Inspect(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to inspect image: %w", err)
// 	}

// 	return imgID, nil
// }

// // setupNetworkNamespaceVolume sets up network, namespace, and volume, returning their IDs.
// func setupNetworkNamespaceVolume(ctx context.Context, conn oci.Conn) (string, string, string, error) {
// 	networkBuilder := net.NewBuilder(conn)
// 	networkID, err := oci.Build(ctx, networkBuilder)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build network: %w", err)
// 	}

// 	namespaceBuilder := namespace.NewBuilder(conn)
// 	namespaceID, err := oci.Build(ctx, namespaceBuilder, networkID)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build namespace: %w", err)
// 	}

// 	volumeBuilder := volume.NewBuilder(conn)
// 	volumeID, err := oci.Build(ctx, volumeBuilder)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build volume: %w", err)
// 	}

// 	return networkID, namespaceID, volumeID, nil
// }

// // manageContainer creates and runs a container with the specified configuration.
// func manageContainer(ctx context.Context, conn oci.Conn, imgID, networkID, namespaceID, volumeID string) error {
// 	cconf := container.Conf{
// 		CDIDevice:     "path/to/gpu",
// 		CPUPROCSLimit: 4,
// 		RAMMBLimit:    2048,
// 		ImageID:       imgID,
// 		NetworkID:     networkID,
// 		NamespaceID:   namespaceID,
// 		VolumeID:      volumeID,
// 	}

// 	containerBuilder := container.NewBuilder(conn)
// 	cntr, err := containerBuilder.Build(ctx, cconf)
// 	if err != nil {
// 		return fmt.Errorf("failed to build container: %w", err)
// 	}

// 	cstat, err := cntr.Stat(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to get container status: %w", err)
// 	}
// 	log.Printf("Container status: %v", cstat)

// 	wc, err := cntr.StdinPipe()
// 	if err != nil {
// 		return fmt.Errorf("failed to get stdin pipe: %w", err)
// 	}
// 	defer wc.Close()

// 	if _, err = wc.Write([]byte("hello world\n")); err != nil {
// 		return fmt.Errorf("failed to write to stdin: %w", err)
// 	}
// 	return nil
// }
