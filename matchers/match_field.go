package matchers

import (
	"reflect"

	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type MatchFieldMatcher struct {
	Name    string
	Matcher types.GomegaMatcher
}

func (m *MatchFieldMatcher) Match(actual interface{}) (success bool, err error) {
	actualVal := reflect.ValueOf(actual)
	if actualVal.Kind() != reflect.Struct {
		err := fmt.Errorf("MatchField matcher requires an actual of kind struct, not %s",
			actualVal.Kind().String())
		return false, err
	}
	fieldVal := actualVal.FieldByName(m.Name)
	if !fieldVal.IsValid() {
		err := fmt.Errorf("field '%s' does not exist on type %s", m.Name, actualVal.Type().Name())
		return false, err
	}

	return m.Matcher.Match(fieldVal.Interface())
}

func (m *MatchFieldMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Field %s of\n%s\ndid not match. %s",
		m.Name, format.Object(actual, 1), m.Matcher.FailureMessage(m.fieldValue(actual)))
}

func (m *MatchFieldMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Field %s of\n%s\nmatched. %s",
		m.Name, format.Object(actual, 1), m.Matcher.FailureMessage(m.fieldValue(actual)))
}

func (m *MatchFieldMatcher) fieldValue(actual interface{}) interface{} {
	return reflect.ValueOf(actual).FieldByName(m.Name).Interface()
}
