package tramonto

import (
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
)

// PeerID returns the identity of the peer
func (t *One) PeerID() (peer.ID, error) {
	return t.node.Identity, nil
}
