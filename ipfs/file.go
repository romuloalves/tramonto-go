package ipfs

import (
	"context"
	"errors"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	files "gx/ipfs/QmQmhotPUzVrMEWNK3x1R5jQ5ZHWyL7tVUrmRPjrBrvyCb/go-ipfs-files"
	cid "gx/ipfs/QmTbxNB1NwDesLmKTscr4udL2tVP7MaxvXnD1D9yX7g3PN/go-cid"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	"io/ioutil"
	"time"
)

const catTimeout = time.Minute

// addContent adds a buffer to IPFS
func addContent(node *core.IpfsNode, content []byte, pin bool) (cid.Cid, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return cid.Cid{}, err
	}

	addCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Adds content to IPFS
	ipfsPath, err := api.Unixfs().Add(addCtx, files.NewBytesFile(content))
	if err != nil {
		return ipfsPath.Cid(), err
	}

	if !pin {
		return ipfsPath.Cid(), nil
	}

	// Pins the file
	if err := api.Pin().Add(addCtx, ipfsPath, options.Pin.Recursive(false)); err != nil {
		return ipfsPath.Cid(), err
	}

	return ipfsPath.Cid(), nil
}

// readContent reads the content in a hash
func readContent(node *core.IpfsNode, path iface.Path) ([]byte, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return []byte{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), catTimeout)
	defer cancel()

	// Gets the content of the file
	f, err := api.Unixfs().Get(ctx, path)
	if err != nil {
		return []byte{}, err
	}

	// Reads the file content
	var file files.File
	switch f := f.(type) {
	case files.File:
		file = f
	default:
		return []byte{}, errors.New("The path is not a file")
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return fileContent, nil
}
