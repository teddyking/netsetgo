package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/teddyking/netsetgo"

	"io/ioutil"
	"net"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("netsetgo", func() {
	Describe("CreateBridge", func() {
		AfterEach(func() {
			cmd := exec.Command("sh", "-c", "ip link delete tower")
			Expect(cmd.Run()).To(Succeed())
		})

		Context("when a device with the provided name doesn't already exist", func() {
			It("creates a bridge device with the provided name", func() {
				err := CreateBridge("tower")
				Expect(err).NotTo(HaveOccurred())

				_, err = net.InterfaceByName("tower")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a device with the provided name already exists", func() {
			BeforeEach(func() {
				cmd := exec.Command("sh", "-c", "ip link add name tower type bridge")
				Expect(cmd.Run()).To(Succeed())
			})

			It("doesn't error", func() {
				err := CreateBridge("tower")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("AddAddressToBridge", func() {
		Context("when the bridge exists", func() {
			BeforeEach(func() {
				cmd := exec.Command("sh", "-c", "ip link add name tower type bridge")
				Expect(cmd.Run()).To(Succeed())
			})

			AfterEach(func() {
				cmd := exec.Command("sh", "-c", "ip link delete tower")
				Expect(cmd.Run()).To(Succeed())
			})

			Context("when the address is valid", func() {
				It("adds the provided address to the provided bridge", func() {
					err := AddAddressToBridge("tower", "10.10.10.1/24")
					Expect(err).NotTo(HaveOccurred())

					stdout := gbytes.NewBuffer()
					command := exec.Command("sh", "-c", "ip addr show tower")
					session, err := gexec.Start(command, stdout, GinkgoWriter)
					Expect(err).NotTo(HaveOccurred())
					Eventually(session).Should(gexec.Exit(0))
					Eventually(stdout).Should(gbytes.Say("10.10.10.1/24"))
				})
			})

			Context("when the address isn't valid", func() {
				It("returns an error", func() {
					err := AddAddressToBridge("tower", "10.10.10.1")
					Expect(err.Error()).To(ContainSubstring("invalid CIDR address"))
				})
			})
		})

		Context("when the bridge doesn't exist", func() {
			It("returns an error", func() {
				err := AddAddressToBridge("tower", "10.10.10.1/24")
				Expect(err.Error()).To(ContainSubstring("no such device"))
			})
		})
	})

	Describe("SetBridgeUp", func() {
		Context("when the bridge exists and has an address", func() {
			BeforeEach(func() {
				cmd := exec.Command("sh", "-c", "ip link add name tower type bridge")
				Expect(cmd.Run()).To(Succeed())
				cmd = exec.Command("sh", "-c", "ip addr add 10.10.10.1/24 dev tower")
				Expect(cmd.Run()).To(Succeed())
			})

			AfterEach(func() {
				cmd := exec.Command("sh", "-c", "ip link delete tower")
				Expect(cmd.Run()).To(Succeed())
			})

			Context("when the bridge isn't already up", func() {
				It("brings the bridge up", func() {
					err := SetBridgeUp("tower")
					Expect(err).NotTo(HaveOccurred())

					Expect("/sys/class/net/tower/carrier").To(BeAnExistingFile())
					carrierFileContents, err := ioutil.ReadFile("/sys/class/net/tower/carrier")
					Expect(err).NotTo(HaveOccurred())
					Expect(string(carrierFileContents)).To(Equal("1\n"))
				})
			})

			Context("when the bridge is already up", func() {
				BeforeEach(func() {
					cmd := exec.Command("sh", "-c", "ip link set tower up")
					Expect(cmd.Run()).To(Succeed())
				})

				It("doesn't error", func() {
					err := SetBridgeUp("tower")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the bridge doesn't exist", func() {
			It("returns an error", func() {
				err := SetBridgeUp("tower")
				Expect(err.Error()).To(ContainSubstring("no such device"))
			})
		})
	})

	Describe("CreateVethPair", func() {
		AfterEach(func() {
			cmd := exec.Command("sh", "-c", "ip link delete veth0") // will implicitly delete veth1 :D
			Expect(cmd.Run()).To(Succeed())
		})

		Context("when a veth pair with the provided name prefix doesn't already exist", func() {
			It("creates a veth pair using the provided name prefix", func() {
				err := CreateVethPair("veth")
				Expect(err).NotTo(HaveOccurred())

				_, err = net.InterfaceByName("veth0")
				Expect(err).NotTo(HaveOccurred())
				_, err = net.InterfaceByName("veth1")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a veth pair with the provided name prefix already exists", func() {
			BeforeEach(func() {
				cmd := exec.Command("sh", "-c", "ip link add veth0 type veth peer name veth1")
				Expect(cmd.Run()).To(Succeed())
			})

			It("doesn't error", func() {
				err := CreateVethPair("veth")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})
