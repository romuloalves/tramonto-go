package tramonto

// Callback represents an async callback function
type Callback interface {
	Invoke(json string)
}

// ShareTestAsync shares a test with IPNS in background
func (t *TramontoOne) ShareTestAsync(ipfsHash, keyName string, callback Callback) {
	go func() {
		// Share with IPNS
		ipnsHash, err := t.ipfs.PublishTest(ipfsHash, keyName)
		if err != nil {
			callback.Invoke("{\"error\":\"" + err.Error() + "\"}")
			return
		}

		// Return the IPNS hash
		callback.Invoke("{\"ipns\":\"" + ipnsHash + "\"}")
	}()
}
