package configurer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo/configurer"
	. "github.com/teddyking/netsetgo/netsetgo_suite_helpers"

	"net"
	"os/exec"

	"code.cloudfoundry.org/guardian/kawasaki/netns"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/teddyking/netsetgo"
)

var _ = Describe("ContainerConfigurer", func() {
	var (
		parentPid, pid      int
		netnsExecer         *netns.Execer
		containerConfigurer *Container
		netConfig           netsetgo.NetworkConfig
	)

	BeforeEach(func() {
		CreateNetNamespace("testNetNamespace")
		parentPid, pid = RunCmdInNetNamespace("testNetNamespace", "sleep 1000")
		CreateVethInNetNamespace("veth", "testNetNamespace")

		netnsExecer = &netns.Execer{}
		containerConfigurer = NewContainerConfigurer(netnsExecer)

		bridgeAddress := "10.10.10.1/24"
		bridgeIP, _, err := net.ParseCIDR(bridgeAddress)
		Expect(err).NotTo(HaveOccurred())

		containerAddress := "10.10.10.10/24"
		containerIP, net, err := net.ParseCIDR(containerAddress)
		Expect(err).NotTo(HaveOccurred())

		netConfig = netsetgo.NetworkConfig{
			BridgeIP:       bridgeIP,
			ContainerIP:    containerIP,
			Subnet:         net,
			VethNamePrefix: "veth",
		}
	})

	AfterEach(func() {
		KillCmd(parentPid)
		DestroyNetNamespace("testNetNamespace")
		DestroyVeth("veth")
	})

	It("assigns the provided address to the container's side of the veth", func() {
		Expect(containerConfigurer.Apply(netConfig, pid)).To(Succeed())

		EnsureOutputForCommand("ip netns exec testNetNamespace ip addr", "10.10.10.10")
	})

	It("brings the veth link up", func() {
		Expect(containerConfigurer.Apply(netConfig, pid)).To(Succeed())

		stdout := gbytes.NewBuffer()
		cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip link list veth1")
		_, err := gexec.Start(cmd, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Consistently(stdout).ShouldNot(gbytes.Say("DOWN"))
	})

	It("adds a default route for network traffic", func() {
		Expect(containerConfigurer.Apply(netConfig, pid)).To(Succeed())

		EnsureOutputForCommand("ip netns exec testNetNamespace ip route show", "default via 10.10.10.1 dev veth1")
	})

	Context("when the network namespace doesn't exist", func() {
		It("returns a descriptive error message", func() {
			err := containerConfigurer.Apply(netConfig, -1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unable to find network namespace for process with pid '-1'"))
		})
	})

	Context("when the veth can't be found", func() {
		BeforeEach(func() {
			netConfig.VethNamePrefix = "vethwillnotbefound"
		})

		It("returns a descriptive error message", func() {
			err := containerConfigurer.Apply(netConfig, pid)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Container veth 'vethwillnotbefound1' not found"))
		})
	})

	Context("when the IP address cannot be assigned", func() {
		BeforeEach(func() {
			cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip addr add 10.10.10.10/24 dev veth1")
			Expect(cmd.Run()).To(Succeed())
		})

		It("returns a descriptive error message", func() {
			err := containerConfigurer.Apply(netConfig, pid)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unable to assign IP address '10.10.10.10' to veth1"))
		})
	})

	Context("when the veth link is already UP", func() {
		BeforeEach(func() {
			cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip link set veth1 up")
			Expect(cmd.Run()).To(Succeed())
		})

		It("doesn't error", func() {
			Expect(containerConfigurer.Apply(netConfig, pid)).To(Succeed())
		})
	})
})
