package tramonto

import "gitlab.com/tramonto-one/go-tramonto/entities"

// ShareTestAsync shares a test with IPNS in background
func (t *TramontoOne) ShareTestAsync(ipfsHash, keyName string, callback entities.Callback) {
	go func() {
		// Share with IPNS
		ipnsHash, err := t.ipfs.PublishToIPNS(ipfsHash, keyName)
		if err != nil {
			callback.Invoke("{\"error\":\"" + err.Error() + "\"}")
			return
		}

		// Return the IPNS hash
		callback.Invoke("{\"ipns\":\"" + ipnsHash + "\"}")
	}()
}
