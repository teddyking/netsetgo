package main

import (
	"flag"

	"github.com/teddyking/netsetgo"
)

func main() {
	var bridgeName string

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.Parse()

	netsetgo.CreateBridge(bridgeName)
}
