# cacophony-dht
DHT implementation for PeARS network

## Usage

Clone this repo locally and run the following commands to build and use
cacophony-dht.

```bash
$ go build -o cacdht ./cmd/cacophony-dht/main.go
$ ./cacdht -h
Let's make some noise

Usage:
  cacophony-dht [flags]

Flags:
      --config string      config file (default is $HOME/.cacophony-dht.yaml)
  -h, --help               help for cacophony-dht
      --log-level string   set log level verbosity (options: debug, info, error, warning) (default "info")
$ .cacdht
INFO[0000] I don't do much yet
```

