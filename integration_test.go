package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo/netsetgo_suite_helpers"

	"fmt"
	"net"
	"os/exec"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo binary", func() {
	Context("when an invalid pid is provided", func() {
		It("exits 1", func() {
			command := exec.Command(pathToNetsetgo)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
		})
	})
	Context("when passed all required args", func() {
		var (
			parentPid int
			pid       int
		)

		BeforeEach(func() {
			CreateNetNamespace("testNetNamespace")
			parentPid, pid = RunCmdInNetNamespace("testNetNamespace", "sleep 1000")

			command := exec.Command(pathToNetsetgo,
				"--bridgeName=tower",
				"--bridgeAddress=10.10.10.1/24",
				"--vethNamePrefix=veth",
				"--containerAddress=10.10.10.2/24",
				fmt.Sprintf("--pid=%d", pid),
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			KillCmd(parentPid)
			DestroyNetNamespace("testNetNamespace")
			DestroyBridge("tower")
			DestroyVeth("veth")
		})

		It("creates a bridge device on the host with the provided name", func() {
			_, err := net.InterfaceByName("tower")
			Expect(err).NotTo(HaveOccurred())
		})

		It("assignes the provided IP address to the bridge", func() {
			bridge, err := net.InterfaceByName("tower")
			Expect(err).NotTo(HaveOccurred())

			bridgeAddresses, err := bridge.Addrs()
			Expect(err).NotTo(HaveOccurred())

			Expect(bridgeAddresses[0].String()).To(Equal("10.10.10.1/24"))
		})

		It("creates a veth pair on the host using the provided name prefix", func() {
			_, err := net.InterfaceByName("veth0")
			Expect(err).NotTo(HaveOccurred())
		})

		It("attaches the host's side of the veth pair to the bridge", func() {
			Expect("/sys/class/net/veth0/master").To(BeAnExistingFile())
		})

		It("puts the container's side of the veth pair into the net ns of the process specified by the provided pid", func() {
			EnsureOutputForCommand("ip netns exec testNetNamespace ip addr", "veth1")
		})

		It("assignss the provided IP address to the container's side of the veth pair", func() {
			EnsureOutputForCommand("ip netns exec testNetNamespace ip addr", "10.10.10.2")
		})

		It("sets the veth link to UP", func() {
			EnsureOutputForCommand("ip link show veth0", "state UP")
		})
	})
})
