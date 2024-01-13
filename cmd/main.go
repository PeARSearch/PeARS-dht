package cmd

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
	Long: `Run a full fledged distributed network`,
	Run: func(cmd *cobra.Command, args []string) {

		config := dht.DefaultConfig()
		config.Id = id
		config.Addr = fmt.Sprintf("0.0.0.0:%s", port)
		config.Timeout = 10 * time.Millisecond
		config.MaxIdle = 100 * time.Millisecond

		// this call starts a GRPC server that currently can be used to communicate with other nodes
		n, err := dht.NewNode(config, nil)
		if err != nil {
			log.Fatal(err)
		}

		// the following code doesn't do anything yet
		pht := dht.NewPht(cmd.Context(), n)

		log.Printf("To connect to the new node, use the address %s", fmt.Sprint(pht.Node.Addr))

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		n.Stop()

		// let us quit grazefully
		os.Exit(0)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&bootstrapAddr, "join-address", "j", "", "Peer address to join the network with")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port for the DHT to listen in locally")
	rootCmd.PersistentFlags().StringVarP(&id, "ID", "s", "0", "Seed to create the peer ID from")
	_ = rootCmd.MarkFlagRequired(port)
}
