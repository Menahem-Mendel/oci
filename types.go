package oci

type Method string

const (
	PULL    Method = "PULL"
	PUSH    Method = "PUSH"
	INSPECT Method = "INSPECT"
	EXEC    Method = "EXEC"
	RUN     Method = "RUN"
)
