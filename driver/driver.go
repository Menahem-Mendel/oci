// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the definitions for the Driver interface and related types in the oci package.
// A Driver is the core interface in this package and is designed to enable interaction with different OCI implementations.
package driver

import (
	"context"
	"io"
	"net"
	"os"
)

// Driver is an interface that defines the behavior of components that
// establish connections with container runtime management daemons. The
// implementations of this interface allow the application to interact with
// different container managers (such as Docker, Podman, etc.), abstracting
// the details of the connection process and providing a unified way of
// managing these connections across different runtimes.
type Driver interface {
	Connector
}

type Connector interface {
	// Open method is responsible for establishing a new connection to a
	// container manager daemon, using the provided context and uri parameters.
	//
	// The context parameter is a context.Context object that can be used
	// to cancel the connection process. This is useful for controlling the
	// lifecycle of the connection, especially in situations where the connection
	// process could potentially block for an indefinite period of time, or
	// in scenarios where control over cancellation, timeout, and deadline
	// behavior is required.
	//
	// The uri parameter specifies the location of the container manager daemon.
	// The format and interpretation of this uri is driver-specific, meaning
	// different Driver implementations may expect different uri formats
	// according to the requirements of the container manager they interact with.
	//
	// The Open method returns a Conn object on success, representing the
	// established connection to the container manager daemon. This Conn object
	// can then be used for further interactions with the container manager daemon.
	//
	// In case the connection process fails, the Open method should return a
	// non-nil error that provides details about the reason of the failure. The
	// error should be descriptive enough to allow callers to understand what went
	// wrong and possibly make informed decisions about error handling and recovery.
	Open(ctx context.Context, uri string) (Conn, error)
}

// Conn represents a connection to a container runtime management daemon.
// This interface abstracts the connection details, allowing the library
// to interact with different container managers in a unified manner,
// regardless of the specific protocols or mechanisms used by each manager
// to maintain and manage connections.
//
// Implementations of the Conn interface are expected to provide specific
// Close and Begin methods that manage the lifecycle of the connection
// (such as initializing, terminating, and cleaning up the connection).
// This is particularly important to prevent resource leaks in long-running
// programs or in scenarios where a large number of connections might be
// created over time.
type Conn interface {
	// Close method should return an error if the closing operation fails,
	// allowing the caller to understand the reason for the failure and potentially
	// take steps to handle or recover from the error condition. The returned
	// error should be descriptive enough to provide meaningful information about
	// the failure.
	//
	// Although not explicitly required by the Conn interface itself, implementations
	// are strongly encouraged to make the Close method idempotent, meaning that
	// multiple calls to Close on the same connection object should be safe and
	// should not result in undefined behavior. Typically, the first call to
	// Close should perform the actual closing operation, and any subsequent
	// calls should have no effect but should return a specific error (or nil)
	// to indicate that the connection has already been closed.
	Close() error

	// Begin method is responsible for initializing the connection to the
	// container manager daemon and preparing it for further interactions.
	// This method should be called before any other operations are performed
	// on the connection after it has been closed.
	//
	// The Begin method returns an error if the initialization process fails.
	// This error should be descriptive enough to provide meaningful information
	// about the reason for the failure, allowing the caller to understand what
	// went wrong and potentially take steps to handle or recover from the error
	// condition.
	//
	// Although not explicitly required by the Conn interface itself, implementations
	// are strongly encouraged to make the Begin method idempotent, meaning that
	// multiple calls to Begin on the same connection object should be safe and
	// should not result in undefined behavior. Typically, the first call to
	// Begin should perform the actual initialization operation, and any subsequent
	// calls should have no effect but should return a specific error (or nil)
	// to indicate that the connection has already been initialized.
	Begin(ctx context.Context) error

	// Prepare method in the Conn interface is used to prepare a service based on the
	// provided service name. It returns an io.ReadWriteCloser that's bound to the connection.
	// This returned object encapsulates the CRD operations that can be performed on the service.
	//
	// For Create (Write) operations:
	// Writing to the io.Writer part of the returned io.ReadWriteCloser can be used for
	// creating resources. For example, in an image service, writing would create an image
	// from the provided data. In the context of a network service, writing might create
	// a network based on provided configurations.
	//
	// For Read operations:
	// The io.Reader part is used for retrieving (inspecting) resources. Reading would
	// provide information about the resource (like an image or network) based on its id
	// or reference.
	//
	// For Delete (Close) operations:
	// The Close method is used for deleting resources. When a resource (like an image or
	// a network) is no longer needed, calling Close() would delete or remove it.
	//
	// If the provided service name does not exist or if the operation fails, the method
	// will return an error.
	//
	// Note: The specific behavior and the kind of data that needs to be written or read
	// depends on the implementation of the specific service.
	// Prepare(service string) (io.ReadWriteCloser, ParserFunc, error)
}

