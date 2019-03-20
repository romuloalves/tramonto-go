package ipfs

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"sync"
)

// OneIPFS represents the IPFS repo to Tramonto One
type OneIPFS struct {
	repoPath string
	node     *core.IpfsNode
	mux      *sync.Mutex
}

// InitializeOneIPFS initializes the One IPFS
func InitializeOneIPFS(path string) (*OneIPFS, error) {
	one := &OneIPFS{
		repoPath: path,
		mux:      new(sync.Mutex),
	}

	return one, nil
}

// isNodeRunning returns if the node is initialized and online
func (t *OneIPFS) isNodeRunning() bool {
	if t.node == nil {
		return false
	}

	return t.node.OnlineMode()
}

// Start starts the node
func (t *OneIPFS) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	// Verifies if the node is already running
	if running := t.isNodeRunning(); running {
		return errors.New("Node is already running")
	}

	// Verifies if the repo is initialized
	if initialized := isRepoInitialized(t.repoPath); !initialized {
		return errors.New("Repo not initialized")
	}

	// Loads the plugins before create the node
	loadPlugins(t.repoPath)

	// Opens the repo
	nodeRepo, err := openRepo(t.repoPath)
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
func (t *OneIPFS) Stop() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.node == nil {
		return errors.New("Node not initilized")
	}

	if running := t.isNodeRunning(); !running {
		return errors.New("Node not running")
	}

	if err := t.node.Close(); err != nil {
		return errors.New("Failed to close the node: " + err.Error())
	}

	t.node = nil

	return nil
}
