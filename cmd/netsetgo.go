package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/teddyking/netsetgo"
	"github.com/teddyking/netsetgo/configurer"
	"github.com/teddyking/netsetgo/device"
)

func main() {
	var bridgeName, bridgeAddress, containerAddress, vethNamePrefix string
	var pid int

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.StringVar(&bridgeAddress, "bridgeAddress", "10.10.10.1/24", "Address to assign to bridge device (CIDR notation)")
	flag.StringVar(&vethNamePrefix, "vethNamePrefix", "veth", "Name prefix for veth devices")
	flag.StringVar(&containerAddress, "containerAddress", "10.10.10.2", "Address to assign to the container (CIDR notation)")
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	bridgeCreator := device.NewBridge()
	vethCreator := device.NewVeth()

	hostConfigurer := configurer.NewHostConfigurer(bridgeCreator, vethCreator)
	netset := netsetgo.New(hostConfigurer)

	bridgeIP, bridgeSubnet, err := net.ParseCIDR(bridgeAddress)
	check(err)

	netConfig := netsetgo.NetworkConfig{
		BridgeName:     bridgeName,
		BridgeIP:       bridgeIP,
		Subnet:         bridgeSubnet,
		VethNamePrefix: vethNamePrefix,
	}

	netset.ConfigureHost(netConfig, pid)
}

func check(err error) {
	if err != nil {
		fmt.Printf("ERROR - %s\n", err.Error())
		os.Exit(1)
	}
}
