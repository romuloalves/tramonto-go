package ipfs

import (
	"context"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
)

// genKey generates a new IPNS key (RSA-2048 default)
// Returns its hash
func genKey(node *core.IpfsNode, name string) (iface.Key, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	// Generates context to generate the IPNS key
	keyCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gets key options
	ipnsKeyOpts := []options.KeyGenerateOption{
		options.Key.Type(options.RSAKey),
		options.Key.Size(options.DefaultRSALen),
	}

	// Generates the key
	key, err := api.Key().Generate(keyCtx, name, ipnsKeyOpts...)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// keyWithName return the key with the given name
func keyWithName(node *core.IpfsNode, name string) (iface.Key, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	// Generates context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// List all keys
	keys, err := api.Key().List(ctx)
	if err != nil {
		return nil, err
	}

	// Loops the keys to verify existence
	for _, key := range keys {
		if key.Name() != name {
			continue
		}

		// Key exists
		return key, nil
	}

	// Do not exists
	return nil, nil
}

// GetKeyWithName returns the key with the given name
func (t *OneIPFS) GetKeyWithName(name string) (bool, string, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	key, err := keyWithName(t.node, name)
	if err != nil {
		return false, "", err
	}

	if key == nil {
		return false, "", nil
	}

	return true, key.ID().Pretty(), nil
}
