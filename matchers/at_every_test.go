package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("AtEvery matcher", func() {
	slice := []int{1, 2, 3}

	It("matches if the matcher is true for every item in a slice", func() {
		Expect(slice).To(AtEvery(BeNumerically(">=", 1)))
		Expect(slice).ToNot(AtEvery(Equal(4)))
	})

	It("fails if the matcher does not match for every item in a slice", func() {
		matcher := AtEvery(Equal(1))
		success, err := matcher.Match(slice)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Match failed at index 1:
Expected
    <int>: 2
to equal
    <int>: 1`))
	})

	It("fails if the slice is empty", func() {
		matcher := AtEvery(Equal(10))
		success, err := matcher.Match([]int{})
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Did not expect empty collection`))
	})

	It("errors if the type is not a slice or map", func() {
		success, err := AtEvery(Equal(1)).Match(123)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("AtEvery matcher requires an actual of type slice"))
	})
})
