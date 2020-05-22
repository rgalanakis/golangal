package matchers_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
	"github.com/rgalanakis/golangal/matchers"
	"net/http/httptest"
	"testing"
)

func TestMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "golangal matchers package Suite")
}

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

var _ = Describe("HaveJsonBodyMatcher", func() {
	newRr := func(body string) *httptest.ResponseRecorder {
		rr := &httptest.ResponseRecorder{
			Body: bytes.NewBufferString(body),
		}
		return rr
	}
	It("can match a matcher", func() {
		resp := newRr(`{"a": 1}`)
		Expect(resp).To(HaveJsonBody(HaveKeyWithValue("a", BeEquivalentTo(1))))
	})
	It("errors for an invalid actual", func() {
		success, err := (&matchers.HaveJsonBodyMatcher{}).Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("actual must be a *httptest.ResponseRecorder"))
	})
	It("fails if the matcher does not match the body", func() {
		resp := newRr(`{"a": 1}`)
		matcher := &matchers.HaveJsonBodyMatcher{Inner: HaveKeyWithValue("b", BeEquivalentTo(2))}
		success, err := matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Expected
    <map[string]interface {} | len:1>: {"a": 1}
to have {key: value} matching
    <map[interface {}]interface {} | len:1>: {"b": {Expected: 2}}`))
	})
	It("fails if the body is not valid json", func() {
		resp := newRr(`{"a`)
		matcher := &matchers.HaveJsonBodyMatcher{}
		success, err := matcher.Match(resp)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		Expect(matcher.FailureMessage(resp)).To(HavePrefix(`Error decoding body: unexpected end of JSON input`))
	})
})

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
