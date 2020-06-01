package matchers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type AtKeyMatcher struct {
	Key      interface{}
	Matcher  types.GomegaMatcher
	notFound bool
	obj      interface{}
}

func (m *AtKeyMatcher) Match(actual interface{}) (success bool, err error) {
	v := reflect.ValueOf(actual)
	if v.Kind() != reflect.Map {
		return false, errors.New("AtKey matcher requires an actual of map")
	}
	found := v.MapIndex(reflect.ValueOf(m.Key))
	if !found.IsValid() {
		m.notFound = true
		return false, nil
	}
	m.obj = found.Interface()
	return m.Matcher.Match(m.obj)
}

func (m *AtKeyMatcher) FailureMessage(actual interface{}) (message string) {
	if m.notFound {
		return fmt.Sprintf("Map\n%s\ndoes not contain key\n%s",
			format.Object(actual, 1),
			format.Object(m.Key, 1))
	}
	return fmt.Sprintf("Matcher failed at map key %s. %s",
		format.Object(m.Key, 0),
		m.Matcher.FailureMessage(m.obj))
}

func (m *AtKeyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.obj)
}
