package peer

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"

	"fmt"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type PeerConfig struct {
	ID         int
	ListenPort int
	Global     bool
	Seed       int
	Target     string
	host       host.Host
}

func NewPeerConfig() PeerConfig {
	return PeerConfig{}
}

func (p *PeerConfig) setHost(h host.Host) {
	p.host = h
}

func (p *PeerConfig) GetHost() host.Host {
	return p.host
}

func (p *PeerConfig) Bootstrap() {
	// TODO this will Bootstrap the local peer
}

func (p *PeerConfig) MakeBasicHost() error {
	var r io.Reader

	if p.Seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(int64(p.Seed)))
	}

	// generate a key pair for the host
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return err
	}

	options := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", p.ListenPort)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
		libp2p.NoSecurity, // may be we don't need this, this inits an insecure connection
	}

	ha, err := libp2p.New(options...)
	p.setHost(ha)

	return err
}

func (p *PeerConfig) MakeRoutedHost(ctx context.Context) error {
	var r io.Reader

	if p.Seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(int64(p.Seed)))
	}

	// generate a key pair for the host
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return err
	}

	options := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", p.ListenPort)),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	basicHost, err := libp2p.New(options...)
	if err != nil {
		return err
	}

	dstore := dsync.MutexWrap(ds.NewMapDatastore()) // simple in-memory data store for the DHT

	dht := dht.NewDHT(ctx, basicHost, dstore)

	routedHost := rhost.Wrap(basicHost, dht)

	// TODO bootstrap code
	//
	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()

	log.Infof("This peer %d can be reached at:", routedHost.ID().Pretty())
	for _, addr := range addrs {
		log.Info(addr.Encapsulate(hostAddr))
	}

	p.setHost(routedHost)

	return nil
}
