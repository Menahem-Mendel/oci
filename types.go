/*
Package oci - Types

Author: Menahem-Mendel Gelfand

Copyright: Copyright 2023, Menahem-Mendel Gelfand

License: This source code is licensed under the BSD 3-Clause License. You may obtain a copy of the License at:
https://opensource.org/licenses/BSD-3-Clause

This file contains various type definitions used in the oci package.
*/

package oci

type Method string

const (
	PULL    Method = "PULL"
	PUSH    Method = "PUSH"
	INSPECT Method = "INSPECT"
	EXEC    Method = "EXEC"
	RUN     Method = "RUN"
)
