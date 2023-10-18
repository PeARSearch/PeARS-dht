package dht

import (
	"errors"
	"fmt"
	protov1 "github.com/PeARSearch/PeARS-dht/pkg/proto/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	emptyNode                = &protov1.Node{}
	emptyRequest             = &protov1.ER{}
	emptyGetResponse         = &protov1.GetResponse{}
	emptySetResponse         = &protov1.SetResponse{}
	emptyDeleteResponse      = &protov1.DeleteResponse{}
	emptyRequestKeysResponse = &protov1.RequestKeysResponse{}
)

func Dial(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, opts...)
}

/*
	Transport enables a node to talk to the other nodes in
	the ring
*/
type Transport interface {
	Start() error
	Stop() error

	//RPC
	GetSuccessor(*protov1.Node) (*protov1.Node, error)
	FindSuccessor(*protov1.Node, []byte) (*protov1.Node, error)
	GetPredecessor(*protov1.Node) (*protov1.Node, error)
	Notify(*protov1.Node, *protov1.Node) error
	CheckPredecessor(*protov1.Node) error
	SetPredecessor(*protov1.Node, *protov1.Node) error
	SetSuccessor(*protov1.Node, *protov1.Node) error

	//Storage
	GetKey(*protov1.Node, string) (*protov1.GetResponse, error)
	SetKey(*protov1.Node, string, string) error
	DeleteKey(*protov1.Node, string) error
	RequestKeys(*protov1.Node, []byte, []byte) ([]*protov1.KV, error)
	DeleteKeys(*protov1.Node, []string) error
}

type GrpcTransport struct {
	config *Config

	timeout time.Duration
	maxIdle time.Duration

	sock *net.TCPListener

	pool    map[string]*grpcConn
	poolMtx sync.RWMutex

	server *grpc.Server

	shutdown int32
}

// func NewGrpcTransport(config *Config) (protov1.ChordClient, error) {
func NewGrpcTransport(config *Config) (*GrpcTransport, error) {

	addr := config.Addr
	// Try to start the listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	pool := make(map[string]*grpcConn)

	// Setup the transport
	grp := &GrpcTransport{
		sock:    listener.(*net.TCPListener),
		timeout: config.Timeout,
		maxIdle: config.MaxIdle,
		pool:    pool,
		config:  config,
	}

	grp.server = grpc.NewServer(config.ServerOpts...)

	// Done
	return grp, nil
}

type grpcConn struct {
	addr       string
	client     protov1.ChordClient
	conn       *grpc.ClientConn
	lastActive time.Time
}

func (g *grpcConn) Close() {
	g.conn.Close()
}

func (g *GrpcTransport) registerNode(node *Node) {
	protov1.RegisterChordServer(g.server, node)
}

func (g *GrpcTransport) GetServer() *grpc.Server {
	return g.server
}

// Gets an outbound connection to a host
func (g *GrpcTransport) getConn(
	addr string,
) (protov1.ChordClient, error) {

	g.poolMtx.RLock()

	if atomic.LoadInt32(&g.shutdown) == 1 {
		g.poolMtx.Unlock()
		return nil, fmt.Errorf("TCP transport is shutdown")
	}

	cc, ok := g.pool[addr]
	g.poolMtx.RUnlock()
	if ok {
		return cc.client, nil
	}

	var conn *grpc.ClientConn
	var err error
	conn, err = Dial(addr, g.config.DialOpts...)
	if err != nil {
		return nil, err
	}

	client := protov1.NewChordClient(conn)
	cc = &grpcConn{addr, client, conn, time.Now()}
	g.poolMtx.Lock()
	if g.pool == nil {
		g.poolMtx.Unlock()
		return nil, errors.New("must instantiate node before using")
	}
	g.pool[addr] = cc
	g.poolMtx.Unlock()

	return client, nil
}

func (g *GrpcTransport) Start() error {
	// Start RPC server
	go g.listen()

	// Reap old connections
	go g.reapOld()

	return nil

}

