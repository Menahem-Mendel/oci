package main

import (
	"context"
	"fmt"
	"log"
	"oci"
	_ "oci/pkg/podman"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	rt, err := oci.Open("podman", "unix:///run/podman/podman.sock")
	if err != nil {
		log.Fatalf("open runtime: %v", err)
	}
	defer rt.Close()

	fmt.Printf("driver=%s caps=%v\n", rt.DriverName(), rt.Capabilities())

	if err := rt.Ping(ctx); err != nil {
		log.Fatalf("ping runtime: %v", err)
	}

	imageID, err := rt.PullImage(ctx, "docker.io/library/alpine:latest", map[string]string{
		oci.OptionPlatform: "linux/amd64",
	})
	if err != nil {
		log.Fatalf("pull image: %v", err)
	}
	fmt.Printf("image id=%s\n", imageID)

	containerID, err := rt.CreateContainer(ctx, map[string]any{
		oci.SpecName:    "oci-example",
		oci.SpecImage:   imageID,
		oci.SpecCommand: []string{"sh", "-lc", "sleep 30"},
		oci.SpecEnv: map[string]string{
			"OCI_EXAMPLE": "1",
		},
	})
	if err != nil {
		log.Fatalf("create container: %v", err)
	}
	fmt.Printf("container id=%s\n", containerID)

	if err := rt.StartContainer(ctx, containerID); err != nil {
		log.Fatalf("start container: %v", err)
	}

	execResult, err := rt.ExecInContainer(ctx, containerID, []string{"/bin/sh", "-lc", "echo done"}, map[string]any{
		oci.ExecOptionTTY: false,
	})
	if err != nil {
		log.Fatalf("exec in container: %v", err)
	}

	fmt.Printf("exec exit=%v stdout=%q stderr=%q\n", execResult[oci.ExecResultExitCode], execResult[oci.ExecResultStdout], execResult[oci.ExecResultStderr])

	_ = rt.StopContainer(ctx, containerID, 2*time.Second)
	_ = rt.RemoveContainer(ctx, containerID, true)
}
