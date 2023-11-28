package image

import (
	"time"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
)

const (
	// AnnotationCreated is the annotation key for the date and time on which the image was built (date-time string as defined by RFC 3339).
	AnnotationCreated = "org.opencontainers.image.created"

	// AnnotationAuthors is the annotation key for the contact details of the people or organization responsible for the image (freeform string).
	AnnotationAuthors = "org.opencontainers.image.authors"

	// AnnotationURL is the annotation key for the URL to find more information on the image.
	AnnotationURL = "org.opencontainers.image.url"

	// AnnotationDocumentation is the annotation key for the URL to get documentation on the image.
	AnnotationDocumentation = "org.opencontainers.image.documentation"

	// AnnotationSource is the annotation key for the URL to get source code for building the image.
	AnnotationSource = "org.opencontainers.image.source"

	// AnnotationVersion is the annotation key for the version of the packaged software.
	// The version MAY match a label or tag in the source code repository.
	// The version MAY be Semantic versioning-compatible.
	AnnotationVersion = "org.opencontainers.image.version"

	// AnnotationRevision is the annotation key for the source control revision identifier for the packaged software.
	AnnotationRevision = "org.opencontainers.image.revision"

	// AnnotationVendor is the annotation key for the name of the distributing entity, organization or individual.
	AnnotationVendor = "org.opencontainers.image.vendor"

	// AnnotationLicenses is the annotation key for the license(s) under which contained software is distributed as an SPDX License Expression.
	AnnotationLicenses = "org.opencontainers.image.licenses"

	// AnnotationRefName is the annotation key for the name of the reference for a target.
	// SHOULD only be considered valid when on descriptors on `index.json` within image layout.
	AnnotationRefName = "org.opencontainers.image.ref.name"

	// AnnotationTitle is the annotation key for the human-readable title of the image.
	AnnotationTitle = "org.opencontainers.image.title"

	// AnnotationDescription is the annotation key for the human-readable description of the software packaged in the image.
	AnnotationDescription = "org.opencontainers.image.description"

	// AnnotationBaseImageDigest is the annotation key for the digest of the image's base image.
	AnnotationBaseImageDigest = "org.opencontainers.image.base.digest"

	// AnnotationBaseImageName is the annotation key for the image reference of the image's base image.
	AnnotationBaseImageName = "org.opencontainers.image.base.name"
)
const (
	// MediaTypeDescriptor specifies the media type for a content descriptor.
	MediaTypeDescriptor = "application/vnd.oci.descriptor.v1+json"

	// MediaTypeLayoutHeader specifies the media type for the oci-layout.
	MediaTypeLayoutHeader = "application/vnd.oci.layout.header.v1+json"

	// MediaTypeImageIndex specifies the media type for an image index.
	MediaTypeImageIndex = "application/vnd.oci.image.index.v1+json"

	// MediaTypeImageManifest specifies the media type for an image manifest.
	MediaTypeImageManifest = "application/vnd.oci.image.manifest.v1+json"

	// MediaTypeImageConfig specifies the media type for the image configuration.
	MediaTypeImageConfig = "application/vnd.oci.image.config.v1+json"

	// MediaTypeEmptyJSON specifies the media type for an unused blob containing the value "{}".
	MediaTypeEmptyJSON = "application/vnd.oci.empty.v1+json"
)

const (
	// MediaTypeImageLayer is the media type used for layers referenced by the manifest.
	MediaTypeImageLayer = "application/vnd.oci.image.layer.v1.tar"

	// MediaTypeImageLayerGzip is the media type used for gzipped layers
	// referenced by the manifest.
	MediaTypeImageLayerGzip = "application/vnd.oci.image.layer.v1.tar+gzip"

	// MediaTypeImageLayerZstd is the media type used for zstd compressed
	// layers referenced by the manifest.
	MediaTypeImageLayerZstd = "application/vnd.oci.image.layer.v1.tar+zstd"
)

