package matchers

import (
	"github.com/onsi/gomega"
	"net/http/httptest"
)

type HaveResponseCodeMatcher struct {
	CodeOrMatcher interface{}
	inner         gomega.OmegaMatcher
	rr            *httptest.ResponseRecorder
}

func (matcher *HaveResponseCodeMatcher) Match(actual interface{}) (bool, error) {
	rr, err := requireRespRec(actual)
	if err != nil {
		return false, err
	}
	if inner, ok := matcher.CodeOrMatcher.(gomega.OmegaMatcher); ok {
		matcher.inner = inner
	} else {
		matcher.inner = gomega.BeEquivalentTo(matcher.CodeOrMatcher)
	}
	matcher.rr = rr
	return matcher.inner.Match(rr.Code)
}

func (matcher *HaveResponseCodeMatcher) FailureMessage(actual interface{}) (message string) {
	msg := matcher.inner.FailureMessage(matcher.rr.Code)
	//msg := format.Message(matcher.rr.Code, "to equal", matcher.CodeOrMatcher)
	msg += "\nBody:\n" + matcher.rr.Body.String()
	return msg
}

func (matcher *HaveResponseCodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic(noNegate)
}
