package golangal_test

import (
	. "github.com/onsi/ginkgo/v2"
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

var _ = Describe("Envvars", func() {
	addEnvVar := golangal.EnvVars()

	It("cleans up env vars (check before)", func() {
		_, exists := os.LookupEnv("SOME_VAR")
		Expect(exists).To(BeFalse())
	})

	It("sets the env var", func() {
		addEnvVar("SOME_VAR", "X")
		Expect(os.Getenv("SOME_VAR")).To(Equal("X"))
	})

	It("cleans up env vars (check after)", func() {
		_, exists := os.LookupEnv("SOME_VAR")
		Expect(exists).To(BeFalse())
	})

	var originalUserValue string
	It("restores original values (check some var exists)", func() {
		originalUserValue = os.Getenv("USER")
		Expect(originalUserValue).ToNot(Equal(""))
	})

	It("restores original values (replace existing var)", func() {
		addEnvVar("USER", "golangal-test-temp")
	})

	It("restores original values (check some var exists)", func() {
		Expect(os.Getenv("USER")).To(Equal(originalUserValue))
	})
})
