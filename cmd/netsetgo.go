package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"code.cloudfoundry.org/guardian/kawasaki/netns"
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
	flag.StringVar(&containerAddress, "containerAddress", "10.10.10.2/24", "Address to assign to the container (CIDR notation)")
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	if pid == 0 {
		fmt.Println("ERROR - netsetgo needs a pid")
		os.Exit(1)
	}

	bridgeCreator := device.NewBridge()
	vethCreator := device.NewVeth()
	netnsExecer := &netns.Execer{}

	hostConfigurer := configurer.NewHostConfigurer(bridgeCreator, vethCreator)
	containerConfigurer := configurer.NewContainerConfigurer(netnsExecer)
	netset := netsetgo.New(hostConfigurer, containerConfigurer)

	bridgeIP, bridgeSubnet, err := net.ParseCIDR(bridgeAddress)
	check(err)

	containerIP, _, err := net.ParseCIDR(containerAddress)
	check(err)

	netConfig := netsetgo.NetworkConfig{
		BridgeName:     bridgeName,
		BridgeIP:       bridgeIP,
		ContainerIP:    containerIP,
		Subnet:         bridgeSubnet,
		VethNamePrefix: vethNamePrefix,
	}

	check(netset.ConfigureHost(netConfig, pid))
	check(netset.ConfigureContainer(netConfig, pid))
}

func check(err error) {
	if err != nil {
		fmt.Printf("ERROR - %s\n", err.Error())
		os.Exit(1)
	}
}