// Non-distributable layer media-types.
//
// Deprecated: Non-distributable layers are deprecated, and not recommended
// for future use. Implementations SHOULD NOT produce new non-distributable
// layers.
// https://github.com/opencontainers/image-spec/pull/965
const (
	// MediaTypeImageLayerNonDistributable is the media type for layers referenced by
	// the manifest but with distribution restrictions.
	//
	// Deprecated: Non-distributable layers are deprecated, and not recommended
	// for future use. Implementations SHOULD NOT produce new non-distributable
	// layers.
	// https://github.com/opencontainers/image-spec/pull/965
	MediaTypeImageLayerNonDistributable = "application/vnd.oci.image.layer.nondistributable.v1.tar"

	// MediaTypeImageLayerNonDistributableGzip is the media type for
	// gzipped layers referenced by the manifest but with distribution
	// restrictions.
	//
	// Deprecated: Non-distributable layers are deprecated, and not recommended
	// for future use. Implementations SHOULD NOT produce new non-distributable
	// layers.
	// https://github.com/opencontainers/image-spec/pull/965
	MediaTypeImageLayerNonDistributableGzip = "application/vnd.oci.image.layer.nondistributable.v1.tar+gzip"

	// MediaTypeImageLayerNonDistributableZstd is the media type for zstd
	// compressed layers referenced by the manifest but with distribution
	// restrictions.
	//
	// Deprecated: Non-distributable layers are deprecated, and not recommended
	// for future use. Implementations SHOULD NOT produce new non-distributable
	// layers.
	// https://github.com/opencontainers/image-spec/pull/965
	MediaTypeImageLayerNonDistributableZstd = "application/vnd.oci.image.layer.nondistributable.v1.tar+zstd"
)

const (
	// ImageLayoutFile is the file name containing ImageLayout in an OCI Image Layout
	ImageLayoutFile = "oci-layout"
	// ImageLayoutVersion is the version of ImageLayout
	ImageLayoutVersion = "1.0.0"
	// ImageIndexFile is the file name of the entry point for references and descriptors in an OCI Image Layout
	ImageIndexFile = "index.json"
	// ImageBlobsDir is the directory name containing content addressable blobs in an OCI Image Layout
	ImageBlobsDir = "blobs"
)

// Image is the JSON structure which describes some basic information about the image.
// This provides the `application/vnd.oci.image.config.v1+json` mediatype when marshalled to JSON.
type Image struct {
	// Created is the combined date and time at which the image was created, formatted as defined by RFC 3339, section 5.6.
	Created *time.Time `json:"created,omitempty"`

	// Author defines the name and/or email address of the person or entity which created and is responsible for maintaining the image.
	Author string `json:"author,omitempty"`

	// Platform describes the platform which the image in the manifest runs on.
	Platform

	// Config defines the execution parameters which should be used as a base when running a container using the image.
	Config ImageConfig `json:"config,omitempty"`

	// RootFS references the layer content addresses used by the image.
	RootFS RootFS `json:"rootfs"`

	// History describes the history of each layer.
	History []History `json:"history,omitempty"`
}

// ImageLayout is the structure in the "oci-layout" file, found in the root
// of an OCI Image-layout directory.
type ImageLayout struct {
	Version string `json:"imageLayoutVersion"`
}

