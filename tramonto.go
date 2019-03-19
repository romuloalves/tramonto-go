package tramonto

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	files "gx/ipfs/QmQmhotPUzVrMEWNK3x1R5jQ5ZHWyL7tVUrmRPjrBrvyCb/go-ipfs-files"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	"sync"

	tramontoIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// One represents the Tramonto One lib
type One struct {
	repoPath string
	repo     repo.Repo
	node     *core.IpfsNode
	mux      *sync.Mutex
}

// NewTramontoOne returns a new instance of Tramonto One library
func NewTramontoOne(path string) (*One, error) {
	return &One{
		repoPath: path,
		mux:      &sync.Mutex{},
	}, nil
}

func (t *One) isRunning() (bool, error) {
	if t.node == nil {
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

	if !fsrepo.IsInitialized(t.repoPath) {
		return errors.New("Repo not initialized")
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

// NewTest publishes a new test to IPFS
func (t *One) NewTest(name, description string) (ipfsHash string, err error) {
	if t.node == nil {
		err = errors.New("node not started")
		return
	}

	coreAPI, err := coreapi.NewCoreAPI(t.node)
	if err != nil {
		return
	}

	// Creates the channel to create the IPNS key in background
	keyCh := make(chan iface.Key, 1)
	go func() {
		// Generates context to generate the IPNS key
		keyCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Gets key options
		ipnsKeyOpts := []options.KeyGenerateOption{options.Key.Type(options.RSAKey), options.Key.Size(options.DefaultRSALen)}

		// Generates the key
		key, err := coreAPI.Key().Generate(keyCtx, name, ipnsKeyOpts...)
		if err != nil {
			return
		}

		keyCh <- key
	}()

	// Generates content to the file
	jsonData := "{\"name\":\"" + name + "\", \"description\": \"" + description + "\"}"

	// Creates context to add
	addCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Adds content to IPFS
	ipfsPath, err := coreAPI.Unixfs().Add(addCtx, files.NewBytesFile([]byte(jsonData)))
	if err != nil {
		return
	}

	// Stores IPFS hash
	ipfsHash = ipfsPath.Cid().Hash().B58String()

	// Pins the file
	if err = coreAPI.Pin().Add(addCtx, ipfsPath); err != nil {
		return
	}

	// Takes the IPNS key
	ipnsKey := <-keyCh

	// Sets IPNS hash in the response

	// Gets publish options (publishing to IPNS)
	ipnsPublishOpts := []options.NamePublishOption{options.Name.Key(ipnsKey.Name())}

	// Publishes to IPNS
	_, err = coreAPI.Name().Publish(addCtx, ipfsPath, ipnsPublishOpts...)

	return
}
