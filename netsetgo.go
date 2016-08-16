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

func bridgeExists(name string) bool {
	_, err := net.InterfaceByName(name)

	return err == nil
}
