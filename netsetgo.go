package netsetgo

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func CreateBridge(name string) error {
	if interfaceExists(name) {
		return nil
	}

	bridgeLinkAttrs := netlink.NewLinkAttrs()
	bridgeLinkAttrs.Name = name

	bridge := &netlink.Bridge{bridgeLinkAttrs}

	return netlink.LinkAdd(bridge)
}

func AddAddressToBridge(name, address string) error {
	bridgeLinkAttrs := netlink.NewLinkAttrs()
	bridgeLinkAttrs.Name = name

	bridge := &netlink.Bridge{bridgeLinkAttrs}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	return netlink.AddrAdd(bridge, addr)
}

func SetBridgeUp(name string) error {
	bridgeLinkAttrs := netlink.NewLinkAttrs()
	bridgeLinkAttrs.Name = name

	bridge := &netlink.Bridge{bridgeLinkAttrs}
	return netlink.LinkSetUp(bridge)
}

func CreateVethPair(namePrefix string) error {
	hostVethName := fmt.Sprintf("%s0", namePrefix)
	containerVethName := fmt.Sprintf("%s1", namePrefix)

	if interfaceExists(hostVethName) {
		return nil
	}

	vethLinkAttrs := netlink.NewLinkAttrs()
	vethLinkAttrs.Name = hostVethName

	veth := &netlink.Veth{
		LinkAttrs: vethLinkAttrs,
		PeerName:  containerVethName,
	}

	return netlink.LinkAdd(veth)
}

func interfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)

	return err == nil
}
