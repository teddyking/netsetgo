package netsetgo_test

import (
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
})