type Listener interface {
	Listen(ctx context.Context, cid string) (net.Conn, error)

	Close() error
}

// type ImageParserFunc func(r io.Reader) (ImageInfo, error)

// type ParserFunc func(r io.Reader) (Info, error)

// func (p ParserFunc) Parse(r io.Reader) (Info, error) {
// 	return p(r)
// }

// type Parser interface {
// 	Parse(r io.Reader) (Info, error)
// }

// Puller interface represents a key abstraction in the realm of container-based and cloud-native
// applications. The term 'pull' is conventionally used to denote the retrieval of an object or resource
// from a remote location. Specifically, within the context of container technology, 'pull' refers to
// the process of fetching a container image from a registry. However, the Puller interface is
// intentionally designed without an explicit prefix (like 'Container' or 'Image'), allowing it to
// be used in diverse contexts where the pull operation is applicable.
type Puller interface {
	// Pull is a method within the Puller interface that encapsulates the logic for fetching or
	// 'pulling' a resource, represented by a string reference, in a given context.
	//
	// The first parameter, 'ctx', is a context.Context object. This carries deadlines, cancellation
	// signals, and other request-specific values across API boundaries and between processes. It is
	// used to manage the lifecycle of the pull operation. For example, if the pulling process is
	// long-running, the context can be used to stop the operation gracefully if it's taking too long
	// or if the client decides to cancel the operation. It's the responsibility of the Pull method's
	// implementation to respect the context's behavior and to regularly check its status during the
	// pulling process.
	//
	// The second parameter, 'reference', is a string that points to the resource to be pulled. In
	// a container environment, this is usually the tag or the identifier of the container image in
	// a remote registry. However, based on the flexible design of this interface, the 'reference'
	// can also refer to other pullable resources, depending on the context where this interface is
	// implemented.
	//
	// The Pull method returns an integer and an error. The integer 'n' represents the number of bytes
	// that have been pulled during the process. This is particularly useful when pulling data from a
	// stream, as it provides the consumer with insights into the amount of data transferred during
	// the pull operation.
	//
	// The error returned by the Pull method indicates the success or failure of the pull operation.
	// If the operation is successful, this value will be nil. If something goes wrong, this error
	// should provide details about the failure, which can include IO errors (issues related to the
	// data stream), context errors (the pull operation being cancelled or timing out), or domain-specific
	// errors (e.g., the resource not being found in the remote registry).
	Pull(ctx context.Context, ref string) (id string, err error)
}

// Pusher is an interface that abstracts the operation of pushing an OCI resource, such as an image,
// to a storage location, often a container registry. Different implementations of this interface
// may work with different OCI compliant runtimes like Docker, Podman, etc., each having their own
// specific methods and mechanisms to push an OCI resource.
type Pusher interface {
	// Push is a method which initiates the process of pushing an OCI resource, identified by its
	// local ID or image name, to a specified location, referenced by 'ref'. The 'id' is a string which uniquely
	// identifies the resource in the local storage of an OCI compliant runtime. This could be an image
	// name or ID. The 'ref' is a reference to the destination where the image needs to be pushed.
	// This could be a URL of a container registry or any other destination supported by the specific runtime.
	//
	// The context parameter is a context.Context object, which is used to provide cancellation
	// and timeout signals. This ensures that the push operation can be controlled and monitored
	// effectively, allowing for graceful termination in case of an unexpected long runtime,
	// system shutdown, or similar scenarios where the operation cannot be allowed to run indefinitely.
	//
	// The Push method returns an integer indicating the number of bytes transferred during
	// the push operation. This might be useful for tracking the amount of data transferred,
	// especially in applications where monitoring or logging of data transfer is required.
	// However, depending on the specific implementation and nature of the operation, this value
	// might not always be accurate or meaningful. Therefore, use of this returned value should
	// be done with understanding of these potential limitations.
	//
	// In case of failure during the push operation, the Push method returns a non-nil error.
	// This error helps the caller to understand what went wrong during the push operation,
	// such as network errors, access permission errors, invalid ID or reference string, etc.
	// Based on the nature and specifics of the error, appropriate error handling and recovery
	// strategies can be implemented.
	Push(ctx context.Context, ref, id string) error
}

