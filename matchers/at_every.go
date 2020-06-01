package matchers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"
)

type AtEveryMatcher struct {
	Matcher types.GomegaMatcher

	length int
	index  int
	obj    interface{}
}

func (m *AtEveryMatcher) Match(actual interface{}) (success bool, err error) {
	actualVal := reflect.ValueOf(actual)
	if actualVal.Kind() != reflect.Slice {
		return false, errors.New("AtEvery matcher requires an actual of type slice")
	}
	m.length = actualVal.Len()
	if m.length == 0 {
		return false, nil
	}
	for i := 0; i < m.length; i++ {
		m.index = i
		m.obj = actualVal.Index(i).Interface()
		if success, err := m.Matcher.Match(m.obj); err != nil {
			return false, err
		} else if !success {
			return false, nil
		}
	}
	return true, nil
}

func (m *AtEveryMatcher) FailureMessage(actual interface{}) (message string) {
	if m.length == 0 {
		return "Did not expect empty collection"
	}
	return fmt.Sprintf("Match failed at index %d:\n%s",
		m.index, m.Matcher.FailureMessage(m.obj))
}

func (m *AtEveryMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.obj)
}
