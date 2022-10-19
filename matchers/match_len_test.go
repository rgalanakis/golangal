package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("MatchLen matcher", func() {
	slice := []int{1, 2, 3}

	It("matches a matcher against the length of an object", func() {
		Expect(slice).To(MatchLen(BeNumerically(">", 1)))
		Expect(slice).ToNot(MatchLen(Equal(1)))
	})

	It("fails if the length of the object does not match the matcher", func() {
		success, err := MatchLen(Equal(1)).Match(slice)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := MatchLen(Equal(1)).FailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Expected length of
    <[]int | len:3, cap:3>: [1, 2, 3]
to match, but failed with
    Expected
        <int>: 3
    to equal
        <int>: 1`))
	})

	It("fails if the length of the object matches the matcher but should not (negated)", func() {
		success, err := MatchLen(Equal(3)).Match(slice)
		Expect(success).To(BeTrue())
		Expect(err).ToNot(HaveOccurred())
		msg := MatchLen(Equal(3)).NegatedFailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Expected length of
    <[]int | len:3, cap:3>: [1, 2, 3]
not to match, but did with
    Expected
        <int>: 3
    not to equal
        <int>: 3`))
	})

	It("errors if the type is not a slice/string/array", func() {
		success, err := MatchLen(Equal(1)).Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError(`MatchLen matcher expects a string/array/map/channel/slice. Got:
    <int>: 5`))
	})
})
