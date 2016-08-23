package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo"

	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo", func() {
	var (
		netset         *Netset
		bridgeName     string
		bridgeAddress  string
		vethNamePrefix string
	)

	BeforeEach(func() {
		bridgeName = "tower"
		bridgeAddress = "10.10.10.1/24"
		vethNamePrefix = "veth"
	})

	JustBeforeEach(func() {
		netset = New(bridgeName, bridgeAddress, vethNamePrefix)
	})

	Describe("CreateBridge", func() {
		AfterEach(func() {
			cleanupTestBridge()
		})

		It("creates a bridge with the provided name", func() {
			bridgeInterface, err := netset.CreateBridge()
			Expect(err).NotTo(HaveOccurred())

			Expect(bridgeInterface.Name).To(Equal("tower"))
		})

		It("assigns the provided address to the bridge", func() {
			bridgeInterface, err := netset.CreateBridge()
			Expect(err).NotTo(HaveOccurred())

			bridgeAddrs, err := bridgeInterface.Addrs()
			Expect(err).NotTo(HaveOccurred())

			Expect(len(bridgeAddrs)).To(Equal(1))
		})

		Context("when a bridge with the provided name already exists", func() {
			BeforeEach(func() {
				createTestBridge()
			})

			It("doesn't error", func() {
				_, err := netset.CreateBridge()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the bridge", func() {
				bridgeInterface, err := netset.CreateBridge()
				Expect(err).NotTo(HaveOccurred())

				Expect(bridgeInterface.Name).To(Equal("tower"))
			})
		})

		Context("when an invalid CIDR address is provided", func() {
			BeforeEach(func() {
				bridgeAddress = "invalid CIDR address"
			})

			It("returns a descriptive error", func() {
				_, err := netset.CreateBridge()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("invalid CIDR address"))
			})
		})
	})

	Describe("CreateVethPair", func() {
		AfterEach(func() {
			cleanupTestVeth()
		})

		It("creates a veth pair using the provided name prefix", func() {
			hostVeth, containerVeth, err := netset.CreateVethPair()
			Expect(err).NotTo(HaveOccurred())

			Expect(hostVeth.Name).To(Equal("veth0"))
			Expect(containerVeth.Name).To(Equal("veth1"))
		})

		Context("when a veth pair using the provided name prefix already exists", func() {
			BeforeEach(func() {
				createTestVeth()
			})

			It("doesn't error", func() {
				_, _, err := netset.CreateVethPair()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the host and container veths", func() {
				hostVeth, containerVeth, err := netset.CreateVethPair()
				Expect(err).NotTo(HaveOccurred())

				Expect(hostVeth.Name).To(Equal("veth0"))
				Expect(containerVeth.Name).To(Equal("veth1"))
			})
		})
	})

	Describe("AttachVethToBridge", func() {
		BeforeEach(func() {
			createTestBridge()
			createTestVeth()
		})

		AfterEach(func() {
			cleanupTestBridge()
			cleanupTestVeth()
		})

		It("attaches the host's side of the veth pair to the bridge", func() {
			err := netset.AttachVethToBridge()
			Expect(err).NotTo(HaveOccurred())

			Expect("/sys/class/net/veth0/master").To(BeAnExistingFile())
		})

		Context("when the bridge doesn't exist", func() {
			BeforeEach(func() {
				bridgeName = "bridgenothere"
			})

			It("returns a descriptive error", func() {
				err := netset.AttachVethToBridge()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Link not found"))
			})
		})

		Context("when the veth pair doesn't exist", func() {
			BeforeEach(func() {
				vethNamePrefix = "vethnothere"
			})

			It("returns a descriptive error", func() {
				err := netset.AttachVethToBridge()
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Link not found"))
			})
		})
	})

	Describe("PlaceVethInNetworkNamespace", func() {
		var (
			parentPid, pid int
		)

		BeforeEach(func() {
			createTestNetNamespace()
			parentPid, pid = runCmdInTestNetNamespace()
			createTestVeth()
		})

		AfterEach(func() {
			killCmdInTestNetNamespace(parentPid)
			cleanupTestVeth()
			cleanupTestNetNamespace()
		})

		It("places the container's side of the veth pair into the namespace using the provided pid", func() {
			err := netset.PlaceVethInNetworkNs(pid, "veth")
			Expect(err).NotTo(HaveOccurred())

			stdout := gbytes.NewBuffer()
			cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip addr")
			_, err = gexec.Start(cmd, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(stdout).Should(gbytes.Say("veth1"))
		})

		Context("when the process doesn't exist", func() {
			It("returns a descriptive error", func() {
				err := netset.PlaceVethInNetworkNs(-1, "veth")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("no such process"))
			})
		})

		Context("when the veth pair doesn't exist", func() {
			It("returns a descriptive error", func() {
				err := netset.PlaceVethInNetworkNs(pid, "vethnothere")
				Expect(err).To(HaveOccurred())

				Expect(err.Error()).To(ContainSubstring("Link not found"))
			})
		})
	})
})