// Returns an outbound TCP connection to the pool
func (g *GrpcTransport) returnConn(o *grpcConn) {
	// Update the last asctive time
	o.lastActive = time.Now()

	// Push back into the pool
	g.poolMtx.Lock()
	defer g.poolMtx.Unlock()
	if atomic.LoadInt32(&g.shutdown) == 1 {
		o.conn.Close()
		return
	}
	g.pool[o.addr] = o
}

// Shutdown the TCP transport
func (g *GrpcTransport) Stop() error {
	atomic.StoreInt32(&g.shutdown, 1)

	// Close all the connections
	g.poolMtx.Lock()

	g.server.Stop()
	for _, conn := range g.pool {
		conn.Close()
	}
	g.pool = nil

	g.poolMtx.Unlock()

	return nil
}

// Closes old outbound connections
func (g *GrpcTransport) reapOld() {
	ticker := time.NewTicker(60 * time.Second)

	for {
		if atomic.LoadInt32(&g.shutdown) == 1 {
			return
		}
		select {
		case <-ticker.C:
			g.reap()
		}

	}
}

func (g *GrpcTransport) reap() {
	g.poolMtx.Lock()
	defer g.poolMtx.Unlock()
	for host, conn := range g.pool {
		if time.Since(conn.lastActive) > g.maxIdle {
			conn.Close()
			delete(g.pool, host)
		}
	}
}

// Listens for inbound connections
func (g *GrpcTransport) listen() {
	g.server.Serve(g.sock)
}

// GetSuccessor the successor ID of a remote node.
func (g *GrpcTransport) GetSuccessor(node *protov1.Node) (*protov1.Node, error) {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return client.GetSuccessor(ctx, emptyRequest)
}

// FindSuccessor the successor ID of a remote node.
func (g *GrpcTransport) FindSuccessor(node *protov1.Node, id []byte) (*protov1.Node, error) {
	// fmt.Println("yo", node.Id, id)
	client, err := g.getConn(node.Addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return client.FindSuccessor(ctx, &protov1.ID{Id: id})
}

// GetPredecessor the successor ID of a remote node.
func (g *GrpcTransport) GetPredecessor(node *protov1.Node) (*protov1.Node, error) {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return client.GetPredecessor(ctx, emptyRequest)
}

func (g *GrpcTransport) SetPredecessor(node *protov1.Node, pred *protov1.Node) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.SetPredecessor(ctx, pred)
	return err
}

func (g *GrpcTransport) SetSuccessor(node *protov1.Node, succ *protov1.Node) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.SetSuccessor(ctx, succ)
	return err
}

func (g *GrpcTransport) Notify(node, pred *protov1.Node) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.Notify(ctx, pred)
	return err

}

func (g *GrpcTransport) CheckPredecessor(node *protov1.Node) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.CheckPredecessor(ctx, &protov1.ID{Id: node.Id})
	return err
}

func (g *GrpcTransport) GetKey(node *protov1.Node, key string) (*protov1.GetResponse, error) {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return client.XGet(ctx, &protov1.GetRequest{Key: key})
}

func (g *GrpcTransport) SetKey(node *protov1.Node, key, value string) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.XSet(ctx, &protov1.SetRequest{Key: key, Value: value})
	return err
}

func (g *GrpcTransport) DeleteKey(node *protov1.Node, key string) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.XDelete(ctx, &protov1.DeleteRequest{Key: key})
	return err
}

func (g *GrpcTransport) RequestKeys(node *protov1.Node, from, to []byte) ([]*protov1.KV, error) {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	val, err := client.XRequestKeys(
		ctx, &protov1.RequestKeysRequest{From: from, To: to},
	)
	if err != nil {
		return nil, err
	}
	return val.Values, nil
}

func (g *GrpcTransport) DeleteKeys(node *protov1.Node, keys []string) error {
	client, err := g.getConn(node.Addr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	_, err = client.XMultiDelete(
		ctx, &protov1.MultiDeleteRequest{Keys: keys},
	)
	return err
}
