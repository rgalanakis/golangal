package golangal

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/rgalanakis/golangal/matchers"
	"io/ioutil"
	"os"
)

type ginkgoHook func(body interface{}, timeout ...float64) bool

// EachTempDir returns a function that returns a temporary directory string when invoked
// that is cached across a single test.
// The temporary directory is cleaned up after the test has finished.
//
// Example:
//
//    tempdir := golangal.EachTempDir()
//    It("writes a file", func() {
//      path := filepath.Join(tempdir(), "some-file.txt")
//      ...
//    })
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

// EnvVars returns a function that can be called to set an environment variable for the duration of the test,
// and unset it (if it did not exist) or restore the original value (if it did exist) after the test.
//
// Example:
//
//    addEnvVar := golangal.EnvVars()
//    It("reads the env", func() {
//      addEnvVar("MY_CONFIG_VAR", "XYZ")
//      ...
//    })
func EnvVars() func(key, value string) {
	var envVars map[string]*string
	ginkgo.BeforeEach(func() {
		envVars = make(map[string]*string, 4)
	})
	ginkgo.AfterEach(func() {
		for k, v := range envVars {
			if v == nil {
				gomega.Expect(os.Unsetenv(k)).To(gomega.Succeed())
			} else {
				gomega.Expect(os.Setenv(k, *v)).To(gomega.Succeed())
			}
		}
	})
	return func(k, v string) {
		orig, exists := os.LookupEnv(k)
		if exists {
			envVars[k] = &orig
		} else {
			envVars[k] = nil
		}
		gomega.Expect(os.Setenv(k, v)).To(gomega.Succeed())
	}
}

func HaveHeader(key string, inner types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.HaveHeaderMatcher{Key: key, Inner: inner}
}

func HaveJsonBody(inner types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.HaveJsonBodyMatcher{Inner: inner}
}
func HaveResponseCode(codeOrMatcher interface{}) types.GomegaMatcher {
	return &matchers.HaveResponseCodeMatcher{CodeOrMatcher: codeOrMatcher}
}
