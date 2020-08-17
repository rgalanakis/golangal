package matchers

import "github.com/pkg/errors"

type NotErrorMatcher struct {
}

func (m *NotErrorMatcher) Match(actual interface{}) (success bool, err error) {
	if e, ok := actual.(error); ok {
		return false, errors.Errorf(
			"MatchError cannot be used with functions that return an error. Got: %s", e)
	}
	return true, nil
}

func (m *NotErrorMatcher) FailureMessage(actual interface{}) (message string) {
	panic("cannot fail")
}

func (m *NotErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic("cannot fail")
}
