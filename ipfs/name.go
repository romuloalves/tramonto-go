package ipfs

import (
	"context"
	"fmt"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	cid "gx/ipfs/QmTbxNB1NwDesLmKTscr4udL2tVP7MaxvXnD1D9yX7g3PN/go-cid"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	nsopts "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options/namesys"
	"time"
)

const ipnsTimeout = time.Second * 30

// PublishIPNS publishes a IPFS hash to IPNS
func PublishIPNS(node *core.IpfsNode, ipfsCid cid.Cid, key iface.Key) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	// Parses IPFS Cid to Path
	ipfsStrPath := fmt.Sprintf("/ipns/%s", ipfsCid.Hash().B58String())
	ipfsPath, err := iface.ParsePath(ipfsStrPath)
	if err != nil {
		return err
	}

	// Gets publish options (publishing to IPNS)
	ipnsPublishOpts := []options.NamePublishOption{
		options.Name.Key(key.Name()),
		options.Name.AllowOffline(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), ipnsTimeout)
	defer cancel()

	// Publishes to IPNS
	if _, err := api.Name().Publish(ctx, ipfsPath, ipnsPublishOpts...); err != nil {
		return err
	}

	return nil
}

// ResolveIPNS resolves the IPNS
func ResolveIPNS(node *core.IpfsNode, hash string) (iface.Path, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	// Generates context with timeout to resolve
	ctx, cancel := context.WithTimeout(context.Background(), ipnsTimeout)
	defer cancel()

	// Options to resolve until an IPFS is found
	nameResolveOpts := []options.NameResolveOption{
		options.Name.ResolveOption(nsopts.Depth(nsopts.UnlimitedDepth)),
		options.Name.ResolveOption(nsopts.DhtTimeout(ipnsTimeout)),
	}

	// Resolves the name
	key := fmt.Sprintf("/ipns/%s", hash)
	ipfsPath, err := api.Name().Resolve(ctx, key, nameResolveOpts...)
	if err != nil {
		return nil, err
	}

	return ipfsPath, nil
}
