package tramonto

import (
	"errors"

	"gitlab.com/tramonto-one/go-tramonto/db"

	oneDb "gitlab.com/tramonto-one/go-tramonto/db"
	oneIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// TramontoOne represents the Tramonto One lib
type TramontoOne struct {
	ipfs *oneIpfs.OneIPFS
	db   *db.OneSQLite
}

// NewTramontoOne returns a new instance of Tramonto One library
func NewTramontoOne(path string) (*TramontoOne, error) {
	// Initializes IPFS
	ipfs, err := oneIpfs.InitializeOneIPFS(path)
	if err != nil {
		return nil, errors.New("Error initializing OneIPFS: " + err.Error())
	}

	// Initializes the database
	db, err := oneDb.OpenOneSQLite(path)
	if err != nil {
		return nil, errors.New("Error opening OneSQLite: " + err.Error())
	}

	tramontoOne := &TramontoOne{
		ipfs: ipfs,
		db:   db,
	}

	return tramontoOne, nil
}

// Setup starts everything in the One instance
func (one *TramontoOne) Setup() error {
	// Initializes the repo if its not
	if err := one.ipfs.InitRepo(); err != nil {
		return errors.New("Error initializing repo: " + err.Error())
	}

	// Starts the IPFS node
	if err := one.ipfs.Start(); err != nil {
		return errors.New("Error starting IPFS: " + err.Error())
	}

	// Migrates database
	if err := one.db.MigrateTables(); err != nil {
		return errors.New("Error migrating IPFS: " + err.Error())
	}

	return nil
}
