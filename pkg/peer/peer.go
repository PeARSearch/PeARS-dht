package peer

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
)

func NewPeer(ctx context.Context, seed int64, port int) (host.Host, error) {

	var randomness io.Reader

	if seed == 0 {
		randomness = rand.Reader
	} else {
		// giving a seed  ensures that the keys stay the same across multiple runs
		randomness = mrand.New(mrand.NewSource(seed))
	}

	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
	if err != nil {
		return nil, err
	}

	// addr, _ := multiaddr.NewMultiaddr()

	options := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Identity(privKey),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(),
		libp2p.EnableAutoRelay(),
		libp2p.NoSecurity, // may be we don't need this, this inits an insecure connection
	}

	return libp2p.New(options...)
}
