package matchers_test

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("MatchCap matcher", func() {
	slice := make([]int, 1, 3)

	It("matches a matcher against the capacity of an object", func() {
		Expect(slice).To(MatchCap(BeNumerically(">", 1)))
		Expect(slice).ToNot(MatchCap(Equal(1)))
	})

	It("fails if the capacity of the object does not match the matcher", func() {
		success, err := MatchCap(Equal(1)).Match(slice)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := MatchCap(Equal(1)).FailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Expected capacity of
    <[]int | len:1, cap:3>: [0]
to match, but failed with
    Expected
        <int>: 3
    to equal
        <int>: 1`))
	})

	It("fails if the capacity of the object matches the matcher but should not (negated)", func() {
		success, err := MatchCap(Equal(3)).Match(slice)
		Expect(success).To(BeTrue())
		Expect(err).ToNot(HaveOccurred())
		msg := MatchCap(Equal(3)).NegatedFailureMessage(slice)
		Expect(msg).To(HaveSuffix(`Expected capacity of
    <[]int | len:1, cap:3>: [0]
not to match, but did with
    Expected
        <int>: 3
    not to equal
        <int>: 3`))
	})

	It("matches a matcher against the capacity of an object with a Cap method", func() {
		Expect(CustomCap{3}).To(MatchCap(3))
		Expect(slice).ToNot(MatchCap(Equal(1)))
		Expect(bytes.NewBufferString("abc")).To(MatchCap(8))
		Expect(bytes.NewBufferString("abc")).ToNot(MatchCap(7))
	})

	It("errors if the type is not a slice/string/array and is not a hasCap", func() {
		success, err := MatchCap(Equal(1)).Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError(`MatchCap matcher expects a string/array/map/channel/slice, or type with a Cap() int method. Got:
    <int>: 5`))
	})
})

type CustomCap struct {
	cap int
}

func (cl CustomCap) Cap() int {
	return cl.cap
}
