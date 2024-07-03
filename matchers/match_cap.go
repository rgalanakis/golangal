package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type MatchCapMatcher struct {
	Matcher types.GomegaMatcher
}

func (matcher *MatchCapMatcher) Match(actual interface{}) (bool, error) {
	capacity, ok := capOf(actual)
	if !ok {
		return false, fmt.Errorf("MatchCap matcher expects a string/array/map/channel/slice, or type with a Cap() int method. Got:\n%s",
			format.Object(actual, 1))
	}
	return matcher.Matcher.Match(capacity)
}

func (matcher *MatchCapMatcher) FailureMessage(actual interface{}) (message string) {
	capacity, _ := capOf(actual)
	return fmt.Sprintf("Expected capacity of\n%s\nto match, but failed with\n%s",
		format.Object(actual, 1),
		format.IndentString(matcher.Matcher.FailureMessage(capacity), 1))
}

func (matcher *MatchCapMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	capacity, _ := capOf(actual)
	return fmt.Sprintf("Expected capacity of\n%s\nnot to match, but did with\n%s",
		format.Object(actual, 1),
		format.IndentString(matcher.Matcher.NegatedFailureMessage(capacity), 1))
}

func capOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	if capper, ok := a.(hasCap); ok {
		return capper.Cap(), true
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Cap(), true
	default:
		return 0, false
	}
}

type hasCap interface {
	Cap() int
}
