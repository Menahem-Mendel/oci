/*
Package oci - Image

Author: Menahem-Mendel Gelfand

  Copyright: Copyright 2023, Menahem-Mendel Gelfand

License: This source code is licensed under the BSD 3-Clause License. You may obtain a copy of the License at:
https://opensource.org/licenses/BSD-3-Clause

This file contains definitions related to OCI images in the oci package.
An Image represents a specific OCI image and provides methods for interacting with it.
*/

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
