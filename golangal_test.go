package golangal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rgalanakis/golangal"
	"os"
	"testing"
)

func TestGolangal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "golangal package Suite")
}

var _ = Describe("EachTempDir and SuiteTempDir", func() {
	tempdir := golangal.EachTempDir()

	It("creates a temporary directory", func() {
		Expect(tempdir()).To(HavePrefix(os.TempDir()))
	})
})
