package matchers

import (
	"github.com/onsi/gomega/types"
	"net/http/httptest"
	"sort"
	"strings"
)

type HaveHeaderMatcher struct {
	Key   string
	Inner types.GomegaMatcher
	rr    *httptest.ResponseRecorder
	got   string
}

func (matcher *HaveHeaderMatcher) Match(actual interface{}) (bool, error) {
	rr, err := requireRespRec(actual)
	if err != nil {
		return false, err
	}
	matcher.rr = rr
	matcher.got = rr.Header().Get(matcher.Key)
	return matcher.Inner.Match(matcher.got)
}

func (matcher *HaveHeaderMatcher) FailureMessage(actual interface{}) (message string) {
	bld := &strings.Builder{}
	if matcher.got == "" {
		bld.WriteString(matcher.Key)
		bld.WriteString(" is missing\nFound: ")
		headers := make([]string, 0, len(matcher.rr.Header()))
		for k := range matcher.rr.Header() {
			headers = append(headers, k)
		}
		sort.Strings(headers)
		bld.WriteString(strings.Join(headers, ", "))
	} else {
		bld.WriteString(matcher.Key)
		bld.WriteString(": ")
		bld.WriteString(matcher.Inner.FailureMessage(matcher.got))
	}
	return bld.String()
}

func (matcher *HaveHeaderMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic(noNegate)
}
