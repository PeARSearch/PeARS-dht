> ⚠️ This library is in its early development stage. Currently it is a very simple [kademlia DHT](https://en.wikipedia.org/wiki/Kademlia). Please come back when we are in beta! 🥹

# PeARS-dht

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/PeARSearch/PeARS-dht)

DHT implementation for PeARS network

## Requirements

> You can click the button on top and open this repo on [gitpod](https://gitpod.io) to get an environment with PeARSd-dht already setup

```sh
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

## Usage

We use the `make` targets defined in `Makefile` to create a multi-node setup locally for testing purposes.

- Each pears-dht node need a port to work, the default being 8080. This is used by the DHT to function
- To start the first node, run the following command:

  ```
  make run
  ```

  The above, will start the DHT in port 8080

