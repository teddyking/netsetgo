package configurer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo/configurer"

	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = PDescribe("ContainerConfigurer", func() {
	var (
		parentPid, pid          int
		configurer              *ContainerConfigurer
		testCmd, invalidTestCmd *exec.Cmd
	)

	BeforeEach(func() {
		createTestNetNamespace()
		parentPid, pid = runCmdInTestNetNamespace()

		testCmd = exec.Command("sh", "-c", "ip link add name testdevice type bridge")
		invalidTestCmd = exec.Command("sh", "-c", "notavalidcommand")
	})

	JustBeforeEach(func() {
		configurer = New(pid)
	})

	AfterEach(func() {
		killCmdInTestNetNamespace(parentPid)
		cleanupTestNetNamespace()
	})

	It("runs commands inside the network namespace of the process identified by the provided pid", func() {
		err := configurer.Exec(testCmd)
		Expect(err).NotTo(HaveOccurred())

		stdout := gbytes.NewBuffer()
		cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace ip link ls")
		_, err = gexec.Start(cmd, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(stdout).Should(gbytes.Say("testdevice"))
	})

	It("assigns the provided address to the veth inside the network namespace of the process identified by the provided pid", func() {

	})

	Context("when the network namespace identified by the provided pid doesn't exist", func() {
		BeforeEach(func() {
			pid = -1
		})

		It("returns a descriptive error", func() {
			err := configurer.Exec(testCmd)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no such file or directory"))
		})
	})

	Context("when the provided *exec.Cmd returns an error", func() {
		It("returns a descriptive error", func() {
			err := configurer.Exec(invalidTestCmd)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("exit status 127"))
		})
	})
})
