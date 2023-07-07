package containers

import (
	"oci/driver"
	"oci/images"
)

type Container struct {
	conf Conf

	stdin  driver.StdinWriter
	stdout driver.StdoutReader
	stderr driver.StderrReader
}

type Conf struct {
	Image images.Conf
}
