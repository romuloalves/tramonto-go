package ipfs

import (
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

// IsNodeRunning returns if the node is initialized and online
func IsNodeRunning(node *core.IpfsNode) (bool, error) {
	if node == nil {
		return false, nil
	}

	return node.OnlineMode(), nil
}
