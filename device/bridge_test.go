package device_test

import (
	"fmt"
	"net"

	. "github.com/teddyking/netsetgo/device"
	"github.com/vishvananda/netlink"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bridge", func() {
	var (
		bridge       *Bridge
		bridgeName   string
		bridgeIP     net.IP
		bridgeSubnet *net.IPNet
	)

	BeforeEach(func() {
		var err error
		bridge = NewBridge()
		bridgeName = "tower"
		bridgeIP, bridgeSubnet, err = net.ParseCIDR("10.10.10.1/24")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(cleanup(bridgeName)).To(Succeed())
	})

	Describe("Create", func() {
		It("creates a bridge with the provided name", func() {
			bridgeInterface, err := bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
			Expect(err).NotTo(HaveOccurred())

			Expect(bridgeInterface.Name).To(Equal(bridgeName))
		})

		It("assigns the provided address to the bridge", func() {
			bridgeInterface, err := bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
			Expect(err).NotTo(HaveOccurred())

			bridgeAddresses, err := bridgeInterface.Addrs()
			Expect(err).NotTo(HaveOccurred())

			Expect(len(bridgeAddresses)).To(Equal(1))
			Expect(bridgeAddresses[0].String()).To(Equal("10.10.10.1/24"))
		})

		Context("when a bridge with the provided name already exists", func() {
			BeforeEach(func() {
				_, err := bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
				Expect(err).NotTo(HaveOccurred())
			})

			It("doesn't error", func() {
				_, err := bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the bridge", func() {
				bridgeInterface, err := bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
				Expect(err).NotTo(HaveOccurred())

				Expect(bridgeInterface.Name).To(Equal(bridgeName))
			})
		})
	})

	Describe("Attach", func() {
		var (
			veth              *Veth
			vethNamePrefix    string
			bridgeInterface   *net.Interface
			hostVethInterface *net.Interface
		)

		BeforeEach(func() {
			var err error

			bridgeInterface, err = bridge.Create(bridgeName, bridgeIP, bridgeSubnet)
			Expect(err).NotTo(HaveOccurred())

			veth = NewVeth()
			vethNamePrefix = "veth"
			hostVethInterface, _, err = veth.Create(vethNamePrefix)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(cleanup(bridgeName)).To(Succeed())
			Expect(cleanup(fmt.Sprintf("%s0", vethNamePrefix))).To(Succeed())
		})

		It("attaches the provided veth to the provided bridge", func() {
			err := bridge.Attach(bridgeInterface, hostVethInterface)
			Expect(err).NotTo(HaveOccurred())

			Expect(fmt.Sprintf("/sys/class/net/%s0/master", vethNamePrefix)).To(BeAnExistingFile())
		})

		Context("when the bridge can't be found", func() {
			BeforeEach(func() {
				Expect(cleanup(bridgeName)).To(Succeed())
			})

			It("returns a descriptive error", func() {
				err := bridge.Attach(bridgeInterface, hostVethInterface)
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal("Link not found"))
			})
		})

		Context("when the veth can't be found", func() {
			BeforeEach(func() {
				Expect(cleanup(fmt.Sprintf("%s0", vethNamePrefix))).To(Succeed())
			})

			It("returns a descriptive error", func() {
				err := bridge.Attach(bridgeInterface, hostVethInterface)
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(Equal("Link not found"))
			})
		})
	})
})

func cleanup(name string) error {
	if _, err := net.InterfaceByName(name); err == nil {
		link, err := netlink.LinkByName(name)
		if err != nil {
			return err
		}
		return netlink.LinkDel(link)
	}
	return nil
}
