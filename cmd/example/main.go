package main

import (
	"context"
	"fmt"
	"log"
	"oci"
	"time"

	// "oci/driver"
	"oci/container"
	"oci/image"

	// "oci/namespace"
	// "oci/net"
	_ "oci/pkg/podman" // Import the Podman adapter
)

func main() {
	// Initialize runtime
	runtime, err := oci.NewRuntime("podman")
	if err != nil {
		log.Fatalf("Error initializing runtime: %v", err)
	}

	// Open connection to the runtime
	conn, err := oci.Open(runtime, "unix:/var/run/podman.sock")
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer conn.Close()

	runtime.Exec("image pull apache:latest")

	puller := image.NewPuller()
	runtime.Serve(ctx, puller)

	func (r Runtime) Serve(s Service) {
		s(Fetcher)
	}

	puller.Pull(ctx, "apache:latest")

	// Pull an image from database through specific container manager
	// Pull(context.Context, conn Conn, dsn string )

	imgSvc := image.NewService(runtime)
	pulled, err := oci.Pull(ctx, conn, imgSvc, "apache:latest")
	inspected, err := oci.Inspect(ctx, conn, imgSvc, pulled.ID)

	cntrSvc := container.NewService(runtime)
	built, err := oci.Build(ctx, conn, cntrSvc, "my-container", pulled.ID, container.BuildOptions{
		CDIDevice:     "path/to/gpu",
		CPUPROCSLimit: 4,
		RAMMBLimit:    2048,
		ImageID:       imgID,
		NetworkID:     networkID,
		NamespaceID:   namespaceID,
		VolumeID:      volumeID,
	})

	cntr, err := oci.Run(ctx, conn, cntrSvc, build.ID, container.RunOptions{
		Detach: true,
		Env:    []string{"FOO=bar"},
		Cmd:    []string{"sleep", "60"},
	})
	cstat, err := oci.Inspect(ctx, conn, cntrSvc, cntr)
	wc, err := cntr.StdinPipe()
	defer wc.Close()
	fmt.Fprintln(wc, "hellow, world!")
}

// // runOperations orchestrates image, network, namespace, and container operations.
// func runOperations(ctx context.Context, conn driver.Conn) error {
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
// func handleImageOperations(ctx context.Context, conn driver.Conn) (string, error) {
// 	oci.Pull(ctx)

// 	type Service struct {
// 		runtime Runtime
// 		service func()
// 	}

// 	service := conn.Begin(context.Background())
// 	service.Exec(image.Pull, "docker.io/library/nginx:latest")

// 	imgid, err := oci.Pull(ctx, conn, "docker.io/library/nginx:latest")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to pull image: %w", err)
// 	}

// 	data, err := image.Inspect(ctx, imgsrv, imgid)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to inspect image: %w", err)
// 	}
// 	_ = data

// 	oci.Transfer(ctx, imgsrv)

// 	return "", nil
// }

// // setupNetworkNamespaceVolume sets up network, namespace, and volume, returning their IDs.
// func setupNetworkNamespaceVolume(ctx context.Context, conn driver.Conn) (string, string, string, error) {
// 	netsrv := net.NewService(conn)
// 	networkID, err := oci.Build(ctx, netsrv)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build network: %w", err)
// 	}

// 	namespaceBuilder := namespace.NewBuilder(conn)
// 	namespaceID, err := oci.Build(ctx, namespaceBuilder, networkID)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build namespace: %w", err)
// 	}

// 	volumeBuilder := namespace.NewBuilder(conn)
// 	volumeID, err := oci.Build(ctx, volumeBuilder)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to build volume: %w", err)
// 	}

// 	return networkID, namespaceID, volumeID, nil
// }

// // manageContainer creates and runs a container with the specified configuration.
// func manageContainer(ctx context.Context, conn driver.Conn, imgID, networkID, namespaceID, volumeID string) error {
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

// // func main() {
// // 	// Define the Docker socket path
// // 	socketPath := "/var/run/docker.sock"

// // 	// Create a custom HTTP transport that uses a Unix socket
// // 	transport := &http.Transport{
// // 		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
// // 			dialer := net.Dialer{}
// // 			return dialer.DialContext(ctx, "unix", socketPath)
// // 		},
// // 	}

// // 	ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
// // 	// Create an HTTP client with the custom transport
// // 	// Pass context timeout to the client timeout, by extracting the deadline from the context and converting it to duration type.
// // 	t, _ := ctx.Deadline()

// // 	client := &http.Client{
// // 		Transport: transport,
// // 		Timeout:   time.Until(t),
// // 	}

// // 	// Define the Docker API endpoint to list containers
// // 	req, err := http.NewRequest(http.MethodGet, "http://localhost/images/json", nil)
// // 	if err != nil {
// // 		log.Fatalf("Failed to create request: %v", err)
// // 	}

// // 	// Send the HTTP request
// // 	resp, err := client.Do(req)
// // 	if err != nil {
// // 		log.Fatalf("Failed to send request: %v", err)
// // 	}
// // 	defer resp.Body.Close()

// // 	// Print the HTTP status
// // 	fmt.Println("Status:", resp.Status)

// // 	// Print the response body
// // 	body := make([]byte, resp.ContentLength)
// // 	_, _ = resp.Body.Read(body)
// // 	fmt.Println(string(body))
// // }

// func main (

// // Declare HTTP client to dockerhub registry
// // and pull image from dockerhub
// client := &http.Client{}
// req, err := http.NewRequest(http.MethodGet, "https://registry.hub.docker.com/v1/search?q=redis&n=1", nil)
// if err != nil {
// 	log.Fatalf("Failed to create request: %v", err)
// }
// resp, err := client.Do(req)

// if err != nil {
// 	log.Fatalf("Failed to send request: %v", err)
// }
// defer resp.Body.Close()

// bs, err := io.ReadAll(resp.Body)
// if err != nil {
// 	log.Fatalf("Failed to read response body: %v", err)
// }
// fmt.Println(string(bs))
// )
