package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type PanicWith struct {
	PanicMatcher types.GomegaMatcher
	panickedWith interface{}
}

func (matcher *PanicWith) Match(actual interface{}) (succeeded bool, err error) {
	if actual == nil {
		return false, fmt.Errorf("PanicWith expects a non-nil actual")
	}

	actualType := reflect.TypeOf(actual)
	if actualType.Kind() != reflect.Func {
		err := fmt.Errorf("PanicWith expects a function. Got:\n%s", format.Object(actual, 1))
		return false, err
	}
	if !(actualType.NumIn() == 0 && actualType.NumOut() == 0) {
		err := fmt.Errorf("PanicWith expects a function with no arguments and no return value. Got:\n%s",
			format.Object(actual, 1))
		return false, err
	}

	defer func() {
		if e := recover(); e != nil {
			matcher.panickedWith = e
			succeeded, err = matcher.PanicMatcher.Match(e)
		}
	}()
	reflect.ValueOf(actual).Call([]reflect.Value{})
	return
}

func (matcher *PanicWith) FailureMessage(actual interface{}) (message string) {
	if matcher.panickedWith == nil {
		return format.Message(actual, "to panic, but did not")
	}
	return format.Message(
		actual,
		fmt.Sprintf("to panic with matcher, but matcher failed:\n%s",
			matcher.PanicMatcher.FailureMessage(matcher.panickedWith)))
}

func (matcher *PanicWith) NegatedFailureMessage(actual interface{}) (message string) {
	panic("Do not use negated PanicWith, it leads to ambiguous test conditions")
}
