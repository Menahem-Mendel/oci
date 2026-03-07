package image

// Index references manifests for multiple platforms.
type Index struct {
	Versioned
	MediaType    string            `json:"mediaType,omitempty"`
	ArtifactType string            `json:"artifactType,omitempty"`
	Manifests    []Descriptor      `json:"manifests"`
	Subject      *Descriptor       `json:"subject,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
}
