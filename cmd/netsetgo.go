package main

import (
	"flag"

	"github.com/teddyking/netsetgo"
)

func main() {
	var bridgeName, bridgeAddress string

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.StringVar(&bridgeAddress, "bridgeAddress", "10.10.10.1/24", "Address to assign to bridge device (CIDR notation)")
	flag.Parse()

	netsetgo.CreateBridge(bridgeName)
	netsetgo.AddAddressToBridge(bridgeName, bridgeAddress)
	netsetgo.SetBridgeUp(bridgeName)
}
