package netsetgo_test

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo binary", func() {
	It("exits with a status code 0", func() {
		command := exec.Command(pathToNetsetgo)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
	})

	Context("when passed --bridgeName <NAME>", func() {
		It("creates a bridge device on the host with the provided name", func() {
			command := exec.Command(pathToNetsetgo, "--bridgeName=tower")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			_, err = net.InterfaceByName("tower")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
