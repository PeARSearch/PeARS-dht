package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PeARSearch/PeARS-dht/pkg/dht"
	"github.com/spf13/cobra"
)

var bootstrapAddr string
var port string
var id string

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the PeARS DHT in a standalone or p2p mode",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

		config := dht.DefaultConfig()
		config.Id = id
		config.Addr = fmt.Sprintf("0.0.0.0:%s", port)
		config.Timeout = 10 * time.Millisecond
		config.MaxIdle = 100 * time.Millisecond
	
		n, err := dht.NewNode(config, nil)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		n.Stop()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&bootstrapAddr, "join-address", "j", "", "Peer address to join the network with")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "Port for the DHT to listen in locally")
	rootCmd.PersistentFlags().StringVarP(&id, "ID", "s", "0", "Seed to create the peer ID from")

	Execute()
}
