package filter

import (
	"fmt"
)

func newMatcherParseError(f string, args ...interface{}) errMatcher {
	return errMatcher{err: fmt.Errorf(f, args...)}
}

type errMatcher struct {
	err error
}

func (e errMatcher) MatcherParseError() bool { return true }
func (e errMatcher) Error() string           { return e.err.Error() }

type matcherParseError interface {
	MatcherParseError() bool
}

func IsMatcherParseError(err error) bool {
	te, ok := err.(matcherParseError)
	return ok && te.MatcherParseError()
}
