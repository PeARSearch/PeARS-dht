package dht

import (
	"errors"
	"fmt"
	"math/big"
	"net/rpc"
)


func dial(addr string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		fmt.Println("error dialing", err)
		return nil
	}
	return client
}

func RPCNotify(addr, notice string) error {
	if addr == "" {
		return errors.New("Notify: rpc address was empty")
	}
	client := dial(addr)
	if client == nil {
		return errors.New("Client was nil")
	}
	defer client.Close()
	var response bool
	return client.Call("Node.Notify", notice, &response)
}

func RPCGetPredecessor(addr string) (string, error) {
	if addr == "" {
		return "", errors.New("FindPredecessor: rpc address was empty")
	}
	client := dial(addr)
	if client == nil {
		return "", errors.New("Client was nil")
	}
	defer client.Close()
	var response string
	err := client.Call("Node.GetPredecessor", false, &response)
	if err != nil {
		return "", err
	}
	if response == "" {
		return "", errors.New("Empty predecessor")
	}
	return response, nil
}

func RPCFindSuccessor(addr string, id *big.Int) (string, error) {
	if addr == "" {
		return "", errors.New("FindSuccessor: rpc address was empty")
	}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer client.Close()
	var response string
	err = client.Call("Node.FindSuccessor", id, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func RPCHealthCheck(addr string) (bool, error) {
	if addr == "" {
		return false, errors.New("HealthCheck: rpc address was empty")
	}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return false, err
	}
	defer client.Close()
	var response int
	err = client.Call("Node.Ping", 101, &response)
	// handle this a bit more gracefully
	if err != nil {
		return false, err
	}
	return true, nil
}
