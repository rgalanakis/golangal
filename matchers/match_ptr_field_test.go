package matchers_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("MatchPtrField matcher", func() {
	type T struct {
		X int
	}

	It("matches if the value of the field equals the given value", func() {
		Expect(&T{4}).To(MatchPtrField("X", 4))
		Expect(&T{4}).ToNot(MatchPtrField("X", 5))
	})

	It("matches if the value of the field passes the given matcher", func() {
		Expect(&T{4}).To(MatchPtrField("X", Equal(4)))
		Expect(&T{4}).ToNot(MatchPtrField("X", Equal(5)))
	})

	It("fails if the value of the field does not match the given value", func() {
		t := &T{5}
		matcher := MatchPtrField("X", 3)
		success, err := matcher.Match(t)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(t)
		Expect(msg).To(HaveSuffix(fmt.Sprintf(`Field X of
    <*matchers_test.T | %p>: {X: 5}
did not match. Expected
    <int>: 5
to equal
    <int>: 3`, t)))
	})

	It("fails if the value of the field does not match the given matcher", func() {
		t := &T{5}
		matcher := MatchPtrField("X", BeNil())
		success, err := matcher.Match(t)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(t)
		Expect(msg).To(HaveSuffix(fmt.Sprintf(`Field X of
    <*matchers_test.T | %p>: {X: 5}
did not match. Expected
    <int>: 5
to be nil`, t)))
	})

	It("errors if the actual type is not a pointer", func() {
		success, err := MatchPtrField("X", Equal(1)).Match(123)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("MatchPtrField matcher requires an actual of kind ptr, not int"))

		success, err = MatchPtrField("X", Equal(1)).Match(T{})
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("MatchPtrField matcher requires an actual of kind ptr, not struct"))
	})

	It("errors if the field does not exist", func() {
		success, err := MatchPtrField("Y", Equal(1)).Match(&T{})
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("field 'Y' does not exist on type T"))
	})
})
