package tramonto

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	"sync"

	tramontoIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// One represents the Tramonto One lib
type One struct {
	repoPath string
	repo     repo.Repo
	node     *core.IpfsNode
	mux      *sync.Mutex
	// database *tramontoDb.OneSQLite
}

// NewTramontoOne returns a new instance of Tramonto One library (automatically initializes the IPFS repo and database)
func NewTramontoOne(path string) (*One, error) {
	tramontoOne := &One{
		repoPath: path,
		mux:      &sync.Mutex{},
	}

	// Initializes the IPFS repo
	if err := tramontoIpfs.InitRepo(path); err != nil {
		return tramontoOne, err
	}

	// Initializes the Tramonto One database
	// db, err := tramontoDb.CreateOneSQLite(path)
	// if err != nil {
	// 	return tramontoOne, err
	// }

	// tramontoOne.database = db

	return tramontoOne, nil
}

// isRunning returns if the node is initialized and online
func (t *One) isRunning() (bool, error) {
	if t.node == nil || t.repo == nil {
		return false, nil
	}

	return t.node.OnlineMode(), nil
}

// Start starts the node
func (t *One) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !fsrepo.IsInitialized(t.repoPath) {
		return errors.New("Repo not initialized")
	}

	// Loads the plugins before create the node
	tramontoIpfs.LoadPlugins(t.repoPath)

	// Gets the repo
	nodeRepo, err := fsrepo.Open(t.repoPath)
	if err != nil {
		return err
	}

	t.repo = nodeRepo

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

	if t.node == nil || t.repo == nil {
		return errors.New("Node or repo not initilized")
	}

	if running, err := t.isRunning(); err != nil || !running {
		return errors.New("Node not running")
	}

	if err := t.node.Close(); err != nil {
		return err
	}

	if err := t.repo.Close(); err != nil {
		return err
	}

	t.node = nil
	t.repo = nil

	return nil
}