type Creator interface {
	Create(ctx context.Context, id string) (resID string, err error)
}

type Starter interface {
	Start(ctx context.Context, id string) error
}

type Stoper interface {
	Stop(ctx context.Context, id string) error
}

type Killer interface {
	Kill(ctx context.Context, id string, signal os.Signal) error
}

type Mounter interface {
	Mount(ctx context.Context, id string, target string) error
}

type HandlerFunc func(ctx context.Context) error

func (h HandlerFunc) ServeOCI(ctx context.Context) error {
	return h(ctx)
}

type Handler interface {
	ServeOCI(ctx context.Context) error
}

type Limiter interface {
	Limit(limit int) error
}

type Option interface {
	Apply(n int) ([]byte, error)
}

// Service could be image, container, network, namespace etc.
type Service interface {
	Inspector
	Lister
	Remover
}

type Lister interface {
	List(ctx context.Context) ([]Inspector, error)
}

type Remover interface {
	Remove(ctx context.Context, id string) error
}

type Inspector interface {
	Stat(ctx context.Context, id string, v any) (b []byte, err error)
}

// Execer is an interface that abstracts the operation of executing commands within the context
// of a container runtime. This could be used with various container runtimes like Docker,
// Podman, containerd, etc., each having their unique methods for executing commands inside running
// containers. The Execer interface encapsulates these operations, allowing the caller to execute a
// command without needing to know the specifics of the underlying runtime.
type Execer interface {
	// Exec is a method that initiates the operation of executing a specific command, identified by
	// the 'cmd' parameter and a variable number of arguments 'args', within a running container
	// or any other suitable runtime environment.
	//
	// The 'ctx' parameter is a context.Context object, which is used to provide control over the
	// lifetime of the operation. With this context, it's possible to cancel the executing operation
	// or set a timeout, thus preventing runaway processes that could consume resources indefinitely.
	// This is particularly important in scenarios where resources are limited, and long-running
	// operations might have a significant impact on system performance.
	//
	// The Exec method returns an error in case the command execution fails. This error provides
	// insight into what went wrong during the operation, such as issues with the command itself,
	// problems with the target container, network errors, access permission errors, etc. This
	// returned error allows the caller to implement appropriate error handling and recovery
	// strategies, depending on the specifics of the error.
	Exec(ctx context.Context, cmd string, args ...string) error
}

// type StdioStreamer interface {
// 	StdinWriter
// 	StdoutReader
// 	StderrReader
// }

type StdinWriter interface {
	StdinPipe(ctx context.Context) (io.WriteCloser, error)
}

type StdoutReader interface {
	StdoutPipe(ctx context.Context) (io.ReadCloser, error)
}

type StderrReader interface {
	StderrPipe(ctx context.Context) (io.ReadCloser, error)
}

// type Stdin struct {
// 	in io.Writer

// 	cancel func()

// 	mu     sync.Mutex
// 	closed bool
// }

// func (s *Stdin) Write(p []byte) (n int, err error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return 0, errors.New("write to closed stdin")
// 	} else if s.in == nil {
// 		return 0, errors.New("stdin not initialized")
// 	}

// 	return s.in.Write(p)
// }

// func (s *Stdin) Close() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return errors.New("stdin already closed")
// 	}

// 	s.closed = true
// 	s.cancel()
// 	return nil
// }

// type Stdout struct {
// 	out io.Reader

// 	cancel func()

// 	mu     sync.Mutex
// 	closed bool
// }

// func (s *Stdout) Read(p []byte) (n int, err error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return 0, errors.New("read from closed stdout")
// 	} else if s.out == nil {
// 		return 0, errors.New("stdout not initialized")
// 	}

// 	return s.out.Read(p)
// }

// func (s *Stdout) Close() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if s.closed {
// 		return errors.New("stdout already closed")
// 	}

// 	s.closed = true
// 	s.cancel()
// 	return nil
// }

// type Stderr struct {
// 	out io.Reader

// 	cancel func()

// 	mu     sync.Mutex
// 	closed bool
// }