// ImageConfig defines the execution parameters which should be used as a base when running a container using an image.
type ImageConfig struct {
	// User defines the username or UID which the process in the container should run as.
	User string `json:"User,omitempty"`

	// ExposedPorts a set of ports to expose from a container running this image.
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`

	// Env is a list of environment variables to be used in a container.
	Env []string `json:"Env,omitempty"`

	// Entrypoint defines a list of arguments to use as the command to execute when the container starts.
	Entrypoint []string `json:"Entrypoint,omitempty"`

	// Cmd defines the default arguments to the entrypoint of the container.
	Cmd []string `json:"Cmd,omitempty"`

	// Volumes is a set of directories describing where the process is likely write data specific to a container instance.
	Volumes map[string]struct{} `json:"Volumes,omitempty"`

	// WorkingDir sets the current working directory of the entrypoint process in the container.
	WorkingDir string `json:"WorkingDir,omitempty"`

	// Labels contains arbitrary metadata for the container.
	Labels map[string]string `json:"Labels,omitempty"`

	// StopSignal contains the system call signal that will be sent to the container to exit.
	StopSignal string `json:"StopSignal,omitempty"`

	// ArgsEscaped
	//
	// Deprecated: This field is present only for legacy compatibility with
	// Docker and should not be used by new image builders.  It is used by Docker
	// for Windows images to indicate that the `Entrypoint` or `Cmd` or both,
	// contains only a single element array, that is a pre-escaped, and combined
	// into a single string `CommandLine`. If `true` the value in `Entrypoint` or
	// `Cmd` should be used as-is to avoid double escaping.
	// https://github.com/opencontainers/image-spec/pull/892
	ArgsEscaped bool `json:"ArgsEscaped,omitempty"`
}

// RootFS describes a layer content addresses
type RootFS struct {
	// Type is the type of the rootfs.
	Type string `json:"type"`

	// DiffIDs is an array of layer content hashes (DiffIDs), in order from bottom-most to top-most.
	DiffIDs []digest.Digest `json:"diff_ids"`
}

// History describes the history of a layer.
type History struct {
	// Created is the combined date and time at which the layer was created, formatted as defined by RFC 3339, section 5.6.
	Created *time.Time `json:"created,omitempty"`

	// CreatedBy is the command which created the layer.
	CreatedBy string `json:"created_by,omitempty"`

	// Author is the author of the build point.
	Author string `json:"author,omitempty"`

	// Comment is a custom message set when creating the layer.
	Comment string `json:"comment,omitempty"`

	// EmptyLayer is used to mark if the history item created a filesystem diff.
	EmptyLayer bool `json:"empty_layer,omitempty"`
}

type Descriptor struct {
	// MediaType is the media type of the object this schema refers to.
	MediaType string `json:"mediaType"`

	// Digest is the digest of the targeted content.
	Digest digest.Digest `json:"digest"`

	// Size specifies the size in bytes of the blob.
	Size int64 `json:"size"`

	// URLs specifies a list of URLs from which this object MAY be downloaded
	URLs []string `json:"urls,omitempty"`

	// Annotations contains arbitrary metadata relating to the targeted content.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Data is an embedding of the targeted content. This is encoded as a base64
	// string when marshalled to JSON (automatically, by encoding/json). If
	// present, Data can be used directly to avoid fetching the targeted content.
	Data []byte `json:"data,omitempty"`

	// Platform describes the platform which the image in the manifest runs on.
	//
	// This should only be used when referring to a manifest.
	Platform *Platform `json:"platform,omitempty"`

	// ArtifactType is the IANA media type of this artifact.
	ArtifactType string `json:"artifactType,omitempty"`
}

// Platform describes the platform which the image in the manifest runs on.
type Platform struct {
	// Architecture field specifies the CPU architecture, for example
	// `amd64` or `ppc64le`.
	Architecture string `json:"architecture"`

	// OS specifies the operating system, for example `linux` or `windows`.
	OS string `json:"os"`

	// OSVersion is an optional field specifying the operating system
	// version, for example on Windows `10.0.14393.1066`.
	OSVersion string `json:"os.version,omitempty"`

	// OSFeatures is an optional field specifying an array of strings,
	// each listing a required OS feature (for example on Windows `win32k`).
	OSFeatures []string `json:"os.features,omitempty"`

	// Variant is an optional field specifying a variant of the CPU, for
	// example `v7` to specify ARMv7 when architecture is `arm`.
	Variant string `json:"variant,omitempty"`
}

// DescriptorEmptyJSON is the descriptor of a blob with content of `{}`.
var DescriptorEmptyJSON = Descriptor{
	MediaType: MediaTypeEmptyJSON,
	Digest:    `sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a`,
	Size:      2,
	Data:      []byte(`{}`),
}

// Index references manifests for various platforms.
// This structure provides `application/vnd.oci.image.index.v1+json` mediatype when marshalled to JSON.
type Index struct {
	specs.Versioned

	// MediaType specifies the type of this document data structure e.g. `application/vnd.oci.image.index.v1+json`
	MediaType string `json:"mediaType,omitempty"`

	// ArtifactType specifies the IANA media type of artifact when the manifest is used for an artifact.
	ArtifactType string `json:"artifactType,omitempty"`

	// Manifests references platform specific manifests.
	Manifests []Descriptor `json:"manifests"`

	// Subject is an optional link from the image manifest to another manifest forming an association between the image manifest and the other manifest.
	Subject *Descriptor `json:"subject,omitempty"`

	// Annotations contains arbitrary metadata for the image index.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Manifest provides `application/vnd.oci.image.manifest.v1+json` mediatype structure when marshalled to JSON.
type Manifest struct {
	specs.Versioned

	// MediaType specifies the type of this document data structure e.g. `application/vnd.oci.image.manifest.v1+json`
	MediaType string `json:"mediaType,omitempty"`

	// ArtifactType specifies the IANA media type of artifact when the manifest is used for an artifact.
	ArtifactType string `json:"artifactType,omitempty"`

	// Config references a configuration object for a container, by digest.
	// The referenced configuration object is a JSON blob that the runtime uses to set up the container.
	Config Descriptor `json:"config"`

	// Layers is an indexed list of layers referenced by the manifest.
	Layers []Descriptor `json:"layers"`

	// Subject is an optional link from the image manifest to another manifest forming an association between the image manifest and the other manifest.
	Subject *Descriptor `json:"subject,omitempty"`

	// Annotations contains arbitrary metadata for the image manifest.
	Annotations map[string]string `json:"annotations,omitempty"`
}
