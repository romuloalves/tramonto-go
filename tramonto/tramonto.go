package tramonto

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"sync"

	oneIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// One represents the Tramonto One lib
type One struct {
	repoPath string
	node     *core.IpfsNode
	mux      *sync.Mutex
	// database *tramontoDb.OneSQLite
}

// NewTramontoOne returns a new instance of Tramonto One library
func NewTramontoOne(path string) (*One, error) {
	tramontoOne := &One{
		repoPath: path,
		mux:      &sync.Mutex{},
	}

	// Initializes the Tramonto One database
	// db, err := tramontoDb.CreateOneSQLite(path)
	// if err != nil {
	// 	return tramontoOne, err
	// }

	// tramontoOne.database = db

	return tramontoOne, nil
}

// Start starts the node
func (t *One) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if initialized := oneIpfs.IsRepoInitialized(t.repoPath); !initialized {
		return errors.New("Repo not initialized")
	}

	// Loads the plugins before create the node
	oneIpfs.LoadPlugins(t.repoPath)

	// Opens the repo
	nodeRepo, err := oneIpfs.OpenRepo(t.repoPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Creates the new node
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Repo:   nodeRepo,
		Online: true,
	})
	if err != nil {
		return err
	}

	t.node = node

	return nil
}

// Stop stops the node
func (t *One) Stop() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.node == nil {
		return errors.New("Node or repo not initilized")
	}

	if running, err := oneIpfs.IsNodeRunning(t.node); err != nil || !running {
		return errors.New("Node not running")
	}

	if err := t.node.Close(); err != nil {
		return err
	}

	t.node = nil

	return nil
}
