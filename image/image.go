package image

import (
	"fmt"
	"strings"
)

// Reference represents an OCI image reference in a normalized split form.
type Reference struct {
	Registry   string
	Repository string
	Tag        string
	Digest     Digest
}

// ParseReference performs light-weight parsing for common image references.
// It does not validate every edge case in distribution/reference grammar.
func ParseReference(ref string) (Reference, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Reference{}, fmt.Errorf("image: empty reference")
	}

	var out Reference

	if at := strings.Index(ref, "@"); at >= 0 {
		out.Digest = Digest(ref[at+1:])
		ref = ref[:at]
	}

	lastSlash := strings.LastIndex(ref, "/")
	lastColon := strings.LastIndex(ref, ":")
	if lastColon > lastSlash {
		out.Tag = ref[lastColon+1:]
		ref = ref[:lastColon]
	}
	if out.Tag == "" && out.Digest == "" {
		out.Tag = "latest"
	}

	parts := strings.Split(ref, "/")
	if len(parts) == 1 {
		out.Repository = parts[0]
		return out, nil
	}

	first := parts[0]
	if strings.Contains(first, ".") || strings.Contains(first, ":") || first == "localhost" {
		out.Registry = first
		out.Repository = strings.Join(parts[1:], "/")
	} else {
		out.Repository = ref
	}

	if out.Repository == "" {
		return Reference{}, fmt.Errorf("image: invalid reference %q", ref)
	}
	return out, nil
}

func (r Reference) String() string {
	repo := r.Repository
	if r.Registry != "" {
		repo = r.Registry + "/" + repo
	}
	if r.Digest != "" {
		return repo + "@" + string(r.Digest)
	}
	if r.Tag != "" {
		return repo + ":" + r.Tag
	}
	return repo
}
