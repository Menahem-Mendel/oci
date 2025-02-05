# OCI Library

**Status:** Under Development 🚧

## Overview

`oci` is a generic driver interface around OCI (Open Container Initiative) based runtime engines. This library provides an abstraction layer for connecting to and interacting with OCI-compliant runtime engines, offering a unified interface regardless of the underlying engine (e.g., Docker, Podman).

This library simplifies the process of developing containerized applications and services, offering consistent interfaces for common operations like pulling images, managing containers, and executing requests. It allows developers to switch between different runtime engines with minimal code changes.

## Features

- **Interface-Based Design:** Enables loose coupling and easy interchangeability of runtime engines.
- **Comprehensive Operations:** Supports operations including images, containers, namespaces, volumes, networks, etc.
- **OCI Specification Compliance:** Executes requests that align with OCI specifications.
- **Standardized Error Handling:** Provides consistent error handling and response structures.
- **Concurrency Safety:** Designed for safe operation in multi-threaded workloads.
- **Multi-Registry Support:** Connects to multiple registries and third-party container managers (Docker, Podman).

## Getting Started

### Installation

To install the `oci` library, use:

```bash
go get github.com/yourusername/oci
```

### Importing the Library

In your Go application, import the `oci` library and any necessary third-party adapters:

```go
import (
	"context"
	"log"
	"oci"

	_ "thirdparty/adapter/name" // Import your driver
)
```

### Creating a New Client

The following example demonstrates how to create a new client and connect to a runtime engine (e.g., Docker, Podman):

```go
package main

import (
	"context"
	"log"
	"oci"

	_ "thirdparty/adapter/docker" // Import Docker adapter
	_ "thirdparty/adapter/podman" // Import Podman adapter
)

func main() {
	// Create a runtime instance using the Docker driver
	runtime, err := oci.NewRuntime("docker")
	if err != nil {
		log.Fatalf("Error initializing Docker runtime: %v", err)
	}

	// Establish a connection to the Docker runtime engine
	conn, err := oci.Open(runtime, "unix:///var/run/docker.sock") // Change to appropriate connection string for your setup
	if err != nil {
		log.Fatalf("Error opening Docker connection: %v", err)
	}
	defer conn.Close()

	// Create a runtime instance using the Podman driver
	podmanRuntime, err := oci.NewRuntime("podman")
	if err != nil {
		log.Fatalf("Error initializing Podman runtime: %v", err)
	}

	// Establish a connection to the Podman runtime engine
	podmanConn, err := oci.Open(podmanRuntime, "unix:///run/podman/podman.sock") // Change to appropriate connection string for your setup
	if err != nil {
		log.Fatalf("Error opening Podman connection: %v", err)
	}
	defer podmanConn.Close()

	// Example using the connection with Docker
	reg, err := oci.NewRegistry("docker.io")
	if err != nil {
		log.Fatalf("Error opening registry: %v", err)
	}

	// Execute a pull request to the Docker registry
	resp, err := oci.Pull(context.Background(), reg, "alpine", oci.WithTag("3.12"), oci.WithPlatform("linux/amd64"))
	if err != nil {
		log.Fatalf("Error pulling image: %v", err)
	}

	log.Println("Pulled image:", resp)
}
```

### Connecting to Multiple Registries

This example shows how to interact with multiple registries, such as Docker Hub and a private registry, using options provided by OCI image and container specs:

