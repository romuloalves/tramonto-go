package entities

import "time"

// Artifact represents an artifact of a test
type Artifact struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CreatedAt   time.Time           `json:"createdAt"`
	Hash        string              `json:"hash"`
	Headers     map[string][]string `json:"headers"`
}

// NewArtifact creates a new artifact
func NewArtifact(name, desc, hash string, headers map[string][]string) (Artifact, error) {
	return Artifact{
		Name:        name,
		Description: desc,
		Hash:        hash,
		CreatedAt:   now(),
		Headers:     headers,
	}, nil
}
