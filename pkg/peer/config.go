package peer

import (
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
