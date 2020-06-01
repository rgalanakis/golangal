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

// HaveResponseCode is a Gomega matcher to ensure an
// *httptest.ResponseRecorder has the expected response code.
// If it does not, the actual code and body are printed
// (the body is really useful information when tests fail).
func HaveResponseCode(codeOrMatcher interface{}) types.GomegaMatcher {
	return &matchers.HaveResponseCodeMatcher{CodeOrMatcher: codeOrMatcher}
}

// PanicWith succeeds if actual is a function that, when invoked, panics.
// The panic message must match the given panicMatcher,
// which is matched against the object used for the panic (usually a string or error).
// Actual must be a function that takes no arguments and returns no results.
func PanicWith(panicMatcher types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.PanicWith{PanicMatcher: panicMatcher}
}

// AtEvery succeeds when every element in a slice matches the given matcher.
// Used to assert an expectation against every element in a collection.
func AtEvery(m types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.AtEveryMatcher{Matcher: m}
}

// AtIndex succeeds when the slice element at the given index matches the given matcher.
func AtIndex(idx int, m types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.AtIndexMatcher{Index: idx, Matcher: m}
}

// AtKey succeeds when the element of a slice at the given key matches the given matcher.
func AtKey(key interface{}, m types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.AtKeyMatcher{Key: key, Matcher: m}
}

// MatchLen matches the length of a collection against a matcher.
// It's like HaveLen, but allows a dynamic length.
// HaveLen(2) would be equivalent to MatchLen(Equal(2)).
func MatchLen(m types.GomegaMatcher) types.GomegaMatcher {
	return &matchers.MatchLenMatcher{Matcher: m}
}

// MatchField matches the the value of the field named name on the actual struct.
// If m is a gomega matcher, it is matched against the field value.
// Otherwise, test against field equality.
//
//     Expect(MyStruct{Field1: 10}).To(MatchField("Field1", 10))
//     Expect(MyStruct{Field1: 10}).To(MatchField("Field1", BeNumerically(">", 5)))
//
// You can match multiple fields by using SatisfyAll or And:
//
//     o := MyStruct{Field1: 10, Field2: true}
//     Expect(o).To(SatisfyAll(MatchField("Field1", 10), MatchField("Field2", BeTrue())))
//
func MatchField(name string, m interface{}) types.GomegaMatcher {
	var matcher types.GomegaMatcher
	if casted, ok := m.(types.GomegaMatcher); ok {
		matcher = casted
	} else {
		matcher = gomega.Equal(m)
	}
	return &matchers.MatchFieldMatcher{Name: name, Matcher: matcher}
}
