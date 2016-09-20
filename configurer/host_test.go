package configurer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo/configurer"

	"errors"
	"net"

	"github.com/teddyking/netsetgo"
	"github.com/teddyking/netsetgo/configurer/configurerfakes"
)

var _ = Describe("Host", func() {
	Describe("Apply", func() {
		var (
			netConfig         netsetgo.NetworkConfig
			pid               int
			fakeBridgeCreator *configurerfakes.FakeBridgeCreator
			fakeVethCreator   *configurerfakes.FakeVethCreator
			hostConfigurer    *Host
		)

		BeforeEach(func() {
			bridgeAddress := "10.10.10.1/24"
			ip, net, err := net.ParseCIDR(bridgeAddress)
			Expect(err).NotTo(HaveOccurred())

			netConfig = netsetgo.NetworkConfig{
				BridgeName:     "tower",
				BridgeIP:       ip,
				Subnet:         net,
				VethNamePrefix: "veth",
			}
			pid = 100
		})

		BeforeEach(func() {
			fakeBridgeCreator = &configurerfakes.FakeBridgeCreator{}
			fakeVethCreator = &configurerfakes.FakeVethCreator{}
		})

		JustBeforeEach(func() {
			hostConfigurer = NewHostConfigurer(fakeBridgeCreator, fakeVethCreator)
		})

		It("creates a bridge", func() {
			Expect(hostConfigurer.Apply(netConfig, pid)).To(Succeed())
			Expect(fakeBridgeCreator.CreateCallCount()).To(Equal(1))
			name, ip, subnet := fakeBridgeCreator.CreateArgsForCall(0)
			Expect(name).To(Equal(netConfig.BridgeName))
			Expect(ip).To(Equal(netConfig.BridgeIP))
			Expect(subnet).To(Equal(netConfig.Subnet))
		})

		It("creates a veth", func() {
			Expect(hostConfigurer.Apply(netConfig, pid)).To(Succeed())
			Expect(fakeVethCreator.CreateCallCount()).To(Equal(1))
			Expect(fakeVethCreator.CreateArgsForCall(0)).To(Equal("veth"))
		})

		It("attaches the host's side of the veth to the bridge", func() {
			fakeBridge := &net.Interface{Name: "fakeBridge"}
			fakeHostVeth := &net.Interface{Name: "fakeHostVeth"}

			fakeBridgeCreator.CreateReturns(fakeBridge, nil)
			fakeVethCreator.CreateReturns(fakeHostVeth, nil, nil)

			Expect(hostConfigurer.Apply(netConfig, pid)).To(Succeed())
			Expect(fakeBridgeCreator.AttachCallCount()).To(Equal(1))

			bridgeArg, hostVethArg := fakeBridgeCreator.AttachArgsForCall(0)
			Expect(bridgeArg).To(Equal(fakeBridge))
			Expect(hostVethArg).To(Equal(fakeHostVeth))
		})

		It("moves the container's side of the veth into the namespace identified by the provided pid", func() {
			fakeContainerVeth := &net.Interface{Name: "fakeContainerVeth"}

			fakeVethCreator.CreateReturns(nil, fakeContainerVeth, nil)

			Expect(hostConfigurer.Apply(netConfig, pid)).To(Succeed())
			Expect(fakeVethCreator.MoveToNetworkNamespaceCallCount()).To(Equal(1))

			containerVethArg, pidArg := fakeVethCreator.MoveToNetworkNamespaceArgsForCall(0)
			Expect(containerVethArg).To(Equal(fakeContainerVeth))
			Expect(pidArg).To(Equal(pid))
		})

		Context("when creating a bridge returns an error", func() {
			BeforeEach(func() {
				fakeBridgeCreator.CreateReturns(nil, errors.New("error creating bridge"))
			})

			It("returns the error", func() {
				Expect(hostConfigurer.Apply(netConfig, pid)).To(MatchError("error creating bridge"))
			})
		})

		Context("when creating a veth returns an error", func() {
			BeforeEach(func() {
				fakeVethCreator.CreateReturns(nil, nil, errors.New("error creating veth"))
			})

			It("returns the error", func() {
				Expect(hostConfigurer.Apply(netConfig, pid)).To(MatchError("error creating veth"))
			})
		})

		Context("when attaching a veth to a bridge returns an error", func() {
			BeforeEach(func() {
				fakeBridgeCreator.AttachReturns(errors.New("error attaching veth to bridge"))
			})

			It("returns the error", func() {
				Expect(hostConfigurer.Apply(netConfig, pid)).To(MatchError("error attaching veth to bridge"))
			})
		})

		Context("when moving the container's side of the veth into the namespace identified by the provided pid errors", func() {
			BeforeEach(func() {
				fakeVethCreator.MoveToNetworkNamespaceReturns(errors.New("error moving veth into namespace"))
			})

			It("returns the error", func() {
				Expect(hostConfigurer.Apply(netConfig, pid)).To(MatchError("error moving veth into namespace"))
			})
		})
	})
})
