package dht

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type Server struct {
	node      *Node
	listener  net.Listener
	listening bool
}

func NewServer(n *Node) *Server {
	return &Server{
		node: n,
	}
}

func (s *Server) Listen() {
	rpc.Register(s.node)
	//	n.create()
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", s.node.Port))
	if e != nil {
		panic(e)
	}
	s.node.create()
	s.listener = l
	s.listening = true
	go http.Serve(l, nil)
}

func (s *Server) Join(addr string) {
	s.Listen()
	s.node.join(addr)
}

func (s *Server) Quit() {
	s.listener.Close()
}

func (s *Server) Listening() bool {
	return s.listening
}

func (s *Server) Debug() string {
	return fmt.Sprintf(`
ID: %v
Listening: %v
Address: %v
Data: %v
Successor: %v
Predecessor: %v
Fingers: %v
`, s.node.Id, s.Listening(), s.node.addr(), s.node.Data, s.node.Successor, s.node.Predecessor, s.node.fingers[1:])

}
