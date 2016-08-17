package netsetgo

import (
	"net"

	"github.com/vishvananda/netlink"
)

func CreateBridge(name string) error {
	if bridgeExists(name) {
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

func bridgeExists(name string) bool {
	_, err := net.InterfaceByName(name)

	return err == nil
}
