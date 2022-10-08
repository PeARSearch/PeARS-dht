# cacophony-dht

DHT implementation for PeARS network

## Requirements

```bash
[docker](https://docs.docker.com/get-docker/)
```

## Usage

We use the `make` targets defined in `Makefile` to create a multi-node setup locally for testing purposes.

- Each pears-dht node need 2 ports to work, the default being 8080 and 8888. The
  former is the port for the DHT to function and the latter is for the REST API
  that accepts data from orchard.
- To start the first node, run the following command:

  ```
  make run-pears
  ```

  > The first time execution may take longer

  The above, will start the DHT in port 8080 and server in port 8888. If you
  want to override these ports, you can run:

  ```
  make run-pears port=5345 serverport=8765
  ```

  Upon the successful completion of the above command, you will find output that
  look like the following:

  ```
    INFO[0000] Creating the basic host for the peer                                                                                                              [2/1952]
    INFO[0000]

    Host ID: QmeZXiuPsLsoNpGPmP3UpdDe4pX1zgv18RXZU8cmwBy7rj
    INFO[0000]

    Connection addresses are:
    INFO[0000]   /ip4/192.168.1.206/tcp/8080/p2p/QmeZXiuPsLsoNpGPmP3UpdDe4pX1zgv18RXZU8cmwBy7rj
    INFO[0000]   /ip4/127.0.0.1/tcp/8080/p2p/QmeZXiuPsLsoNpGPmP3UpdDe4pX1zgv18RXZU8cmwBy7rj
    INFO[0000]

    **Please use one of the above IDs for other nodes to join**

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] GET    /health                   --> github.com/PeARSearch/cacophony-dht/pkg/client.Setup.func1 (3 handlers)
    [GIN-debug] GET    /search/:word             --> github.com/PeARSearch/cacophony-dht/pkg/client.Setup.func2 (3 handlers)
    [GIN-debug] GET    /store/:word/:url         --> github.com/PeARSearch/cacophony-dht/pkg/client.Setup.func3 (3 handlers)
    [GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
    Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
  ```

  Note one of the `ipv4` addresses listed. For eg: ` /ip4/192.168.1.206/tcp/8080/p2p/QmeZXiuPsLsoNpGPmP3UpdDe4pX1zgv18RXZU8cmwBy7rj`

  - Create a second node that joins the network by running:

  ```
  make run-pears port=7878 serverport=9089 contact=/ip4/192.168.1.206/tcp/8080/p2p/QmeZXiuPsLsoNpGPmP3UpdDe4pX1zgv18RXZU8cmwBy7rj
  ```

  Make sure the `port` and `serverport` are not already in use. If you want to
  connect to more that one contacts, you can give the ips in a comman separated list.
