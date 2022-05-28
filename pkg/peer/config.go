package peer

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type PeerConfig struct {
	ID         int
	ListenPort int
	Global     bool
	Seed       int
	host       host.Host
	Contacts   string
	dht        *dht.IpfsDHT
}

func NewPeerConfig() *PeerConfig {
	return &PeerConfig{}
}

func (p *PeerConfig) SetDht(dht *dht.IpfsDHT) {
	p.dht = dht
}

func (p *PeerConfig) PutData(ctx context.Context, word string, url string) error {
	return p.dht.PutValue(ctx, "/v/nandaka", []byte(url))
}
