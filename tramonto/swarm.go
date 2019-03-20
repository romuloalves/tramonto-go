package tramonto

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	"time"
)

// GetConnectedPeersQuantity returns the quantity of connected peers in the swarm
func (t *One) GetConnectedPeersQuantity() (int, error) {
	// Verifies if node is running
	if running, err := t.isRunning(); err != nil || !running {
		return 0, errors.New("Node not running")
	}

	// Locks context
	t.mux.Lock()
	defer t.mux.Unlock()

	// Generates IPFS api
	coreAPI, err := coreapi.NewCoreAPI(t.node)
	if err != nil {
		return 0, err
	}

	// Generates context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gets connections
	connections, err := coreAPI.Swarm().Peers(ctx)
	if err != nil {
		return 0, err
	}

	// Returns just the length
	return len(connections), nil
}
