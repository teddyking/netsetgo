package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"net"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo binary", func() {
	Context("when passed all required args", func() {
		var (
			parentPid int
			pid       int
		)

		BeforeEach(func() {
			createTestNetNamespace()
			parentPid, pid = runCmdInTestNetNamespace()

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
			killCmdInTestNetNamespace(parentPid)
			cleanupTestNetNamespace()
			cleanupTestBridge()
			cleanupTestVeth()
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
			stdout := gbytes.NewBuffer()
			cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip addr")
			_, err := gexec.Start(cmd, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(stdout).Should(gbytes.Say("veth1"))
		})

		It("assignss the provided IP address to the container's side of the veth pair", func() {
			stdout := gbytes.NewBuffer()
			cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip addr")
			_, err := gexec.Start(cmd, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(stdout).Should(gbytes.Say("10.10.10.2"))
		})

		It("sets the veth link to UP", func() {
			stdout := gbytes.NewBuffer()
			cmd := exec.Command("sh", "-c", "ip link show veth0")
			_, err := gexec.Start(cmd, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(stdout).Should(gbytes.Say("state UP"))
		})
	})
})
