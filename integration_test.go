package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"net"
	"os/exec"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo binary", func() {
	Context("when passed all required args", func() {
		BeforeEach(func() {
			command := exec.Command(pathToNetsetgo, "--bridgeName=tower", "--bridgeAddress=10.10.10.1/24", "--vethNamePrefix=v")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			cmd := exec.Command("sh", "-c", "ip link delete tower")
			Expect(cmd.Run()).To(Succeed())
			cmd = exec.Command("sh", "-c", "ip link delete v0") // will implicitly delete v1 :D
			Expect(cmd.Run()).To(Succeed())
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

		It("sets the bridge link up", func() {
			Expect("/sys/class/net/tower/carrier").To(BeAnExistingFile())
			carrierFileContents, err := ioutil.ReadFile("/sys/class/net/tower/carrier")
			Expect(err).NotTo(HaveOccurred())
			Eventually(string(carrierFileContents)).Should(Equal("1\n"))
		})

		It("creates a veth pair on the host using the provided name prefix", func() {
			_, err := net.InterfaceByName("v0")
			Expect(err).NotTo(HaveOccurred())
			_, err = net.InterfaceByName("v1")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
