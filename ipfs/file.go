package ipfs

import (
	"context"
	"errors"
	"io/ioutil"
	"time"

	cid "github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/interface-go-ipfs-core/options"
	ifacePath "github.com/ipfs/interface-go-ipfs-core/path"
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
func readContent(node *core.IpfsNode, path ifacePath.Path) ([]byte, error) {
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
