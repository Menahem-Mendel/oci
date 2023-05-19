package oci

import (
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewImage() *Image { return new(Image) }

type Image struct {
	v1.Image
}

type ImageConfig v1.ImageConfig

type ImageLayoutj v1.ImageLayout
