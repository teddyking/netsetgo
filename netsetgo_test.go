package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo"

	"errors"

	"github.com/teddyking/netsetgo/netsetgofakes"
)

var _ = Describe("netsetgo", func() {
	var (
		fakeHostConfigurer      *netsetgofakes.FakeConfigurer
		fakeContainerConfigurer *netsetgofakes.FakeConfigurer
		netConfig               NetworkConfig
		pid                     int
		netset                  *Netset
	)

	BeforeEach(func() {
		fakeHostConfigurer = &netsetgofakes.FakeConfigurer{}
		fakeContainerConfigurer = &netsetgofakes.FakeConfigurer{}
		netConfig = NetworkConfig{BridgeName: "tower"}
		pid = 100

		netset = New(fakeHostConfigurer, fakeContainerConfigurer)
	})

	Describe("ConfigureHost()", func() {
		It("configures the host", func() {
			Expect(netset.ConfigureHost(netConfig, pid)).To(Succeed())
			Expect(fakeHostConfigurer.ApplyCallCount()).To(Equal(1))

			netConfigArg, pidArg := fakeHostConfigurer.ApplyArgsForCall(0)
			Expect(netConfigArg).To(Equal(netConfig))
			Expect(pidArg).To(Equal(pid))
		})

		Context("when the HostConfigurer returns an error", func() {
			BeforeEach(func() {
				fakeHostConfigurer.ApplyReturns(errors.New("error configuring host"))
			})

			It("returns the error", func() {
				Expect(netset.ConfigureHost(netConfig, pid)).NotTo(Succeed())
			})
		})
	})

	Describe("ConfigureContainer()", func() {
		It("configures the container", func() {
			Expect(netset.ConfigureContainer(netConfig, pid)).To(Succeed())
			Expect(fakeContainerConfigurer.ApplyCallCount()).To(Equal(1))

			netConfigArg, pidArg := fakeContainerConfigurer.ApplyArgsForCall(0)
			Expect(netConfigArg).To(Equal(netConfig))
			Expect(pidArg).To(Equal(pid))
		})

		Context("when the ContainerConfigurer returns an error", func() {
			BeforeEach(func() {
				fakeContainerConfigurer.ApplyReturns(errors.New("error configuring container"))
			})

			It("returns the error", func() {
				Expect(netset.ConfigureContainer(netConfig, pid)).NotTo(Succeed())
			})
		})
	})
})
