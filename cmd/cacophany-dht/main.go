package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p-core/host"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PeARSearch/cacophony-dht/pkg/client"
	"github.com/PeARSearch/cacophony-dht/pkg/peer"
	// nolint:typecheck
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var logLevel string
var cfgFile string

var peerConfig = peer.NewPeerConfig()

var pearCmd = &cobra.Command{ // nolint:gochecknoglobals
	PersistentPreRunE: configLogger,
	Use:               "cacophony-dht",
	Short:             "Let's make some noise",
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(cmd.Context())

		log.Debug("Creating the basic host for the peer")

		host, err := peer.NewPeer(ctx, int64(peerConfig.Seed), peerConfig.ListenPort)
		if err != nil {
			return err
		}

		// initialize the Peer

		log.Infof("Host ID: %s", host.ID().Pretty())
		log.Info("Connection addresses are:")

		for _, addr := range host.Addrs() {
			log.Infof("  %s/p2p/%s", addr, host.ID().Pretty())
		}

		// initialize the DHT
		var bootstrapPeARS []multiaddr.Multiaddr
		if len(peerConfig.Contacts) > 0 {
			for _, id := range strings.Split(peerConfig.Contacts, ",") {
				multiaddress, err := multiaddr.NewMultiaddr(id)
				if err != nil {
					log.Errorf("Peer %s couldn't be connected to: %-v", id, err)
					continue
				}

				bootstrapPeARS = append(bootstrapPeARS, multiaddress)
			}
		}

		dht, err := peer.NewPearDHT(ctx, host, bootstrapPeARS)
		if err != nil {
			return err
		}

		result, err := dht.GetValue(ctx, "test")
		if err != nil {
			log.Info(err)

		}

		log.Info(result)
		r := client.Setup()
		addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8888)
		log.Debugf("ready to serve at %s", addr)

		go r.Run(addr)

		// TODO add the dht value and all to a config
		// TODO go run an API server to input data to

		run(host, cancel)

		return nil
	},
}

func run(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Printf("\rExiting...\n")

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}

func configLogger(cmd *cobra.Command, args []string) error {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithField("log-level", logLevel).Fatal("incorrect log level")

		return fmt.Errorf("incorrect log level")
	}

	log.SetLevel(lvl)
	log.WithField("log-level", logLevel).Debug("log level configured")

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))

			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", "CDHT", envVarSuffix))
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)

			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
		}
	})
}

func init() {
	v := readConfigFile()

	pearCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l",
		"info", "set log level verbosity (options: debug, info, error, warning)")
	bindFlags(pearCmd, v)

	pearCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file "+
		"(default is $HOME/.cacophony-dht.yaml)")

	pearCmd.Flags().IntVarP(&peerConfig.ListenPort, "port", "p", 0, "port to listen to")
	pearCmd.Flags().StringVarP(&peerConfig.Contacts, "contacts", "t", "", "target peers to dial(give a comma separated list)")
	pearCmd.Flags().IntVarP(&peerConfig.Seed, "seed", "s", 0, "random seed for id generation")

	pearCmd.MarkFlagRequired("listen") // we require the port to bind the service to
}

func readConfigFile() *viper.Viper {
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".pear-dht")
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Info(err)
		}
	}

	v.SetEnvPrefix("PDHT")

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	return v
}

func main() {
	if err := pearCmd.Execute(); err != nil {
		log.WithError(err).Fatal("error in the cli. Exiting")
		os.Exit(1)
	}
}
