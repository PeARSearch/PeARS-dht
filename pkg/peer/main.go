package peer

import (
	"context"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"sync"

	"fmt"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	libPeer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
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
	Peers      []string
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

func (p *PeerConfig) getPeerInfo() []libPeer.AddrInfo {
	pinfos := make([]libPeer.AddrInfo, len(p.Peers))
	for i, addr := range p.Peers {
		maddr := multiaddr.StringCast(addr)
		p, err := libPeer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}
		pinfos[i] = *p
	}
	return pinfos
}

func (p *PeerConfig) Bootstrap(ctx context.Context) error {
	peers := p.getPeerInfo()
	log.Info(peers)
	if len(peers) < 1 {
		return fmt.Errorf("Not enough peers to bootstrap with")
	}

	log.Infof("We have peers %s", peers)

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup
	for _, pr := range peers {
		wg.Add(1)
		go func(pr libPeer.AddrInfo) {
			defer wg.Done()
			defer log.Println(ctx, "bootstrapDial", p.host, pr.ID)
			log.Printf("%s bootstrapping to %s", p.host, pr.ID)

			log.Info("test")
			log.Info(p.host.Peerstore())
			p.host.Peerstore().AddAddrs(pr.ID, pr.Addrs, peerstore.PermanentAddrTTL)
			if err := p.host.Connect(ctx, pr); err != nil {
				log.Println(ctx, "bootstrapDialFailed", pr.ID)
				log.Printf("failed to bootstrap with %v: %s", pr.ID, err)
				errs <- err
				return
			}
			log.Println(ctx, "bootstrapDialSuccess", pr.ID)
			log.Printf("bootstrapped with %v", pr.ID)
		}(pr)
	}

	wg.Wait()
	close(errs)
	count := 0
	var err error
	for err = range errs {
		if err != nil {
			count++
		}
	}
	if count == len(peers) {
		return fmt.Errorf("failed to bootstrap. %s", err)
	}

	return nil
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
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", p.ListenPort)),
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

	p.host = rhost.Wrap(basicHost, dht)

	if len(p.Peers) > 0 {
		err = p.Bootstrap(ctx)
		if err != nil {
			return err
		}
	}

	err = dht.Bootstrap(ctx)

	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", p.host.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := p.host.Addrs()

	log.Infof("This peer %d can be reached at:", p.host.ID().Pretty())
	for _, addr := range addrs {
		log.Info(addr.Encapsulate(hostAddr))
	}

	p.setHost(p.host)

	return nil
}
