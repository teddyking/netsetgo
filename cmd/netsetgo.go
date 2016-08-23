package main

import (
	"flag"

	"github.com/teddyking/netsetgo"
)

func main() {
	var bridgeName, bridgeAddress, vethNamePrefix string
	var pid int

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.StringVar(&bridgeAddress, "bridgeAddress", "10.10.10.1/24", "Address to assign to bridge device (CIDR notation)")
	flag.StringVar(&vethNamePrefix, "vethNamePrefix", "veth", "Name prefix for veth devices")
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	netset := netsetgo.New(bridgeName, bridgeAddress, vethNamePrefix)

	netset.CreateBridge()
	netset.CreateVethPair()
	netset.AttachVethToBridge()
	netset.PlaceVethInNetworkNs(pid, vethNamePrefix)
}
