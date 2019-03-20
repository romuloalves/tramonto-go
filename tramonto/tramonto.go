package tramonto

import (
	"errors"

	oneIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// TramontoOne represents the Tramonto One lib
type TramontoOne struct {
	ipfs *oneIpfs.OneIPFS
}

// NewTramontoOne returns a new instance of Tramonto One library
func NewTramontoOne(path string) (*TramontoOne, error) {
	// Initializes IPFS
	ipfs, err := oneIpfs.InitializeOneIPFS(path)
	if err != nil {
		return nil, errors.New("Error initializing OneIPFS: " + err.Error())
	}

	tramontoOne := &TramontoOne{
		ipfs: ipfs,
	}

	return tramontoOne, nil
}

// Setup starts everything in the One instance
func (one *TramontoOne) Setup() error {
	// Starts the IPFS repo+node
	if err := one.ipfs.Start(); err != nil {
		return errors.New("Error starting IPFS: " + err.Error())
	}

	return nil
}
