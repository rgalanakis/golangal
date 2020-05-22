package matchers

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/types"
	"net/http/httptest"
)

type HaveJsonBodyMatcher struct {
	Inner          types.GomegaMatcher
	rr             *httptest.ResponseRecorder
	actualBodyJson interface{}
	decodeErr      error
}

func (matcher *HaveJsonBodyMatcher) Match(actual interface{}) (bool, error) {
	rr, err := requireRespRec(actual)
	if err != nil {
		return false, err
	}
	matcher.rr = rr
	if err := json.Unmarshal(rr.Body.Bytes(), &matcher.actualBodyJson); err != nil {
		matcher.decodeErr = err
		return false, nil
	}
	return matcher.Inner.Match(matcher.actualBodyJson)
}

func (matcher *HaveJsonBodyMatcher) FailureMessage(actual interface{}) (message string) {
	if matcher.decodeErr != nil {
		return fmt.Sprintf("Error decoding body: %+v", matcher.decodeErr)
	}
	return matcher.Inner.FailureMessage(matcher.actualBodyJson)
}

func (matcher *HaveJsonBodyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic(noNegate)
}
