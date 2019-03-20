package files

import "time"

// Artifact represents an artifact of a test
type Artifact struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// NewArtifact creates a new artifact
func NewArtifact(name, desc string) (Artifact, error) {
	return Artifact{
		Name:        name,
		Description: desc,
		CreatedAt:   now(),
	}, nil
}
