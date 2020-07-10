package matchers

import (
	"github.com/hashicorp/go-multierror"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
)

type BeCausedByMatcher struct {
	MatchError types.GomegaMatcher
}

func (m *BeCausedByMatcher) Match(actual interface{}) (success bool, err error) {
	success, err = m.MatchError.Match(actual)
	if success || err != nil {
		return
	}
	// We now need to use recursive matching for everything.
	// Eventually we'll unwind and hit the normal MatchError.Match call above.
	actualErr := actual.(error)
	if unwrapped := unwrap(actualErr); unwrapped != actualErr && m.matchRecursive(unwrapped) {
		return true, nil
	}
	if cause := errors.Cause(actualErr); cause != actualErr && m.matchRecursive(cause) {
		return true, nil
	}
	if me, ok := actualErr.(*multierror.Error); ok {
		for _, e := range me.Errors {
			if m.matchRecursive(e) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *BeCausedByMatcher) matchRecursive(err error) bool {
	if err == nil {
		return false
	}
	m2 := &BeCausedByMatcher{MatchError: m.MatchError}
	success, _ := m2.Match(err)
	return success
}

func (m *BeCausedByMatcher) FailureMessage(actual interface{}) (message string) {
	return m.MatchError.FailureMessage(actual)
}

func (m *BeCausedByMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	panic("bad idea to use a matcher like this negated, it's too specific")
}

// Copied from Go 1.13 errors.Unwrap
func unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
