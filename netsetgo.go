package netsetgo

import "net"

type NetworkConfig struct {
	BridgeName     string
	BridgeIP       net.IP
	ContainerIP    net.IP
	Subnet         *net.IPNet
	VethNamePrefix string
}

//go:generate counterfeiter . Configurer
type Configurer interface {
	Apply(netConfig NetworkConfig, pid int) error
}

type Netset struct {
	HostConfigurer      Configurer
	ContainerConfigurer Configurer
}

func New(hostConfigurer, containerConfigurer Configurer) *Netset {
	return &Netset{
		HostConfigurer:      hostConfigurer,
		ContainerConfigurer: containerConfigurer,
	}
}

func (n *Netset) ConfigureHost(netConfig NetworkConfig, pid int) error {
	return n.HostConfigurer.Apply(netConfig, pid)
}

func (n *Netset) ConfigureContainer(netConfig NetworkConfig, pid int) error {
	return n.ContainerConfigurer.Apply(netConfig, pid)
}
