package configurer

import (
	"net"

	"github.com/teddyking/netsetgo"
)

//go:generate counterfeiter . BridgeCreator
type BridgeCreator interface {
	Create(name string, ip net.IP, subnet *net.IPNet) (*net.Interface, error)
	Attach(bridge, hostVeth *net.Interface) error
}

//go:generate counterfeiter . VethCreator
type VethCreator interface {
	Create(vethNamePrefix string) (*net.Interface, *net.Interface, error)
	MoveToNetworkNamespace(containerVeth *net.Interface, pid int) error
}

type Host struct {
	BridgeCreator BridgeCreator
	VethCreator   VethCreator
}

func NewHostConfigurer(bridgeCreator BridgeCreator, vethCreator VethCreator) *Host {
	return &Host{
		BridgeCreator: bridgeCreator,
		VethCreator:   vethCreator,
	}
}

func (h *Host) Apply(netConfig netsetgo.NetworkConfig, pid int) error {
	bridge, err := h.BridgeCreator.Create(netConfig.BridgeName, netConfig.BridgeIP, netConfig.Subnet)
	if err != nil {
		return err
	}

	hostVeth, containerVeth, err := h.VethCreator.Create(netConfig.VethNamePrefix)
	if err != nil {
		return err
	}

	err = h.BridgeCreator.Attach(bridge, hostVeth)
	if err != nil {
		return err
	}

	err = h.VethCreator.MoveToNetworkNamespace(containerVeth, pid)
	if err != nil {
		return err
	}

	return nil
}
