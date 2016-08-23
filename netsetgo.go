package netsetgo

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

type Netset struct {
	bridgeLinkAttrs netlink.LinkAttrs
	bridgeAddress   string
	vethNamePrefix  string
}

func New(bridgeName, bridgeAddress, vethNamePrefix string) *Netset {
	return &Netset{
		bridgeLinkAttrs: netlink.LinkAttrs{
			Name: bridgeName,
		},
		bridgeAddress:  bridgeAddress,
		vethNamePrefix: vethNamePrefix,
	}
}

func (n *Netset) CreateBridge() (*net.Interface, error) {
	if interfaceExists(n.bridgeLinkAttrs.Name) {
		return net.InterfaceByName(n.bridgeLinkAttrs.Name)
	}

	bridgeLink := &netlink.Bridge{n.bridgeLinkAttrs}

	if err := netlink.LinkAdd(bridgeLink); err != nil {
		return nil, err
	}

	bridgeAddress, err := netlink.ParseAddr(n.bridgeAddress)
	if err != nil {
		return nil, err
	}

	if err := netlink.AddrAdd(bridgeLink, bridgeAddress); err != nil {
		return nil, err
	}

	return net.InterfaceByName(n.bridgeLinkAttrs.Name)
}

func (n *Netset) CreateVethPair() (*net.Interface, *net.Interface, error) {
	hostVethName := fmt.Sprintf("%s0", n.vethNamePrefix)
	containerVethName := fmt.Sprintf("%s1", n.vethNamePrefix)

	if interfaceExists(hostVethName) {
		return vethInterfacesByName(hostVethName, containerVethName)
	}

	vethLinkAttrs := netlink.NewLinkAttrs()
	vethLinkAttrs.Name = hostVethName

	veth := &netlink.Veth{
		LinkAttrs: vethLinkAttrs,
		PeerName:  containerVethName,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		return nil, nil, err
	}

	return vethInterfacesByName(hostVethName, containerVethName)
}

func (n *Netset) AttachVethToBridge() error {
	hostVethName := fmt.Sprintf("%s0", n.vethNamePrefix)

	bridge, err := netlink.LinkByName(n.bridgeLinkAttrs.Name)
	if err != nil {
		return err
	}

	hostVeth, err := netlink.LinkByName(hostVethName)
	if err != nil {
		return err
	}

	return netlink.LinkSetMaster(hostVeth, bridge.(*netlink.Bridge))
}

func (n *Netset) PlaceVethInNetworkNs(pid int, vethNamePrefix string) error {
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

func vethInterfacesByName(hostVethName, containerVethName string) (*net.Interface, *net.Interface, error) {
	hostVeth, err := net.InterfaceByName(hostVethName)
	if err != nil {
		return nil, nil, err
	}

	containerVeth, err := net.InterfaceByName(containerVethName)
	if err != nil {
		return nil, nil, err
	}

	return hostVeth, containerVeth, nil
}
