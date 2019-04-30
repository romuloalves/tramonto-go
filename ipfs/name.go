package ipfs

import (
	"context"
	"fmt"
	"time"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/interface-go-ipfs-core/options"
	nsopts "github.com/ipfs/interface-go-ipfs-core/options/namesys"
	ifacePath "github.com/ipfs/interface-go-ipfs-core/path"
)

const ipnsTimeout = time.Minute * 3

// publishIPNS publishes a IPFS hash to IPNS
func publishIPNS(node *core.IpfsNode, ipfsCid cid.Cid, keyName string) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	// Parses IPFS Cid to Path
	ipfsStrPath := fmt.Sprintf("/ipns/%s", ipfsCid.Hash().B58String())
	ipfsPath := ifacePath.New(ipfsStrPath)

	// Gets publish options (publishing to IPNS)
	ipnsPublishOpts := []options.NamePublishOption{
		options.Name.Key(keyName),
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

// resolveIPNS resolves the IPNS
func resolveIPNS(node *core.IpfsNode, hash string) (ifacePath.Path, error) {
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
