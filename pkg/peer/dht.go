package peer

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

func NewPearDHT(ctx context.Context, host host.Host, bootstrapPears []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
	var options []dht.Option
	var wg sync.WaitGroup

	if len(bootstrapPears) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kadht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = kadht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	for _, peerAddr := range bootstrapPears {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Warn(err)
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Errorf("Error while connect to peer %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}

	wg.Wait()

	return kadht, nil
}
