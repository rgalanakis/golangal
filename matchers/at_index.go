package matchers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type AtIndexMatcher struct {
	Index   int
	Matcher types.GomegaMatcher

	badLen     bool
	failingIdx int
	obj        interface{}
}

func (m *AtIndexMatcher) Match(actual interface{}) (success bool, err error) {
	actualVal := reflect.ValueOf(actual)
	k := actualVal.Kind()
	if k != reflect.String && k != reflect.Slice && k != reflect.Array {
		return false, errors.New("AtIndex matcher requires an actual of string, slice, or array")
	}
	if m.Index >= actualVal.Len() {
		m.badLen = true
		return false, nil
	}
	m.obj = actualVal.Index(m.Index).Interface()
	m.failingIdx = m.Index
	return m.Matcher.Match(m.obj)
}

func (m *AtIndexMatcher) FailureMessage(actual interface{}) (message string) {
	if m.badLen {
		return fmt.Sprintf("Slice\n%s\nis too short to match against index %d",
			format.Object(actual, 1), m.Index)
	}
	return fmt.Sprintf("Matcher failed at slice index %d. %s", m.failingIdx, m.Matcher.FailureMessage(m.obj))
}

func (m *AtIndexMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.obj)
}
