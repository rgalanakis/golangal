package golangal

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

type ginkgoHook func(body interface{}, timeout ...float64) bool

// EachTempDir returns a function that returns a temporary directory string when invoked
// that is cached across a single test.
// The temporary directory is cleaned up after the test has finished.
func EachTempDir() func() string {
	return eachTempDir("each", ginkgo.BeforeEach, ginkgo.AfterEach)
}

func eachTempDir(scope string, before, after ginkgoHook) func() string {
	var tempdir string
	before(func() {
		td, err := ioutil.TempDir("", "galangal-"+scope)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		tempdir = td
	})
	after(func() {
		gomega.Expect(os.RemoveAll(tempdir)).To(gomega.Succeed())
	})
	return func() string {
		return tempdir
	}
}
