package dht

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	m = 161 // 161 so we can use 1-based indexing
)

type Node struct {
	Address     string
	Port        string
	Id          *big.Int
	Successor   string
	Predecessor string
	// TODO(nvn): Add a better datastore for the nodes
	Data        map[string]string
	fingers     [m]string
	next        int
}

type PutArgs struct {
	Key, Val string
}

// TODO(nvn): Make hash function configurable
func NewNode(port string) *Node {
	addr := GetAddress()

	return &Node{
		Address: addr,
		Port:    port,
		Id:      Hash(fmt.Sprintf("%v:%v", addr, port)),
		Data:    make(map[string]string),
	}
}

func (n *Node) addr() string {
	return fmt.Sprintf("%v:%v", n.Address, n.Port)
}

func (n *Node) Ping(one int, two *int) error {
	*two = 42
	return nil
}

func (n *Node) Put(args PutArgs, success *bool) error {
	n.Data[args.Key] = args.Val
	*success = true
	return nil
}

func (n *Node) Get(key string, response *string) error {
	*response = n.Data[key]
	return nil
}

func (n *Node) Delete(key string, response *bool) error {
	delete(n.Data, key)
	*response = true
	return nil
}

func (n *Node) FindSuccessor(id *big.Int, successor *string) error {
	// If the id I'm looking for falls between me and my successor
	// Then the data for this id will be found on my successor
	if InclusiveBetween(n.Id, id, Hash(n.Successor)) {
		*successor = n.Successor
		return nil
	}
	var err error
	*successor, err = RPCFindSuccessor(n.Successor, id)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Notify(addr string, response *bool) error {
	if n.Predecessor == "" || ExclusiveBetween(Hash(n.Predecessor), Hash(addr), n.Id) {
		n.Predecessor = addr
	}
	return nil
}

func (n *Node) GetPredecessor(none bool, addr *string) error {
	*addr = n.Predecessor
	return nil
}

func (n *Node) join(addr string) {
	n.Predecessor = ""
	// This is saying connect to who i'm trying to join and find my successor!
	addr, err := RPCFindSuccessor(addr, Hash(n.addr()))
	if err != nil {
		fmt.Printf("Error in join %v\n", err)
		return
	}
	n.Successor = addr

}

func (n *Node) stabalize() {
	// Successor's predecessor
	predecessor, err := RPCGetPredecessor(n.Successor)
	if err == nil {
		if ExclusiveBetween(n.Id, Hash(predecessor), Hash(n.Successor)) {
			n.Successor = predecessor
		}
	}
	err = RPCNotify(n.Successor, n.addr())
	if err != nil {
		fmt.Println(err)
	}
}

func (n *Node) checkPredecessor() {
	up, err := RPCHealthCheck(n.Predecessor)
	if err != nil || !up {
		n.Predecessor = ""
	}
}

func (n *Node) fixFingers() {
	n.next += 1
	if n.next > m-1 {
		n.next = 1
	}
	var resp string
	id := FingerEntry(n.addr(), n.next)
	err := n.FindSuccessor(id, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == "" {
		return
	}
	for InclusiveBetween(n.Id, id, Hash(resp)) {
		n.fingers[n.next] = resp
		n.next += 1
		if n.next > m-1 {
			n.next = 0
			break
		}
		id = FingerEntry(n.addr(), n.next)
	}
}

func (n *Node) create() {
	n.Predecessor = ""
	n.Successor = n.addr()
	go n.stabalizeOften()
	go n.checkPredecessorOften()
	go n.fixFingersOften()
}

func (n *Node) fixFingersOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.fixFingers()
		}
	}
}

func (n *Node) checkPredecessorOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.checkPredecessor()
		}
	}
}

func (n *Node) stabalizeOften() {
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			n.stabalize()
		}
	}
}

// inclusive id ε (left, right]
func InclusiveBetween(left, id, right *big.Int) bool {
	// if the right is bigger than the left then we know it doesn't cross zero
	if right.Cmp(left) == 1 {
		return left.Cmp(id) == -1 && id.Cmp(right) <= 0
	}
	return left.Cmp(id) == -1 || id.Cmp(right) <= 0
}

// inclusive id ε (left, right)
func ExclusiveBetween(left, id, right *big.Int) bool {
	if right.Cmp(left) == 1 {
		return left.Cmp(id) == -1 && id.Cmp(right) < 0
	}
	return left.Cmp(id) == -1 || id.Cmp(right) < 0
}

func Hash(in string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(in))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

const keySize = sha1.Size * 8

var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

func FingerEntry(start string, fingerentry int) *big.Int {
	id := Hash(start)
	two := big.NewInt(2)
	exponent := big.NewInt(int64(fingerentry) - 1)
	two.Exp(two, exponent, nil)
	id.Add(id, two)
	return id.Mod(id, hashMod)
}

func GetAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, interf := range interfaces {
		flags := interf.Flags

		// get only not loopback and up interfaces
		if flags&(net.FlagLoopback|flags&net.FlagUp) == net.FlagUp {
			addrs, err := interf.Addrs()
			if err != nil {
				panic(err)
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						return ip4.String()
					}
				}
			}
		}

	}
	return ""
}