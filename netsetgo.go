package netsetgo

import "net"

type NetworkConfig struct {
	BridgeName     string
	BridgeIP       net.IP
	Subnet         *net.IPNet
	VethNamePrefix string
}

//go:generate counterfeiter . HostConfigurer
type HostConfigurer interface {
	Apply(netConfig NetworkConfig, pid int) error
}

type Netset struct {
	HostConfigurer HostConfigurer
}

func New(hostConfigurer HostConfigurer) *Netset {
	return &Netset{
		HostConfigurer: hostConfigurer,
	}
}

func (n *Netset) ConfigureHost(netConfig NetworkConfig, pid int) error {
	return n.HostConfigurer.Apply(netConfig, pid)
}
