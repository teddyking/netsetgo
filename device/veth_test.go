package device_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo/device"
	. "github.com/teddyking/netsetgo/netsetgo_suite_helpers"

	"fmt"
	"net"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/vishvananda/netlink"
)

var _ = Describe("Veth", func() {
	var (
		veth           *Veth
		vethNamePrefix string
	)

	BeforeEach(func() {
		vethNamePrefix = "veth"
		veth = NewVeth()
	})

	AfterEach(func() {
		Expect(cleanup(fmt.Sprintf("%s0", vethNamePrefix))).To(Succeed())
	})

	Describe("Create", func() {
		It("creates a veth using the provided name prefix", func() {
			hostVeth, containerVeth, err := veth.Create(vethNamePrefix)
			Expect(err).NotTo(HaveOccurred())

			Expect(hostVeth.Name).To(Equal(fmt.Sprintf("%s0", vethNamePrefix)))
			Expect(containerVeth.Name).To(Equal(fmt.Sprintf("%s1", vethNamePrefix)))
		})

		It("brings the veth link up", func() {
			_, _, err := veth.Create(vethNamePrefix)
			Expect(err).NotTo(HaveOccurred())

			stdout := gbytes.NewBuffer()
			cmd := exec.Command("sh", "-c", "ip link show veth0")
			_, err = gexec.Start(cmd, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Consistently(stdout).ShouldNot(gbytes.Say("DOWN"))
		})

		Context("when a veth pair using the provided name prefix already exists", func() {
			BeforeEach(func() {
				_, _, err := veth.Create(vethNamePrefix)
				Expect(err).NotTo(HaveOccurred())
			})

			It("doesn't error", func() {
				_, _, err := veth.Create(vethNamePrefix)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the host and container veths", func() {
				hostVeth, containerVeth, err := veth.Create(vethNamePrefix)
				Expect(err).NotTo(HaveOccurred())

				Expect(hostVeth.Name).To(Equal("veth0"))
				Expect(containerVeth.Name).To(Equal("veth1"))
			})

			Context("and the link is already up", func() {
				BeforeEach(func() {
					link, err := netlink.LinkByName(fmt.Sprintf("%s0", vethNamePrefix))
					Expect(err).NotTo(HaveOccurred())
					Expect(netlink.LinkSetUp(link)).To(Succeed())
				})

				It("doesn't error", func() {
					_, _, err := veth.Create(vethNamePrefix)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns the host and container veths", func() {
					hostVeth, containerVeth, err := veth.Create(vethNamePrefix)
					Expect(err).NotTo(HaveOccurred())

					Expect(hostVeth.Name).To(Equal("veth0"))
					Expect(containerVeth.Name).To(Equal("veth1"))
				})
			})
		})
	})

	Describe("MoveToNetworkNamespace", func() {
		var (
			containerVeth *net.Interface
			parentPid     int
			pid           int
		)

		BeforeEach(func() {
			var err error
			_, containerVeth, err = veth.Create(vethNamePrefix)
			Expect(err).NotTo(HaveOccurred())

			CreateNetNamespace("testNetNamespace")
			parentPid, pid = RunCmdInNetNamespace("testNetNamespace", "sleep 1000")
		})

		AfterEach(func() {
			KillCmd(parentPid)
			DestroyNetNamespace("testNetNamespace")
		})

		It("moves the container's side of the veth into the namespace identified by the pid", func() {
			err := veth.MoveToNetworkNamespace(containerVeth, pid)
			Expect(err).NotTo(HaveOccurred())

			EnsureOutputForCommand("ip netns exec testNetNamespace ip addr", fmt.Sprintf("%s1", vethNamePrefix))
		})

		Context("when the veth doesn't exist", func() {
			It("returns a descriptive error", func() {
				nonexistentVeth := &net.Interface{Name: "nonexistentVeth"}
				err := veth.MoveToNetworkNamespace(nonexistentVeth, pid)

				Expect(err.Error()).To(ContainSubstring("Link not found"))
			})
		})

		Context("when the process doesn't exist", func() {
			It("returns a descriptive error", func() {
				err := veth.MoveToNetworkNamespace(containerVeth, -1)

				Expect(err.Error()).To(ContainSubstring("no such process"))
			})
		})
	})
})
