package image

// Manifest is an OCI image manifest.
type Manifest struct {
	Versioned
	MediaType    string            `json:"mediaType,omitempty"`
	ArtifactType string            `json:"artifactType,omitempty"`
	Config       Descriptor        `json:"config"`
	Layers       []Descriptor      `json:"layers"`
	Subject      *Descriptor       `json:"subject,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}
