package matchers_test

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"
	. "github.com/rgalanakis/golangal/errmatch"
)

var _ = Describe("BeCausedBy matcher", func() {
	e := errors.New("blah")

	It("matches if the error passes MatchError", func() {
		Expect(e).To(BeCausedBy(e))
		Expect(e).To(BeCausedBy("blah"))
	})

	It("matches if the error is wrapped by pkgerrors.Wrap", func() {
		Expect(pkgerrors.Wrap(e, "foo")).To(BeCausedBy("blah"))
		Expect(pkgerrors.Wrap(e, "foo")).To(BeCausedBy("foo: blah"))
		Expect(pkgerrors.Wrap(e, "foo")).To(BeCausedBy(e))
	})

	It("matches if the error is in a multierror", func() {
		Expect(multierror.Append(e)).To(BeCausedBy("blah"))
		Expect(multierror.Append(e)).To(BeCausedBy(e))
		Expect(multierror.Append(pkgerrors.New("foo"), e)).To(BeCausedBy(e))
	})

	It("matches if the error is wrapped and a member of a multierror", func() {
		Expect(multierror.Append(
			pkgerrors.New("foo"),
			pkgerrors.Wrap(e, "spam"),
		)).To(BeCausedBy(e))
	})

	It("fails if actual does not match", func() {
		matcher := BeCausedBy("foo")
		success, err := matcher.Match(e)
		Expect(success).To(BeFalse())
		Expect(err).ToNot(HaveOccurred())
		msg := matcher.FailureMessage(e)
		Expect(msg).To(HavePrefix(`Expected
    <*errors.errorString | `))
		Expect(msg).To(HaveSuffix(`>: {s: "blah"}
to match error
    <string>: foo`))
	})

	It("errors if actual is nil", func() {
		matcher := BeCausedBy(e)
		success, err := matcher.Match(nil)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("Expected an error, got nil"))
	})

	It("errors if actual is not an error", func() {
		matcher := BeCausedBy("abc")
		success, err := matcher.Match(5)
		Expect(success).To(BeFalse())
		Expect(err).To(MatchError("Expected an error.  Got:\n    <int>: 5"))
	})
})
