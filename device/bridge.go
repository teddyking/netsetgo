package device

import (
	"net"

	"github.com/vishvananda/netlink"
)

type Bridge struct{}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) Create(name string, ip net.IP, subnet *net.IPNet) (*net.Interface, error) {
	if interfaceExists(name) {
		return net.InterfaceByName(name)
	}

	linkAttrs := netlink.LinkAttrs{Name: name}
	link := &netlink.Bridge{linkAttrs}

	if err := netlink.LinkAdd(link); err != nil {
		return nil, err
	}

	address := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: subnet.Mask}}
	if err := netlink.AddrAdd(link, address); err != nil {
		return nil, err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return nil, err
	}

	return net.InterfaceByName(name)
}

func (b *Bridge) Attach(bridge, hostVeth *net.Interface) error {
	bridgeLink, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		return err
	}

	hostVethLink, err := netlink.LinkByName(hostVeth.Name)
	if err != nil {
		return err
	}

	return netlink.LinkSetMaster(hostVethLink, bridgeLink.(*netlink.Bridge))
}

func interfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)

	return err == nil
}
