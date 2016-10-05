package netsetgo_suite_helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func CreateNetNamespace(name string) {
	cmdString := fmt.Sprintf("ip netns add %s", name)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())
}

func CreateVethInNetNamespace(vethNamePrefix, netNamespaceName string) {
	cmdString := fmt.Sprintf("ip link add %s0 type veth peer name %s1", vethNamePrefix, vethNamePrefix)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())

	cmdString = fmt.Sprintf("ip link set %s1 netns %s", vethNamePrefix, netNamespaceName)
	cmd = exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())

	cmdString = fmt.Sprintf("ip link set %s0 up", vethNamePrefix)
	cmd = exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())
}

func DestroyBridge(name string) {
	cmdString := fmt.Sprintf("ip link delete %s", name)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())
}

func DestroyNetNamespace(name string) {
	cmdString := fmt.Sprintf("ip netns delete %s", name)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())
}

func DestroyVeth(namePrefix string) {
	cmdString := fmt.Sprintf("ip link delete %s0", namePrefix)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Run()).To(Succeed())
}

func KillCmd(pid int) {
	process, err := os.FindProcess(pid)
	Expect(err).NotTo(HaveOccurred())

	Expect(process.Kill()).To(Succeed())
}

func RunCmdInNetNamespace(netNamespaceName string, cmdPathAndArgs string) (int, int) {
	cmdString := fmt.Sprintf("ip netns exec %s %s", netNamespaceName, cmdPathAndArgs)
	cmd := exec.Command("sh", "-c", cmdString)
	Expect(cmd.Start()).To(Succeed())

	parentPid := cmd.Process.Pid

	// super gross
	cmd = exec.Command("sh", "-c", fmt.Sprintf("ps --ppid %d | tail -n 1 | awk '{print $1}'", parentPid))
	pidBytes, err := cmd.Output()
	Expect(err).NotTo(HaveOccurred())

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	Expect(err).NotTo(HaveOccurred())

	return parentPid, pid
}

func EnsureOutputForCommand(command, expectedOutput string) {
	stdout := gbytes.NewBuffer()
	cmd := exec.Command("sh", "-c", command)
	_, err := gexec.Start(cmd, stdout, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(stdout).Should(gbytes.Say(expectedOutput))
}
