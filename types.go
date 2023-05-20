// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package oci

type Method string

const (
	PULL    Method = "PULL"
	PUSH    Method = "PUSH"
	INSPECT Method = "INSPECT"
	EXEC    Method = "EXEC"
	RUN     Method = "RUN"
)
