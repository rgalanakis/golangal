package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("MatchField matcher", func() {
	type T struct {
		X int
	}

	It("matches if the value of the field equals the given value", func() {
		Expect(T{4}).To(MatchField("X", 4))
		Expect(T{4}).ToNot(MatchField("X", 5))
	})

	It("matches if the value of the field passes the given matcher", func() {
		Expect(T{4}).To(MatchField("X", Equal(4)))
		Expect(T{4}).ToNot(MatchField("X", Equal(5)))
	})

	It("fails if the value of the field does not match the given value", func() {
		t := T{5}
		matcher := MatchField("X", 3)
		success, err := matcher.Match(t)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(t)
		Expect(msg).To(HaveSuffix(`Field X of
    <matchers_test.T>: {X: 5}
did not match. Expected
    <int>: 5
to equal
    <int>: 3`))
	})

	It("fails if the value of the field does not match the given matcher", func() {
		t := T{5}
		matcher := MatchField("X", BeNil())
		success, err := matcher.Match(t)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(t)
		Expect(msg).To(HaveSuffix(`Field X of
    <matchers_test.T>: {X: 5}
did not match. Expected
    <int>: 5
to be nil`))
	})

	It("errors if the actual type is not a struct", func() {
		success, err := MatchField("X", Equal(1)).Match(123)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("MatchField matcher requires an actual of kind struct, not int"))
	})

	It("errors if the field does not exist", func() {
		success, err := MatchField("Y", Equal(1)).Match(T{})
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("field 'Y' does not exist on type T"))
	})
})
