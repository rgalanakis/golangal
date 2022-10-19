package matchers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
	"github.com/rgalanakis/golangal/matchers"
	"net/http/httptest"
)

var _ = Describe("HaveHeaderMatcher", func() {
	var success bool
	var err error
	It("can match positives", func() {
		resp := &httptest.ResponseRecorder{}
		resp.Header().Add("Test-Header", "somestring")
		Expect(resp).To(HaveHeader("Test-Header", ContainSubstring("somestring")))
	})
	It("can match negatives", func() {
		resp := &httptest.ResponseRecorder{}
		resp.Header().Add("Test-Header", "somestring")
		Expect(resp).ToNot(HaveHeader("Test-Header", ContainSubstring("otherstring")))
	})
	It("errors for an invalid actual", func() {
		success, err = (&matchers.HaveHeaderMatcher{}).Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("actual must be a *httptest.ResponseRecorder"))
	})
	It("fails if inner matcher does not match", func() {
		resp := &httptest.ResponseRecorder{}
		resp.Header().Add("Test-Header", "somestring")
		matcher := &matchers.HaveHeaderMatcher{Key: "Test-Header", Inner: ContainSubstring("otherstring")}
		success, err = matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Test-Header: Expected
    <string>: somestring
to contain substring
    <string>: otherstring`))
	})
	It("fails if header is not present", func() {
		resp := &httptest.ResponseRecorder{}
		resp.Header().Add("Test-Header", "s")
		resp.Header().Add("Another-Header", "s")
		matcher := &matchers.HaveHeaderMatcher{Key: "Other-Header", Inner: ContainSubstring("s")}
		success, err = matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Other-Header is missing
Found: Another-Header, Test-Header`))
	})
})
