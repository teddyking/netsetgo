package netsetgo_test

import (
	"net"
	"os/exec"

	. "github.com/teddyking/netsetgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("netsetgo", func() {
	Describe("CreateBridge", func() {
		AfterEach(func() {
			cmd := exec.Command("sh", "-c", "ip link delete tower")
			Expect(cmd.Run()).To(Succeed())
		})

		Context("when a device with the provided name doesn't already exist", func() {
			It("creates a bridge device with the provided name", func() {
				err := CreateBridge("tower")
				Expect(err).NotTo(HaveOccurred())

				_, err = net.InterfaceByName("tower")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a device with the provided name already exists", func() {
			BeforeEach(func() {
				cmd := exec.Command("sh", "-c", "ip link add name tower type bridge")
				Expect(cmd.Run()).To(Succeed())
			})

			It("doesn't error", func() {
				err := CreateBridge("tower")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
