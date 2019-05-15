package ipfs

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	cid "github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/interface-go-ipfs-core/options"
	ifacePath "github.com/ipfs/interface-go-ipfs-core/path"
	oneCrypto "gitlab.com/tramonto-one/go-tramonto/crypto"
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

	// Pins the content
	if err = pin(node, path); err != nil {
		return nil, errors.New("Error pinning the object: " + err.Error())
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

// pin will pin the ipfs in the node
func pin(node *core.IpfsNode, path ifacePath.Path) error {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := []options.PinAddOption{
		options.Pin.Recursive(true),
	}

	// Pins the item
	if err = api.Pin().Add(ctx, path, options...); err != nil {
		return err
	}

	return nil
}

// ReadArtifact will read the artifact of the specific hash
func (oneIpfs *OneIPFS) ReadArtifact(ipfsHash, secret string) ([]byte, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	// Parses IPFS hash to Path
	ipfsHashPath := fmt.Sprintf("/ipfs/%s", ipfsHash)
	ipfsPath := ifacePath.New(ipfsHashPath)

	// Reads the content of the IPFS hash
	content, err := readContent(oneIpfs.node, ipfsPath)
	if err != nil {
		return nil, err
	}

	// Decrypts the content and just returns it
	return oneCrypto.DecryptArtifact(secret, content)
}

// UploadArtifact updates an artifact to IPFS
func (oneIpfs *OneIPFS) UploadArtifact(content []byte, secret string) (string, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	// Encrypts the content
	encryptedContent, err := oneCrypto.EncryptArtifact(secret, content)
	if err != nil {
		return "", errors.New("Could not encrypt artifact: " + err.Error())
	}

	// Uploads to IPFS
	cid, err := addContent(oneIpfs.node, encryptedContent, true)
	if err != nil {
		return "", errors.New("Could not upload file to IPFS: " + err.Error())
	}

	return cid.Hash().B58String(), nil
}
