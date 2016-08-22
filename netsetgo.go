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

func AttachVethToBridge(bridgeName, vethNamePrefix string) error {
	bridge, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	hostVethName := fmt.Sprintf("%s0", vethNamePrefix)
	hostVeth, err := netlink.LinkByName(hostVethName)
	if err != nil {
		return err
	}

	return netlink.LinkSetMaster(hostVeth, bridge.(*netlink.Bridge))
}

func PlaceVethInNetworkNamespace(pid int, vethNamePrefix string) error {
	containerVethName := fmt.Sprintf("%s1", vethNamePrefix)

	containerVeth, err := netlink.LinkByName(containerVethName)
	if err != nil {
		return err
	}
	return netlink.LinkSetNsPid(containerVeth, pid)
}

func interfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)

	return err == nil
}
