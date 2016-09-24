package device

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

type Veth struct{}

func NewVeth() *Veth {
	return &Veth{}
}

func (v *Veth) Create(namePrefix string) (*net.Interface, *net.Interface, error) {
	hostVethName := fmt.Sprintf("%s0", namePrefix)
	containerVethName := fmt.Sprintf("%s1", namePrefix)

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

	if err := netlink.LinkSetUp(veth); err != nil {
		return nil, nil, err
	}

	return vethInterfacesByName(hostVethName, containerVethName)
}

func (v *Veth) MoveToNetworkNamespace(containerVeth *net.Interface, pid int) error {
	containerVethLink, err := netlink.LinkByName(containerVeth.Name)
	if err != nil {
		return err
	}

	return netlink.LinkSetNsPid(containerVethLink, pid)
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
