package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("AtIndex matcher", func() {
	slice := []int{1, 2, 3}

	It("matches a matcher against the item at an index in a slice", func() {
		Expect(slice).To(AtIndex(0, Equal(1)))
		Expect(slice).ToNot(AtIndex(1, Equal(1)))
	})

	It("fails if the value at the given index does not match", func() {
		matcher := AtIndex(0, Equal(10))
		success, err := matcher.Match(slice)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Matcher failed at slice index 0. Expected
    <int>: 1
to equal
    <int>: 10`))
	})

	It("fails if the index is out of bounds", func() {
		matcher := AtIndex(5, Equal(1))
		success, err := matcher.Match(slice)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(slice)).To(HaveSuffix(`Slice
    <[]int | len:3, cap:3>: [1, 2, 3]
is too short to match against index 5`))
	})

	It("errors if the type is not a slice/string/array", func() {
		success, err := AtIndex(0, Equal(1)).Match(123)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("AtIndex matcher requires an actual of string, slice, or array"))
	})
})
