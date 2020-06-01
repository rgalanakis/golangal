package matchers_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
	"github.com/rgalanakis/golangal/matchers"
)

var _ = Describe("PanicWith matcher", func() {

	It("matches panic messages", func() {
		Expect(func() {
			panic("uh oh")
		}).To(PanicWith(Equal("uh oh")))
	})

	It("matches panic objects", func() {
		Expect(func() {
			panic(errors.New("some error"))
		}).To(PanicWith(MatchError(HavePrefix("some "))))
	})

	It("fails if the function does not panic", func() {
		m := &matchers.PanicWith{PanicMatcher: MatchError(HavePrefix("some "))}
		success, err := m.Match(func() {})
		Expect(err).To(Not(HaveOccurred()))
		Expect(success).To(BeFalse())
	})

	It("fails if the panic does not match the matcher", func() {
		m := &matchers.PanicWith{PanicMatcher: Equal("this")}
		success, err := m.Match(func() {
			panic("that")
		})
		Expect(err).To(Not(HaveOccurred()))
		Expect(success).To(BeFalse())
	})
})
