package files

import (
	"testing"
	"time"
)

func TestNewMetadata(t *testing.T) {
	currentTime := time.Now()
	now = func() time.Time {
		return currentTime
	}

	assertName := "TR0001"
	assertDescription := "Desc"
	assertRevision := 1
	assertCreatedAt := currentTime

	metadata, err := NewMetadata("TR0001", "Desc")
	if err != nil {
		t.Error(err)
	}

	if metadata.Name != assertName {
		t.Error("metadata name is wrong")
	}

	if metadata.Description != assertDescription {
		t.Error("metadata description is wrong")
	}

	if metadata.Revision != assertRevision {
		t.Error("metadata revision is wrong")
	}

	if metadata.CreatedAt != assertCreatedAt {
		t.Error("metadata createdAt is wrong")
	}
}

func TestConvertMetadataToJSON(t *testing.T) {
	currentTime := time.Now()
	now = func() time.Time {
		return currentTime
	}

	jsonTime, err := now().MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	assert := "{\"name\":\"TR0001\",\"description\":\"My description!\",\"revision\":1,\"createdAt\":" + string(jsonTime) + "}"

	metadata, err := NewMetadata("TR0001", "My description!")
	if err != nil {
		t.Error(err)
	}

	json, err := metadata.ToJSON()
	if err != nil {
		t.Error(err)
	}

	jsonString := string(json)

	if jsonString != assert {
		t.Error("metadata convertion is wrong", jsonString, assert)
	}
}
