package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type MatchLenMatcher struct {
	Matcher types.GomegaMatcher
}

func (matcher *MatchLenMatcher) Match(actual interface{}) (bool, error) {
	length, ok := lengthOf(actual)
	if !ok {
		return false, fmt.Errorf("MatchLen matcher expects a string/array/map/channel/slice, or type with a Len() int method. Got:\n%s",
			format.Object(actual, 1))
	}
	return matcher.Matcher.Match(length)
}

func (matcher *MatchLenMatcher) FailureMessage(actual interface{}) (message string) {
	length, _ := lengthOf(actual)
	return fmt.Sprintf("Expected length of\n%s\nto match, but failed with\n%s",
		format.Object(actual, 1),
		format.IndentString(matcher.Matcher.FailureMessage(length), 1))
}

func (matcher *MatchLenMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	length, _ := lengthOf(actual)
	return fmt.Sprintf("Expected length of\n%s\nnot to match, but did with\n%s",
		format.Object(actual, 1),
		format.IndentString(matcher.Matcher.NegatedFailureMessage(length), 1))
}

func lengthOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	if lengther, ok := a.(hasLen); ok {
		return lengther.Len(), true
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	default:
		return 0, false
	}
}

type hasLen interface {
	Len() int
}
