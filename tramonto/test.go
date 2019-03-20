package tramonto

import (
	"strings"

	oneFiles "gitlab.com/tramonto-one/go-tramonto/files"
	oneIpfs "gitlab.com/tramonto-one/go-tramonto/ipfs"
)

// OneTest represents a test in the Tramonto One
type OneTest struct {
	IpfsHash string            `json:"ipfsHash"`
	IpnsHash string            `json:"ipnsHash"`
	Secret   string            `json:"secret"`
	Metadata oneFiles.Metadata `json:"metadata,omitempty"`
}

// NewTest publishes a new test to IPFS
func (t *One) NewTest(name, description string) (OneTest, error) {
	// Locks context
	t.mux.Lock()
	defer t.mux.Unlock()

	// Verifies if node is running
	if running, err := oneIpfs.IsNodeRunning(t.node); err != nil || !running {
		return OneTest{}, err
	}

	testResult := OneTest{}

	// Generates content to the file
	metadata, err := oneFiles.NewMetadata(name, description)
	if err != nil {
		return testResult, err
	}

	testResult.Metadata = metadata

	// Converts metadata to JSON
	jsonData, err := metadata.ToJSON()
	if err != nil {
		return testResult, err
	}

	// Adds content to IPFS
	ipfsCid, err := oneIpfs.AddContent(t.node, jsonData, true)
	if err != nil {
		return testResult, err
	}

	// Stores IPFS hash
	testResult.IpfsHash = ipfsCid.Hash().B58String()

	// Generates IPNS key
	key, err := oneIpfs.GenKey(t.node, name)
	if err != nil {
		return testResult, err
	}

	// Stores IPNS hash
	testResult.IpnsHash = key.ID().Pretty()

	// Publishes to IPNS
	if err := oneIpfs.PublishIPNS(t.node, ipfsCid, key); err != nil {
		return testResult, err
	}

	return testResult, nil
}

// GetTest returns the test with the given IPNS
func (t *One) GetTest(ipns, secret string) (OneTest, error) {
	// Locks
	t.mux.Lock()
	defer t.mux.Unlock()

	if running, err := oneIpfs.IsNodeRunning(t.node); err != nil || !running {
		return OneTest{}, err
	}

	testResult := OneTest{
		IpnsHash: ipns,
		Secret:   secret,
	}

	// Gets the IPFS hash
	ipfsPath, err := oneIpfs.ResolveIPNS(t.node, ipns)
	if err != nil {
		return testResult, err
	}

	testResult.IpfsHash = strings.Split(ipfsPath.String(), "/")[2]

	// Reads the content of the hash
	fileContent, err := oneIpfs.ReadContent(t.node, ipfsPath)
	if err != nil {
		return testResult, err
	}

	// Transforms the json string to a Metadata struct
	metadata, err := oneFiles.MetadataFromJSON(fileContent)
	if err != nil {
		return testResult, err
	}

	testResult.Metadata = metadata

	return testResult, nil
}
