package tramonto

import (
	"errors"

	"gitlab.com/tramonto-one/go-tramonto/entities"
)

// CreateTest creates a new test
// Uploads to IPFS and inserts in the database
func (t *TramontoOne) CreateTest(name, description string) (entities.Test, error) {
	testResult := entities.NewEmptyTest()

	// Generates the test content
	metadata, err := entities.NewMetadata(name, description)
	if err != nil {
		return testResult, err
	}

	testResult.Metadata = metadata

	// Upload to IPFS
	ipfsHash, err := t.ipfs.UploadTest(metadata)
	if err != nil {
		return testResult, errors.New("Erro uploading to IPFS: " + err.Error())
	}

	testResult.Ipfs = ipfsHash

	// TODO: Add to database

	return testResult, nil
}

// GetTestByIPFS returns a single test by its IPFS hash
func (t *TramontoOne) GetTestByIPFS(ipfsHash, secret string) (entities.Test, error) {
	// Get Metadata from IPFS
	metadata, err := t.ipfs.GetTestByIPFS(ipfsHash, secret)
	if err != nil {
		return entities.Test{}, errors.New("Cannot read from IPFS: " + err.Error())
	}

	// Return the Test
	return entities.Test{
		Ipfs:           ipfsHash,
		IpnsKeyCreated: false, // TODO: Verify about find IPNS key with the test name
		Secret:         secret,
		Metadata:       metadata,
	}, nil
}

// GetTestByIPNS returns a single test by its IPNS hash
func (t *TramontoOne) GetTestByIPNS(ipnsHash, secret string) (entities.Test, error) {
	// Get Metadata from IPNS
	metadata, err := t.ipfs.GetTestByIPNS(ipnsHash, secret)
	if err != nil {
		return entities.Test{}, errors.New("Cannot read from IPNS: " + err.Error())
	}

	// Return the Test
	return entities.Test{
		Ipns:           ipnsHash,
		IpnsKeyCreated: true, // TODO: Verify about find IPNS key with the test name
		Secret:         secret,
		Metadata:       metadata,
	}, nil
}

// ShareTest shares a test with IPNS
func (t *TramontoOne) ShareTest(ipfsHash, testName string) (string, error) {
	// Share with IPNS
	ipnsHash, err := t.ipfs.PublishTest(ipfsHash, testName)
	if err != nil {
		return ipnsHash, errors.New("Error sharing test: " + err.Error())
	}

	// Return the IPNS hash
	return ipnsHash, nil
}
