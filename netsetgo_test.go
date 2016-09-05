package netsetgo_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo"

	"github.com/teddyking/netsetgo/netsetgofakes"
)

var _ = Describe("netsetgo", func() {
	Describe("ConfigureHost()", func() {
		var (
			fakeHostConfigurer *netsetgofakes.FakeHostConfigurer
			netConfig          NetworkConfig
			pid                int
			netset             *Netset
		)

		BeforeEach(func() {
			fakeHostConfigurer = &netsetgofakes.FakeHostConfigurer{}
			netConfig = NetworkConfig{BridgeName: "tower"}
			pid = 100

			netset = New(fakeHostConfigurer)
		})

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
})
