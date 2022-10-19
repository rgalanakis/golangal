package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
)

var _ = Describe("AtKey matcher", func() {
	subject := map[string]int{"a": 5, "b": 6}

	It("matches a matcher against the value at the given key", func() {
		Expect(subject).To(AtKey("a", Equal(5)))
		Expect(subject).ToNot(AtKey("a", Equal(6)))
	})

	It("errors if the type is not a map", func() {
		success, err := AtKey(0, Equal(1)).Match(123)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("AtKey matcher requires an actual of map"))
	})

	It("fails if the value at the given key does not match", func() {
		matcher := AtKey("a", Equal(1))
		success, err := matcher.Match(subject)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(subject)
		Expect(msg).To(HaveSuffix(`Matcher failed at map key <string>: "a". Expected
    <int>: 5
to equal
    <int>: 1`))
	})

	It("fails if the key is not present", func() {
		matcher := AtKey("c", Equal(1))
		success, err := matcher.Match(subject)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(subject)).To(MatchRegexp(`Map
    <map[string]int | len:2>: .*
does not contain key
    <string>: c`))
	})
})
