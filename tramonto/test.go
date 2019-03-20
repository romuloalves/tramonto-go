package tramonto

import (
	"context"
	"errors"
	"fmt"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/core/coreapi"
	files "gx/ipfs/QmQmhotPUzVrMEWNK3x1R5jQ5ZHWyL7tVUrmRPjrBrvyCb/go-ipfs-files"
	iface "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core"
	"gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options"
	nsopts "gx/ipfs/QmXLwxifxwfc2bAwq6rdjbYqAsGzWsDE9RM5TWMGtykyj6/interface-go-ipfs-core/options/namesys"
	"io/ioutil"
	"strings"
	"time"

	oneFiles "gitlab.com/tramonto-one/go-tramonto/files"
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
	testResult := OneTest{}

	// Locks context
	t.mux.Lock()
	defer t.mux.Unlock()

	// Verifies if node is running
	if running, err := t.isRunning(); err != nil || !running {
		return testResult, errors.New("Node isn't running")
	}

	// Generates API
	coreAPI, err := coreapi.NewCoreAPI(t.node)
	if err != nil {
		return testResult, err
	}

	// Creates the channel to create the IPNS key in background
	keyCh := make(chan iface.Key, 1)
	go func() {
		// Generates context to generate the IPNS key
		keyCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Gets key options
		ipnsKeyOpts := []options.KeyGenerateOption{options.Key.Type(options.RSAKey), options.Key.Size(options.DefaultRSALen)}

		// Generates the key
		key, err := coreAPI.Key().Generate(keyCtx, name, ipnsKeyOpts...)
		if err != nil {
			return
		}

		keyCh <- key
	}()

	// Generates content to the file
	metadata, err := oneFiles.NewMetadata(name, description)
	if err != nil {
		return testResult, err
	}

	// Sets metadata to the test
	testResult.Metadata = metadata

	// Converts metadata to JSON
	jsonData, err := metadata.ToJSON()
	if err != nil {
		return testResult, err
	}

	// Creates context to add
	addCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Adds content to IPFS
	ipfsPath, err := coreAPI.Unixfs().Add(addCtx, files.NewBytesFile(jsonData))
	if err != nil {
		return testResult, err
	}

	testResult.IpfsHash = ipfsPath.Cid().Hash().B58String()

	// Pins the file
	if err := coreAPI.Pin().Add(addCtx, ipfsPath); err != nil {
		return testResult, err
	}

	// Takes the IPNS key
	ipnsKey := <-keyCh

	// Sets IPNS hash in the response
	testResult.IpnsHash = ipnsKey.ID().Pretty()

	// Gets publish options (publishing to IPNS)
	ipnsPublishOpts := []options.NamePublishOption{
		options.Name.Key(ipnsKey.Name()),
		options.Name.AllowOffline(true),
	}

	// Publishes to IPNS
	if _, err := coreAPI.Name().Publish(addCtx, ipfsPath, ipnsPublishOpts...); err != nil {
		return testResult, err
	}

	return testResult, nil
}

// GetTest returns the test with the given IPNS
func (t *One) GetTest(ipns, secret string) (OneTest, error) {
	testResult := OneTest{
		IpnsHash: ipns,
		Secret:   secret,
	}

	// Locks
	t.mux.Lock()
	defer t.mux.Unlock()

	if running, err := t.isRunning(); err != nil || !running {
		return testResult, err
	}

	// Generates API
	coreAPI, err := coreapi.NewCoreAPI(t.node)
	if err != nil {
		return testResult, err
	}

	ipnsTimeout := time.Second * 30

	// Generates context with timeout to resolve
	ctx, cancel := context.WithTimeout(context.Background(), ipnsTimeout)
	defer cancel()

	// Options to resolve until an IPFS is found
	nameResolveOpts := []options.NameResolveOption{
		options.Name.ResolveOption(nsopts.Depth(nsopts.UnlimitedDepth)),
		options.Name.ResolveOption(nsopts.DhtTimeout(ipnsTimeout)),
	}

	// Resolves the name
	key := fmt.Sprintf("/ipns/%s", ipns)
	ipfsPath, err := coreAPI.Name().Resolve(ctx, key, nameResolveOpts...)
	if err != nil {
		return testResult, err
	}

	// Gets the IPFS hash
	testResult.IpfsHash = strings.Split(ipfsPath.String(), "/")[2]

	// Gets the content of the file
	f, err := coreAPI.Unixfs().Get(ctx, ipfsPath)
	if err != nil {
		return testResult, err
	}
	defer f.Close()

	// Reads the file content
	var file files.File
	switch f := f.(type) {
	case files.File:
		file = f
	default:
		return testResult, errors.New("Test metadata is not a file")
	}

	fileContent, err := ioutil.ReadAll(file)
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
