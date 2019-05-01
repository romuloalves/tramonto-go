package ipfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	ifacePath "github.com/ipfs/interface-go-ipfs-core/path"

	oneCrypto "gitlab.com/tramonto-one/go-tramonto/crypto"
	"gitlab.com/tramonto-one/go-tramonto/entities"
)

// UploadTest uploads a test to IPFS
// Returns the IPFS hash
func (oneIpfs *OneIPFS) UploadTest(metadata entities.Metadata, secret string) (string, error) {
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

	encryptedData, err := oneCrypto.EncryptConfigFile(secret, jsonRepresentation)
	if err != nil {
		return "", errors.New("Error encrypting data: " + err.Error())
	}

	// Uploads json to IPFS
	ipfsCid, err := addContent(oneIpfs.node, encryptedData, true)
	if err != nil {
		return "", errors.New("Error adding content: " + err.Error())
	}

	return ipfsCid.Hash().B58String(), nil
}

func getTestByIPFS(node *core.IpfsNode, path ifacePath.Path, secret string) (entities.Metadata, error) {
	// Reads content
	content, err := readContent(node, path)
	if err != nil {
		return entities.Metadata{}, errors.New("Error reading content: " + err.Error())
	}

	// Pins IPFS
	if err = pin(node, path); err != nil {
		return entities.Metadata{}, errors.New("Error pinning: " + err.Error())
	}

	decryptedData, err := oneCrypto.DecryptConfigFile(secret, content)
	if err != nil {
		return entities.Metadata{}, errors.New("Error decrypting data: " + err.Error())
	}

	// Parses from json
	var metadata entities.Metadata
	if err := json.Unmarshal(decryptedData, &metadata); err != nil {
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
	ipfsPath := ifacePath.New(ipfsHashPath)

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

// PublishToIPNS publishes a test with IPNS
// Returns the IPNS hash
func (oneIpfs *OneIPFS) PublishToIPNS(ipfsHash, keyName string) (string, error) {
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
