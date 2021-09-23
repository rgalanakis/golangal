package matchers

import (
	"reflect"

	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type MatchPtrFieldMatcher struct {
	Name    string
	Matcher types.GomegaMatcher
}

func (m *MatchPtrFieldMatcher) Match(actual interface{}) (success bool, err error) {
	actualPtrVal := reflect.ValueOf(actual)
	if actualPtrVal.Kind() != reflect.Ptr {
		err := fmt.Errorf("MatchPtrField matcher requires an actual of kind ptr, not %s",
			actualPtrVal.Kind().String())
		return false, err
	}
	actualVal := actualPtrVal.Elem()
	fieldVal := actualVal.FieldByName(m.Name)
	if !fieldVal.IsValid() {
		err := fmt.Errorf("field '%s' does not exist on type %s", m.Name, actualVal.Type().Name())
		return false, err
	}

	return m.Matcher.Match(fieldVal.Interface())
}

func (m *MatchPtrFieldMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Field %s of\n%s\ndid not match. %s",
		m.Name, format.Object(actual, 1), m.Matcher.FailureMessage(m.fieldValue(actual)))
}

func (m *MatchPtrFieldMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Field %s of\n%s\nmatched. %s",
		m.Name, format.Object(actual, 1), m.Matcher.FailureMessage(m.fieldValue(actual)))
}

func (m *MatchPtrFieldMatcher) fieldValue(actual interface{}) interface{} {
	return reflect.ValueOf(actual).Elem().FieldByName(m.Name).Interface()
}
