package tramonto

// Callback represents an async callback function
type Callback interface {
	SendResult(json string)
}

// ShareTestAsync shares a test with IPNS in background
func (t *TramontoOne) ShareTestAsync(ipfsHash, keyName string, callback Callback) {
	go func() {
		// Share with IPNS
		ipnsHash, err := t.ipfs.PublishTest(ipfsHash, keyName)
		if err != nil {
			callback.SendResult("{\"error\":\"" + err.Error() + "\"}")
			return
		}

		// Return the IPNS hash
		callback.SendResult("{\"ipns\":\"" + ipnsHash + "\"}")
	}()
}
