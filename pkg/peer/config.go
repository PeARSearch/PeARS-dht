package peer

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type PeerConfig struct {
	ID         int
	ListenPort int
	ServerPort int
	Global     bool
	Seed       int
	host       host.Host
	Contacts   string
	dht        *dht.IpfsDHT
	// dht2       *dualdht.DHT
}

func NewPeerConfig() *PeerConfig {
	return &PeerConfig{}
}

func (p *PeerConfig) SetDht(dht *dht.IpfsDHT) {
	p.dht = dht
}

func (p *PeerConfig) PutData(ctx context.Context, word string, url string) error {
	return p.dht.PutValue(ctx, fmt.Sprintf("/v/%s", word), []byte(url))
}

func (p *PeerConfig) GetData(ctx context.Context, word string) ([]byte, error) {
	data, err := p.dht.GetClosestPeers(ctx, fmt.Sprintf("/v/%s", word))
	fmt.Println(data)
	fmt.Println(err)

	return p.dht.GetValue(ctx, fmt.Sprintf("/v/%s", word))
}
