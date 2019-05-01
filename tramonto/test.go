package tramonto

import (
	"encoding/json"
	"errors"

	oneCrypto "gitlab.com/tramonto-one/go-tramonto/crypto"
	"gitlab.com/tramonto-one/go-tramonto/entities"
)

// CreateTest creates a new test
// Uploads to IPFS and inserts in the database
func (t *TramontoOne) CreateTest(name, description string) ([]byte, error) {
	testResult := entities.NewEmptyTest()

	// Generates the test content
	metadata, err := entities.NewMetadata(name, description)
	if err != nil {
		return nil, err
	}

	testResult.Metadata = metadata

	// Creates a secret
	secret, err := oneCrypto.GenerateSecret()
	if err != nil {
		return nil, errors.New("Error generating secret: " + err.Error())
	}

	testResult.Secret = secret

	// Upload to IPFS
	ipfsHash, err := t.ipfs.UploadTest(metadata, secret)
	if err != nil {
		return nil, errors.New("(IPFS) Error uploading: " + err.Error())
	}

	testResult.Ipfs = ipfsHash

	// Adds test to database
	if err := t.db.InsertTest(testResult); err != nil {
		return nil, errors.New("(Database) Error inserting to the database: " + err.Error())
	}

	jsonResponse, err := json.Marshal(testResult)
	if err != nil {
		return nil, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonResponse, nil
}

// GetTests gets all the tests from the database
func (t *TramontoOne) GetTests() ([]byte, error) {
	// Finds tests
	tests, err := t.db.FindTests()
	if err != nil {
		return nil, errors.New("(Database) Error finding tests: " + err.Error())
	}

	jsonData, err := json.Marshal(tests)
	if err != nil {
		return nil, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}

// GetTestByIPFS returns a single test by its IPFS hash
func (t *TramontoOne) GetTestByIPFS(ipfsHash, secret string) ([]byte, error) {
	// Get Metadata from IPFS
	metadata, err := t.ipfs.GetTestByIPFS(ipfsHash, secret)
	if err != nil {
		return nil, errors.New("(IPFS) Cannot read from IPFS: " + err.Error())
	}

	ipnsKeyExists, ipnsKey, err := t.ipfs.GetKeyWithName(metadata.Name)
	if err != nil {
		return nil, errors.New("Error verifing IPNS key: " + err.Error())
	}

	test := entities.Test{
		Ipfs:           ipfsHash,
		IpnsKeyCreated: ipnsKeyExists,
		Secret:         secret,
		Metadata:       metadata,
	}

	if ipnsKeyExists {
		test.Ipns = ipnsKey
	}

	// Return the Test
	jsonData, err := json.Marshal(test)
	if err != nil {
		return nil, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}

// GetTestByIPNS returns a single test by its IPNS hash
func (t *TramontoOne) GetTestByIPNS(ipnsHash, secret string) ([]byte, error) {
	// Get Metadata from IPNS
	ipfsHash, metadata, err := t.ipfs.GetTestByIPNS(ipnsHash, secret)
	if err != nil {
		return nil, errors.New("(IPNS) Cannot read from IPNS: " + err.Error())
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
		return nil, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}

// ShareTest shares a test with IPNS
func (t *TramontoOne) ShareTest(ipfsHash, testName string) (string, error) {
	// Share with IPNS
	ipnsHash, err := t.ipfs.PublishToIPNS(ipfsHash, testName)
	if err != nil {
		return "", errors.New("(IPNS) Error sharing test: " + err.Error())
	}

	// Saves the IPNS hash in the database
	if err = t.db.SaveSharedTest(ipfsHash, ipnsHash); err != nil {
		return "", errors.New("(Database) Error saving hash: " + err.Error())
	}

	// Return the IPNS hash
	return ipnsHash, nil
}

// AddMember adds a new member to an existing test
func (t *TramontoOne) AddMember(ipns, name, email, role string) ([]byte, error) {
	// Finds test in the database
	test, err := t.db.FindTestByIpns(ipns)
	if err != nil {
		return nil, errors.New("(Database) Test not found: " + err.Error())
	}

	// Reads test config file from IPFS
	ipfsTest, err := t.ipfs.GetTestByIPFS(test.Ipfs, test.Secret)
	if err != nil {
		return nil, errors.New("(IPFS) Test not found: " + err.Error())
	}

	// Creates the member entity
	newMember, err := entities.NewMember(name, email, role)
	if err != nil {
		return nil, errors.New("Error creating member: " + err.Error())
	}

	// Adds the member to the metadata
	if err = ipfsTest.AddMember(newMember); err != nil {
		return nil, errors.New("Error adding member: " + err.Error())
	}

	// Uploads the test to IPFS
	newIpfsHash, err := t.ipfs.UploadTest(ipfsTest, test.Secret)
	if err != nil {
		return nil, errors.New("(IPFS) Error uploading test: " + err.Error())
	}

	// Publishes the new Metadata to IPNS
	// We should update the database just after a succeded publish to IPNS
	if _, err := t.ipfs.PublishToIPNS(newIpfsHash, test.Metadata.Name); err != nil {
		return nil, errors.New("(IPNS) Error publishing: " + err.Error())
	}

	// Updates the database
	if err = t.db.UpdateIPFSHash(ipns, newIpfsHash); err != nil {
		return nil, errors.New("(Database) Error updating data: " + err.Error())
	}

	// Return the Test
	jsonData, err := json.Marshal(ipfsTest)
	if err != nil {
		return nil, errors.New("Error parsing to json: " + err.Error())
	}

	return jsonData, nil
}
