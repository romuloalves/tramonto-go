package tramonto

import (
	"errors"

	"gitlab.com/tramonto-one/go-tramonto/entities"

	"gitlab.com/tramonto-one/go-tramonto/db"

	oneDb "gitlab.com/tramonto-one/go-tramonto/db"
	oneHttp "gitlab.com/tramonto-one/go-tramonto/http"
	oneIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// TramontoOne represents the Tramonto One lib
type TramontoOne struct {
	ipfs *oneIpfs.OneIPFS
	db   *db.OneSQLite
	http *oneHttp.OneHTTP
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
		return nil, errors.New("Error initializing OneSQLite: " + err.Error())
	}

	// Configures HTTP server
	http, err := oneHttp.InitializeHTTPServer()
	if err != nil {
		return nil, errors.New("Error initializing OneHTTP: " + err.Error())
	}

	tramontoOne := &TramontoOne{
		ipfs: ipfs,
		db:   db,
		http: http,
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

	// Configures endpoints
	one.http.AddGetArtifact(func(ipns, artifactHash string) (entities.Artifact, []byte, error) {
		return one.GetArtifact(ipns, artifactHash)
	})

	one.http.AddPostArtifact(func(ipns, name, description string, file []byte, fileHeaders map[string][]string) ([]byte, error) {
		return one.AddArtifact(ipns, name, description, file, fileHeaders)
	})

	// Starts HTTP server
	go func() {
		if err := one.http.Start(); err != nil {
			panic(err)
		}
	}()

	return nil
}
