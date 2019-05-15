package entities

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Metadata represents the metadata file of a test
type Metadata struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Revision    int        `json:"revision,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	Artifacts   []Artifact `json:"artifacts"`
	Members     []Member   `json:"members"`
}

// NewMetadata creates a new Metadata instance
func NewMetadata(name, description string) (Metadata, error) {
	return Metadata{
		Name:        name,
		Description: description,
		Revision:    1,
		CreatedAt:   now(),
		Artifacts:   []Artifact{},
		Members:     []Member{},
	}, nil
}

// MetadataFromJSON returns a metadata from JSON
func MetadataFromJSON(metadataJSON []byte) (Metadata, error) {
	var metadata Metadata
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}

// ToJSON converts the metadata file to a JSON
func (m *Metadata) ToJSON() ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return json, nil
}

// AddArtifact adds a new artifact to the test
func (m *Metadata) AddArtifact(name, description, hash string, fileHeaders map[string][]string) error {
	artifact, err := NewArtifact(name, description, hash, fileHeaders)
	if err != nil {
		return err
	}

	m.Artifacts = append(m.Artifacts, artifact)

	return nil
}

// AddMember adds the new member to the metadata
func (m *Metadata) AddMember(newMember Member) error {
	lowerName, lowerEmail := strings.ToLower(newMember.Name), strings.ToLower(newMember.Email)

	// Validates if a members with this data already exists
	for _, member := range m.Members {
		if strings.ToLower(member.Email) != lowerEmail || strings.ToLower(member.Name) != lowerName {
			continue
		}

		return errors.New("Member with this name or email already exists")
	}

	m.Members = append(m.Members, newMember)

	return nil
}
