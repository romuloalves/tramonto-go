package tramonto

import (
	"encoding/json"
	"errors"

	"gitlab.com/tramonto-one/go-tramonto/entities"
)

// CreateTest creates a new test
// Uploads to IPFS and inserts in the database
func (t *TramontoOne) CreateTest(name, description string) ([]byte, error) {
	testResult := entities.NewEmptyTest()

	// Generates the test content
	metadata, err := entities.NewMetadata(name, description)
	if err != nil {
		return []byte{}, err
	}

	// testResult.Metadata = metadata

	// Upload to IPFS
	ipfsHash, err := t.ipfs.UploadTest(metadata)
	if err != nil {
		return []byte{}, errors.New("Erro uploading to IPFS: " + err.Error())
	}

	testResult.Ipfs = ipfsHash

	// Adds test to database
	if err := t.db.InsertTest(testResult); err != nil {
		return []byte{}, errors.New("Error inserting to the database: " + err.Error())
	}

	jsonResponse, err := json.Marshal(testResult)
	if err != nil {
		return []byte{}, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonResponse, nil
}

// GetTests gets all the tests from the database
func (t *TramontoOne) GetTests() ([]byte, error) {
	// Finds tests
	tests, err := t.db.FindTests()
	if err != nil {
		return []byte{}, errors.New("Error finding tests: " + err.Error())
	}

	jsonData, err := json.Marshal(tests)
	if err != nil {
		return []byte{}, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}

// GetTestByIPFS returns a single test by its IPFS hash
func (t *TramontoOne) GetTestByIPFS(ipfsHash, secret string) ([]byte, error) {
	// Get Metadata from IPFS
	metadata, err := t.ipfs.GetTestByIPFS(ipfsHash, secret)
	if err != nil {
		return []byte{}, errors.New("Cannot read from IPFS: " + err.Error())
	}

	ipnsKeyExists, ipnsKey, err := t.ipfs.GetKeyWithName(metadata.name)
	if err != nil {
		return []byte{}, errors.New("Error verifing IPNS key: " + err.Error())
	}

	test := entities.Test{
		Ipfs:           ipfsHash,
		IpnsKeyCreated: ipnsKeyExists,
		Secret:         secret,
		Metadata:       metadata,
	}

	if ipnsKeyExists {
		test.Ipns = ipnsKey.Ipns
	}

	// Return the Test
	jsonData, err := json.Marshal(test)
	if err != nil {
		return []byte{}, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}

// GetTestByIPNS returns a single test by its IPNS hash
func (t *TramontoOne) GetTestByIPNS(ipnsHash, secret string) ([]byte, error) {
	// Get Metadata from IPNS
	ipfsHash, metadata, err := t.ipfs.GetTestByIPNS(ipnsHash, secret)
	if err != nil {
		return []byte{}, errors.New("Cannot read from IPNS: " + err.Error())
	}

	test := entities.Test{
		Ipns:           ipnsHash,
		Ipfs:           ipfsHash,
		IpnsKeyCreated: true,
		Secret:         secret,
		Metadata:       metadata,
	}

	// Return the Test
	jsonData, err := json.Marshal(test)
	if err != nil {
		return []byte{}, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
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
