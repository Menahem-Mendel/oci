package oci

import "errors"

var (
	ErrUnregisteredDriver = errors.New("oci: unregistered driver")
	ErrRuntimeClosed      = errors.New("oci: runtime is closed")
	ErrNotConnected       = errors.New("oci: runtime is not connected")
	ErrUnsupported        = errors.New("oci: operation not supported by driver")
	ErrInvalidReference   = errors.New("oci: invalid reference")
	ErrInvalidArgument    = errors.New("oci: invalid argument")
)

type Capability string

const (
	CapabilityPing    Capability = "ping"
	CapabilityPull    Capability = "pull"
	CapabilityPush    Capability = "push"
	CapabilityInspect Capability = "inspect"
	CapabilityList    Capability = "list"
	CapabilityRemove  Capability = "remove"
	CapabilityCreate  Capability = "create"
	CapabilityStart   Capability = "start"
	CapabilityStop    Capability = "stop"
	CapabilityExec    Capability = "exec"
)

// Common option/spec/result keys used by runtime and adapters.
const (
	OptionPlatform        = "platform"
	OptionAllTags         = "all_tags"
	OptionInsecureSkipTLS = "insecure_skip_tls_verify"
	OptionUsername        = "username"
	OptionPassword        = "password"

	SpecName        = "name"
	SpecImage       = "image"
	SpecCommand     = "command"
	SpecEnv         = "env"
	SpecLabels      = "labels"
	SpecWorkingDir  = "working_dir"
	SpecNetworkMode = "network_mode"

	ExecOptionEnv   = "env"
	ExecOptionTTY   = "tty"
	ExecOptionStdin = "stdin"

	ExecResultExitCode = "exit_code"
	ExecResultStdout   = "stdout"
	ExecResultStderr   = "stderr"
)

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func cloneAnyMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
