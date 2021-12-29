package peer

import (
	"bufio"
	"context"

	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ha.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func doEcho(s network.Stream) error {

	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Infof("read %s", str)

	_, err = s.Write([]byte(str))

	return err
}

func Listener(ctx context.Context, ha host.Host, port int) {
	// this bit runs a listening server in port `port`
	address := getHostAddress(ha)

	log.Infof("Peer %s starts listening on port %d", address, port)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Info("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Error(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
}

func Sender(ctx context.Context, ha host.Host, targetPeer string) error {
	fullAddr := getHostAddress(ha)

	log.Infof("Sender has address %s", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Info("sender received new stream")
		if err := doEcho(s); err != nil {
			log.Error(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	ipfsaddr, err := multiaddr.NewMultiaddr(targetPeer)
	if err != nil {
		return err
	}

	pid, err := ipfsaddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		return err
	}

	peerid, err := peer.Decode(pid)
	if err != nil {
		return err
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, peerstore.PermanentAddrTTL)

	log.Info("sender opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/echo/1.0.0")
	if err != nil {
		return err
	}

	log.Println("sender saying hello")
	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		return err
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}

	log.Info("read reply: %q\n", out)

	return nil
}
