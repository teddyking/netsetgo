package configurer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestConfigurer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configurer Suite")
}

func createTestBridge() {
	cmd := exec.Command("sh", "-c", "ip link add name tower type bridge")
	Expect(cmd.Run()).To(Succeed())
}

func createTestVeth() {
	cmd := exec.Command("sh", "-c", "ip link add veth0 type veth peer name veth1")
	Expect(cmd.Run()).To(Succeed())
}

func createTestNetNamespace() {
	cmd := exec.Command("sh", "-c", "ip netns add testNetNamespace")
	Expect(cmd.Run()).To(Succeed())
}

func cleanupTestBridge() {
	cmd := exec.Command("sh", "-c", "ip link delete tower")
	Expect(cmd.Run()).To(Succeed())
}

func cleanupTestVeth() {
	cmd := exec.Command("sh", "-c", "ip link delete veth0") // will implicitly delete veth1 :D
	Expect(cmd.Run()).To(Succeed())
}

func cleanupTestNetNamespace() {
	cmd := exec.Command("sh", "-c", "ip netns delete testNetNamespace")
	Expect(cmd.Run()).To(Succeed())
}

func addTestIPToTestBridge() {
	cmd := exec.Command("sh", "-c", "ip addr add 10.10.10.1/24 dev tower")
	Expect(cmd.Run()).To(Succeed())
}

func setTestBridgeUp() {
	cmd := exec.Command("sh", "-c", "ip link set tower up")
	Expect(cmd.Run()).To(Succeed())
}

func runCmdInTestNetNamespace() (int, int) {
	cmd := exec.Command("sh", "-c", "ip netns exec testNetNamespace sleep 1000")
	Expect(cmd.Start()).To(Succeed())

	parentPid := cmd.Process.Pid

	cmd = exec.Command("sh", "-c", fmt.Sprintf("ps --ppid %d | tail -n 1 | awk '{print $1}'", parentPid))
	pidBytes, err := cmd.Output()
	Expect(err).NotTo(HaveOccurred())

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	Expect(err).NotTo(HaveOccurred())

	return parentPid, pid
}

func killCmdInTestNetNamespace(pid int) {
	process, err := os.FindProcess(pid)
	Expect(err).NotTo(HaveOccurred())

	Expect(process.Kill()).To(Succeed())
}
