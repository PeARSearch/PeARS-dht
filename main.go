package main

import (
	"log"

	"github.com/PeARSearch/PeARS-dht/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
