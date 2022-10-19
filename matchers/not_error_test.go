package matchers_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("NotError matcher", func() {
	f := func(e error) (int, error) {
		return 1, e
	}

	fe := func(e error) error {
		return e
	}

	It("matches if the second value is not an error", func() {
		Expect(f(nil)).To(NotError())
		Expect(fe(nil)).To(NotError())
	})

	It("fails the actual is an error", func() {
		success, err := NotError().Match(errors.New("hi"))
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError(`MatchError cannot be used with functions that return an error. Got: hi`))
	})
})
