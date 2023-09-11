package api

import (
	"context"

	"github.com/PeARSearch/pears-dht/pkg/dht"
)

type DHT struct {
	Node   *dht.Node
	Server *dht.Server
}

func NewPearsDHT(ctx context.Context, port string) *DHT {
	node := dht.NewNode(port)
	server := dht.NewServer(node)

	return &DHT{
		Node: node,
		Server: server,
	}
}

func (d DHT) Join(ctx context.Context, addr string) error {
	d.Server.Join(addr)

	return nil
}
