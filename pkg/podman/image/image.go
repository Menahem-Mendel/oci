package image

import (
	"context"
	"oci/driver"
)

type Puller struct {
}

func (p *Puller) Init(pi driver.Puller) {
	pi = p
}

func (p *Puller) Pull(ctx context.Context, ref string) error {
	return nil
}
