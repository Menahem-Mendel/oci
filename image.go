// Copyright 2023, Menahem-Mendel Gelfand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains definitions related to OCI images.
// An Image represents a specific OCI image and provides methods for interacting with it.
package oci

import (
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Image v1.Image

func NewImage() *Image { return new(Image) }

type ImageConfig v1.ImageConfig

type ImageLayoutj v1.ImageLayout

type Descriptor v1.Descriptor

type History v1.History

type Index v1.Index

type Manifest v1.Manifest

type Platform v1.Platform

type RootFS v1.RootFS
