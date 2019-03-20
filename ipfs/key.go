package ipfs

import (
	"context"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
)

// GenKey generates a new IPNS key (RSA-2048 default)
// Returns its hash
func GenKey(node *core.IpfsNode, name string) (iface.Key, error) {
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
