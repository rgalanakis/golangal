package golangal

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/rgalanakis/golangal/internal"
	"github.com/rgalanakis/golangal/matchers"
	"io/ioutil"
	"os"
)

type ginkgoHook func(args ...interface{}) bool

// EachTempDir returns a function that returns a temporary directory string when invoked
// that is cached across a single test.
// The temporary directory is cleaned up after the test has finished.
//
// Example:
//
//	tempdir := golangal.EachTempDir()
//	It("writes a file", func() {
//	  path := filepath.Join(tempdir(), "some-file.txt")
//	  ...
//	})
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
//	addEnvVar := golangal.EnvVars()
//	It("reads the env", func() {
//	  addEnvVar("MY_CONFIG_VAR", "XYZ")
//	  ...
//	})
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

func HaveHeader(key string, m interface{}) gomega.OmegaMatcher {
	return &matchers.HaveHeaderMatcher{Key: key, Inner: internal.CoerceToMatcher(m)}
}

func HaveJsonBody(m interface{}) gomega.OmegaMatcher {
	return &matchers.HaveJsonBodyMatcher{Inner: internal.CoerceToMatcher(m)}
}

// HaveResponseCode is a Gomega matcher to ensure an
// *httptest.ResponseRecorder has the expected response code.
// If it does not, the actual code and body are printed
// (the body is really useful information when tests fail).
func HaveResponseCode(codeOrMatcher interface{}) gomega.OmegaMatcher {
	return &matchers.HaveResponseCodeMatcher{CodeOrMatcher: codeOrMatcher}
}

// AtEvery succeeds when every element in a slice matches the given matcher.
// Used to assert an expectation against every element in a collection.
func AtEvery(m interface{}) gomega.OmegaMatcher {
	return &matchers.AtEveryMatcher{Matcher: internal.CoerceToMatcher(m)}
}

// AtIndex succeeds when the slice element at the given index matches the given matcher.
func AtIndex(idx int, m interface{}) gomega.OmegaMatcher {
	return &matchers.AtIndexMatcher{Index: idx, Matcher: internal.CoerceToMatcher(m)}
}

// AtKey succeeds when the element of a slice at the given key matches the given matcher.
func AtKey(key interface{}, m interface{}) gomega.OmegaMatcher {
	return &matchers.AtKeyMatcher{Key: key, Matcher: internal.CoerceToMatcher(m)}
}

// MatchLen matches the length of a collection against a matcher.
// It's like HaveLen, but allows a dynamic length.
// HaveLen(2) would be equivalent to MatchLen(Equal(2)).
// MatchLen also works with any type with a Len() int method.
func MatchLen(m interface{}) gomega.OmegaMatcher {
	return &matchers.MatchLenMatcher{Matcher: internal.CoerceToMatcher(m)}
}

// MatchCap matches the length of a collection against a matcher.
// It's like HaveCap, but allows a dynamic length.
// HaveCap(2) would be equivalent to MatchCap(Equal(2)).
// MatchCap also works with any type with a Cap() int method.
func MatchCap(m interface{}) gomega.OmegaMatcher {
	return &matchers.MatchCapMatcher{Matcher: internal.CoerceToMatcher(m)}
}

// MatchField matches the value of the field named name on the actual struct.
// If m is a gomega matcher, it is matched against the field value.
// Otherwise, test against field equality.
//
//	Expect(MyStruct{Field1: 10}).To(MatchField("Field1", 10))
//	Expect(MyStruct{Field1: 10}).To(MatchField("Field1", BeNumerically(">", 5)))
//
// You can match multiple fields by using SatisfyAll or And:
//
//	o := MyStruct{Field1: 10, Field2: true}
//	Expect(o).To(SatisfyAll(MatchField("Field1", 10), MatchField("Field2", BeTrue())))
func MatchField(name string, m interface{}) gomega.OmegaMatcher {
	return &matchers.MatchFieldMatcher{Name: name, Matcher: internal.CoerceToMatcher(m)}
}

// MatchPtrField is same as MatchField, but the actual object should be a pointer, not a struct.
func MatchPtrField(name string, m interface{}) gomega.OmegaMatcher {
	return &matchers.MatchPtrFieldMatcher{Name: name, Matcher: internal.CoerceToMatcher(m)}
}

// NotError is like gomega's Succeed matcher, except it handles functions which
// return multiple values. The docs say this:
//
// You should not use a function with multiple return values (like DoSomethingHard) with Succeed.
// Matchers are only passed the first value provided to Ω/Expect,
// the subsequent arguments are handled by Ω and Expect as outlined above.
// As a result of this behavior Ω(DoSomethingHard()).ShouldNot(Succeed()) would never pass.
//
// We can circumvent this by using a noop matcher, and allowing Gomega's internal Expect behavior
// to assert that an error wasn't returned (this matcher will fail if it receives an error).
//
//	Expect(DoSomethingHard()).To(NotError())
func NotError() gomega.OmegaMatcher {
	return &matchers.NotErrorMatcher{}
}
