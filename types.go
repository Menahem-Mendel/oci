// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package oci

import "errors"

type Method string

const (
	PULL    Method = "PULL"
	PUSH    Method = "PUSH"
	INSPECT Method = "INSPECT"
	EXEC    Method = "EXEC"
	RUN     Method = "RUN"
)

type Kind string

const (
	IMAGE     Kind = "IMAGE"
	CONTAINER Kind = "CONTAINER"
	NETWORK   Kind = "NETWORK"
	POD       Kind = "POD"
)

// ErrConnClosed indicates that the connection to the OCI runtime has been closed.
var ErrConnClosed = errors.New("oci: connection is already closed")

// ErrNoConnection indicates that there is no connection established to the OCI runtime.
var ErrNoConnection = errors.New("oci: no connection established")

// ErrDriverMismatch indicates that the specified driver does not match the OCI runtime.
var ErrDriverMismatch = errors.New("oci: driver mismatch")

// ErrContextCanceled indicates that the context has been canceled.
var ErrContextCanceled = errors.New("oci: context canceled")

// ErrInvalidContainerStatus is returned when an operation is not valid for the current status of the container.
var ErrInvalidContainerStatus = errors.New("oci: invalid container status")

// ErrInvalidContainerStatus is returned when an operation is not valid for the current status of the container.
var ErrHandlerNotFound = errors.New("oci: handler not found")

// ErrInvalidContainerStatus is returned when an operation is not valid for the current status of the container.
var ErrHandlerRedefined = errors.New("oci: handler redefined")

// ErrUnsupportedOperation is returned when an operation is not supported by the OCI runtime service.
var ErrUnsupportedOperation = errors.New("oci: unsupported operation")

// ErrUnsupportedService is returned when a service is not supported by the OCI runtime.
var ErrUnsupportedService = errors.New("oci: unsupported service")

// ErrUnsupportedService is returned when a service is not supported by the OCI runtime.
var ErrUnregisteredDriver = errors.New("oci: unsupported service")

// invalid errors
var (
	// ErrInvalidMethod indicates that an invalid method was specified in a request to the OCI runtime.
	ErrInvalidMethod = errors.New("oci: invalid method")

	// ErrInvalidRef indicates that an invalid reference was specified in a request to the OCI runtime.
	ErrInvalidRef = errors.New("oci: invalid reference")

	// ErrInvalidKind indicates that an invalid kind was specified in a request to the OCI runtime.
	ErrInvalidKind = errors.New("oci: invalid kind")

	// ErrInvalidBody indicates that an invalid body was provided in a request to the OCI runtime.
	ErrInvalidBody = errors.New("oci: invalid body")
)

// not found errors
var (
	// ErrContainerNotFound is returned when the specified container could not be found.
	ErrImageNotFound = errors.New("oci: image not found")

	// ErrContainerNotFound is returned when the specified container could not be found.
	ErrContainerNotFound = errors.New("oci: container not found")
)
