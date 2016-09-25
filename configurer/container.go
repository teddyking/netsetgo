package configurer

import (
	"fmt"
	"net"
	"os"

	"code.cloudfoundry.org/guardian/kawasaki/netns"
	"github.com/teddyking/netsetgo"
	"github.com/vishvananda/netlink"
)

type Container struct {
	NetnsExecer *netns.Execer
}

func NewContainerConfigurer(netnsExecer *netns.Execer) *Container {
	return &Container{
		NetnsExecer: netnsExecer,
	}
}

func (c *Container) Apply(netConfig netsetgo.NetworkConfig, pid int) error {
	netnsFile, err := os.Open(fmt.Sprintf("/proc/%d/ns/net", pid))
	defer netnsFile.Close()
	if err != nil {
		return fmt.Errorf("Unable to find network namespace for process with pid '%d'", pid)
	}

	cbFunc := func() error {
		containerVethName := fmt.Sprintf("%s1", netConfig.VethNamePrefix)
		link, err := netlink.LinkByName(containerVethName)
		if err != nil {
			return fmt.Errorf("Container veth '%s' not found", containerVethName)
		}

		addr := &netlink.Addr{IPNet: &net.IPNet{IP: netConfig.ContainerIP, Mask: netConfig.Subnet.Mask}}
		err = netlink.AddrAdd(link, addr)
		if err != nil {
			return fmt.Errorf("Unable to assign IP address '%s' to %s", netConfig.ContainerIP, containerVethName)
		}

		if err := netlink.LinkSetUp(link); err != nil {
			return err
		}

		route := &netlink.Route{
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: link.Attrs().Index,
			Gw:        netConfig.BridgeIP,
		}

		return netlink.RouteAdd(route)
	}

	return c.NetnsExecer.Exec(netnsFile, cbFunc)
}
