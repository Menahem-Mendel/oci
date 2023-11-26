package container

import (
	"oci/driver"
	"oci/image"
)

type Container struct {
	conf Conf

	stdin  driver.StdinWriter
	stdout driver.StdoutReader
	stderr driver.StderrReader
}

type Conf struct {
	Image image.Conf
}