```go
package main

import (
	"context"
	"log"
	"oci"
	"oci/image"

	_ "thirdparty/adapter/docker" // Import Docker adapter
)

func main() {
	// Initialize Docker runtime
	runtime, err := oci.NewRuntime("docker")
	if err != nil {
		log.Fatalf("Error initializing Docker runtime: %v", err)
	}

	// Connect to the Docker engine
	conn, err := oci.Open(runtime, "unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("Error opening Docker connection: %v", err)
	}
	defer conn.Close()

	// Connect to Docker Hub
	dockerHub, err := oci.NewRegistry("docker.io")
	if err != nil {
		log.Fatalf("Error connecting to Docker Hub: %v", err)
	}

	// Connect to a private registry
	privateRegistry, err := oci.NewRegistry("https://myprivateregistry.com")
	if err != nil {
		log.Fatalf("Error connecting to private registry: %v", err)
	}

	// Pull an image from Docker Hub using specific options
	resp, err := oci.Pull(context.Background(), dockerHub, "nginx", image.WithTag("1.19"), image.WithPlatform("linux/amd64"), image.WithDigest("sha256:1234567890abcdef..."))
	if err != nil {
		log.Fatalf("Error pulling image from Docker Hub: %v", err)
	}
	log.Println("Pulled image from Docker Hub:", resp)

	// Push the same image to the private registry
	pushResp, err := oci.Push(context.Background(), privateRegistry, resp.Image, "mynginx", image.WithTag("latest"))
	if err != nil {
		log.Fatalf("Error pushing image to private registry: %v", err)
	}
	log.Println("Pushed image to private registry:", pushResp)
}
```

### Executing Commands with OCI Specification Options

After setting up the client, you can execute OCI-compliant commands using options provided by OCI image and container specs. For example, to pull an image:

```go
// Pull an image from a registry with specific options
resp, err := oci.Pull(context.Background(), reg, "alpine", oci.WithTag("3.12"), oci.WithPlatform("linux/amd64"), oci.WithArchitecture("amd64"))
if err != nil {
	log.Fatalf("Error pulling image: %v", err)
}
log.Println("Pulled image:", resp)
```

### Using Container Specifications

Here's an example showing how to create and run a container using the OCI container specifications:

```go
package main

import (
	"context"
	"log"
	"oci"

	_ "thirdparty/adapter/docker" // Import Docker adapter
)

func main() {
	// Initialize Docker runtime
	runtime, err := oci.NewRuntime("docker")
	if err != nil {
		log.Fatalf("Error initializing Docker runtime: %v", err)
	}

	// Connect to the Docker engine
	conn, err := oci.Open(runtime, "unix:///var/run/docker.sock")
	if err != nil {
		log.Fatalf("Error opening Docker connection: %v", err)
	}
	defer conn.Close()

	// Define container options
	containerOpts := []oci.ContainerOption{
		oci.WithImage("alpine:3.12"),
		oci.WithCommand([]string{"sh", "-c", "echo Hello, World!"}),
		oci.WithNetworkMode("bridge"),
		oci.WithEnv([]string{"MY_ENV_VAR=my_value"}),
	}

	// Create and start the container
	container, err := oci.CreateContainer(context.Background(), conn, "my-container", containerOpts...)
	if err != nil {
		log.Fatalf("Error creating container: %v", err)
	}

	err = container.Start(context.Background())
	if err != nil {
		log.Fatalf("Error starting container: %v", err)
	}

	log.Println("Container started successfully!")
}
```

## Interfaces

### Driver

This interface must be implemented by different runtime engines (adapters). It provides methods for connecting to the runtime engine and handling requests.

### Conn

This interface represents a connection to the runtime engine and must be implemented by the adapters. It provides methods to close the connection and perform various operations.

## Writing Custom Adapters

To support a runtime engine not currently covered by this library, implement the `Driver` and `Conn` interfaces. Below is a basic structure of a custom adapter:

```go
type MyDriver struct {
	// Custom fields here
}

func (d *MyDriver) Open(uri string) (oci.Conn, error) {
	// Implement connection logic for your runtime engine...
}
```

## Contributing

Contributions to `oci` are welcomed! Whether it's bug reports, feature requests, or pull requests, we appreciate all help in improving this library. Please read our [contributing guide](CONTRIBUTING.md) to get started.

## License

`oci` is open-source software licensed under the BSD 3-Clause License. See the [LICENSE](LICENSE) file for more details.