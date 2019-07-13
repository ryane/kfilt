package filter_test

import (
	"errors"
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
)

func TestIsMatcherParseError(t *testing.T) {
	e := errors.New("test error")

	if filter.IsMatcherParseError(e) {
		t.Errorf("%q is not a matcher parse error", e.Error())
	}

	pe := fakeMatcherParserError{e}
	if !filter.IsMatcherParseError(pe) {
		t.Errorf("%q is a matcher parse error", e.Error())
	}
}

type fakeMatcherParserError struct{ err error }

func (e fakeMatcherParserError) MatcherParseError() bool { return true }
func (e fakeMatcherParserError) Error() string           { return e.err.Error() }
