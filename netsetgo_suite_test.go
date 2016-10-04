package netsetgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
)

var pathToNetsetgo string

func TestNetsetgo(t *testing.T) {
	BeforeSuite(func() {
		var err error
		pathToNetsetgo, err = gexec.Build("github.com/teddyking/netsetgo/cmd")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Netsetgo Suite")
}
