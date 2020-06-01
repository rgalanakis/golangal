package matchers_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
	"github.com/rgalanakis/golangal/matchers"
	"net/http/httptest"
)

var _ = Describe("HaveResponseCodeMatcher", func() {
	It("can match a code", func() {
		rr := &httptest.ResponseRecorder{Code: 422}
		Expect(rr).To(HaveResponseCode(422))
	})
	It("can match a matcher", func() {
		rr := &httptest.ResponseRecorder{Code: 422}
		Expect(rr).To(HaveResponseCode(BeNumerically(">", 400)))
	})
	It("errors for an invalid actual", func() {
		success, err := (&matchers.HaveResponseCodeMatcher{}).Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("actual must be a *httptest.ResponseRecorder"))
	})
	It("fails with a clear message", func() {
		resp := &httptest.ResponseRecorder{Code: 422, Body: bytes.NewBufferString("abc")}
		matcher := &matchers.HaveResponseCodeMatcher{CodeOrMatcher: 301}
		success, err := matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Expected
    <int>: 422
to be equivalent to
    <int>: 301
Body:
abc`))
	})
	It("can fail with a matcher", func() {
		resp := &httptest.ResponseRecorder{Code: 422}
		matcher := &matchers.HaveResponseCodeMatcher{CodeOrMatcher: BeNumerically("<", 400)}
		success, err := matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Expected
    <int>: 422
to be <
    <int>: 400
Body:
<nil>`))
	})
})
