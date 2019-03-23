package ipfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core"
	cid "gx/ipfs/QmTbxNB1NwDesLmKTscr4udL2tVP7MaxvXnD1D9yX7g3PN/go-cid"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"strings"

	"gitlab.com/tramonto-one/go-tramonto/entities"
)

// UploadTest uploads a test to IPFS
// Returns the IPFS hash
func (oneIpfs *OneIPFS) UploadTest(metadata entities.Metadata) (string, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	// Verifies if node is running
	if running := oneIpfs.isNodeRunning(); !running {
		return "", errors.New("Node is not running")
	}

	// Converts the metadata to json
	jsonRepresentation, err := metadata.ToJSON()
	if err != nil {
		return "", errors.New("Erro converting metadata to json: " + err.Error())
	}

	// Uploads json to IPFS
	ipfsCid, err := addContent(oneIpfs.node, jsonRepresentation, true)
	if err != nil {
		return "", errors.New("Error adding content: " + err.Error())
	}

	return ipfsCid.Hash().B58String(), nil
}

func getTestByIPFS(node *core.IpfsNode, path iface.Path, secret string) (entities.Metadata, error) {
	// Reads content
	content, err := readContent(node, path)
	if err != nil {
		return entities.Metadata{}, errors.New("Error reading content: " + err.Error())
	}

	// Parses from json
	var metadata entities.Metadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return entities.Metadata{}, errors.New("Error parsing to json: " + err.Error())
	}

	// Return
	return metadata, nil
}

// GetTestByIPFS returns a test Metadata by IPFS
func (oneIpfs *OneIPFS) GetTestByIPFS(hash, secret string) (entities.Metadata, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	// Parses IPFS hash to Path
	ipfsHashPath := fmt.Sprintf("/ipfs/%s", hash)
	ipfsPath, err := iface.ParsePath(ipfsHashPath)
	if err != nil {
		return entities.Metadata{}, errors.New("Error parsing path: " + err.Error())
	}

	// Reads and returns
	return getTestByIPFS(oneIpfs.node, ipfsPath, secret)
}

// GetTestByIPNS returns a test Metadata by IPFS
func (oneIpfs *OneIPFS) GetTestByIPNS(hash, secret string) (string, entities.Metadata, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	// Resolves IPNS to IPFS
	ipfsPath, err := resolveIPNS(oneIpfs.node, hash)
	if err != nil {
		return "", entities.Metadata{}, errors.New("Error resolving IPNS: " + err.Error())
	}

	// Reads and returns
	test, err := getTestByIPFS(oneIpfs.node, ipfsPath, secret)
	if err != nil {
		return "", test, err
	}

	ipfsHash := strings.Split(ipfsPath.String(), "/")[2]

	return ipfsHash, test, nil
}

// PublishTest publishes a test with IPNS
// Returns the IPNS hash
func (oneIpfs *OneIPFS) PublishTest(ipfsHash, keyName string) (string, error) {
	oneIpfs.mux.Lock()
	defer oneIpfs.mux.Unlock()

	ipnsHash := ""
	keyExists := false

	// Verifies key existence
	if keyName != "" {
		// Tries to retrieve key by its name
		key, err := keyWithName(oneIpfs.node, keyName)
		if err != nil {
			return ipnsHash, errors.New("Could not find IPNS keys")
		}

		// Key exists
		if key != nil {
			keyExists = true

			ipnsHash = key.ID().Pretty()
		}
	}

	// Have no key
	if !keyExists {
		// Generate IPNS key
		key, err := genKey(oneIpfs.node, keyName)
		if err != nil {
			return ipnsHash, errors.New("Error generating key: " + err.Error())
		}

		ipnsHash = key.ID().Pretty()
	}

	// Parses IPFS to CID
	ipfsCid, err := cid.Parse(ipfsHash)
	if err != nil {
		return ipfsHash, errors.New("Error parsing hash to cid: " + err.Error())
	}

	// Publish to IPNS
	if err := publishIPNS(oneIpfs.node, ipfsCid, keyName); err != nil {
		return ipnsHash, errors.New("Error publishing IPNS: " + err.Error())
	}

	// Return the IPNS hash
	return ipnsHash, nil
}
