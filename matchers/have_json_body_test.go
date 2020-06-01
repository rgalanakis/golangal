package matchers_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/rgalanakis/golangal"
	"github.com/rgalanakis/golangal/matchers"
	"net/http/httptest"
)

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
