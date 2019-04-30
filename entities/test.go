package entities

// Test represents a test in the Tramonto One
type Test struct {
	// IPFS hash
	Ipfs string `json:"ipfs"`

	// IPNS hash
	Ipns string `json:"ipns"`

	// If IPNS key was created
	IpnsKeyCreated bool `json:"ipnsKeyCreated"`

	// If the current node is owner of the archieve
	IsOwner bool `json:"isOwner"`

	// Secret to decrypt the files
	Secret string `json:"secret"`

	// Metadata informations
	Metadata Metadata `json:"metadata,omitempty"`
}

// NewEmptyTest instances a new empty test
func NewEmptyTest() Test {
	return Test{
		IpnsKeyCreated: false,
		IsOwner:        true,
	}
}
