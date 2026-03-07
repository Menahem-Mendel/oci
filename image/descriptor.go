package image

// Digest is a content digest in the form "algorithm:encoded".
type Digest string

// Versioned is the common OCI schema version envelope.
type Versioned struct {
	SchemaVersion int `json:"schemaVersion"`
}

// Descriptor describes referenced OCI content.
type Descriptor struct {
	MediaType    string            `json:"mediaType"`
	Digest       Digest            `json:"digest"`
	Size         int64             `json:"size"`
	URLs         []string          `json:"urls,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	Data         []byte            `json:"data,omitempty"`
	Platform     *Platform         `json:"platform,omitempty"`
	ArtifactType string            `json:"artifactType,omitempty"`
}

// Platform describes the target platform.
type Platform struct {
	Architecture string   `json:"architecture"`
	OS           string   `json:"os"`
	OSVersion    string   `json:"os.version,omitempty"`
	OSFeatures   []string `json:"os.features,omitempty"`
	Variant      string   `json:"variant,omitempty"`
}

var DescriptorEmptyJSON = Descriptor{
	MediaType: MediaTypeEmptyJSON,
	Digest:    "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
	Size:      2,
	Data:      []byte("{}"),
}
