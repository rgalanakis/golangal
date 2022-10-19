package errmatch

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/rgalanakis/golangal/errmatch/matchers"
)

// BeCausedBy is a robust error matcher that works with a number of error libraries,
// including Go 1.13 errors, pkg/errors, and hashicorp/multierror.
// The expected value is used to construct a gomega.MatchError matcher,
// and then it is used to check one of the following:
//
//   - actual passes MatchError(expected)
//   - actual can be unwrapped, and the unwrapped value passes BeCausedBy.
//   - pkg/errors.Cause(actual) passes MatchError(expected).
//     This would be hit when you use pkg/errors.Wrap.
//     Note that this only looks for the 'cause' error,
//     and would not work to match error messages of intermediate wrapped errors.
//     However, you can use BeCausedBy("outer: inner") (or MatchError) to match the full message.
//   - any of ((*multierror.Error)actual).Errors passes BeCausedBy.
//     That is, if you use errors.Wrap in conjunction with multierror,
//     the individual errors in multierror are unwrapped with errors.Cause.
//
// If all of those do not pass, BeCausedBy fails.
func BeCausedBy(expected interface{}) types.GomegaMatcher {
	return &matchers.BeCausedByMatcher{MatchError: gomega.MatchError(expected)}
}
