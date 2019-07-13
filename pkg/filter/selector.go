package filter

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Selector struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string
}

func (s *Selector) Match(u unstructured.Unstructured) bool {
	gvk := u.GroupVersionKind()

	if s.Group != "" && !strings.EqualFold(s.Group, gvk.Group) {
		return false
	}
	if s.Version != "" && !strings.EqualFold(s.Version, gvk.Version) {
		return false
	}
	if s.Kind != "" && !strings.EqualFold(s.Kind, gvk.Kind) {
		return false
	}
	if s.Name != "" && !strings.EqualFold(s.Name, u.GetName()) {
		return false
	}
	if s.Namespace != "" {
		ns := u.GetNamespace()
		if strings.ToLower(s.Namespace) == "default" {
			if ns != "" && !strings.EqualFold(s.Namespace, ns) {
				return false
			}
		} else if !strings.EqualFold(s.Namespace, ns) {
			return false
		}
	}

	return true
}

func NewSelector(q string) (Selector, error) {
	m := Selector{}
	criteria := strings.Split(q, ",")

	if len(criteria) == 0 {
		return m, newMatcherError("invalid matcher %q. query is required", q)
	}

	for _, criterion := range criteria {
		parts := strings.Split(criterion, "=")
		if len(parts) != 2 {
			return m, newMatcherError("invalid matcher %q. Should be in the format %q", criterion, "key=value")
		}

		key, val := strings.ToLower(parts[0]), parts[1]

		switch key {
		case "kind", "k":
			m.Kind = val
		case "name", "n":
			m.Name = val
		default:
			return m, newMatcherError("invalid matcher %q. key should be one of %q", criterion, validMatcherKeys())
		}
	}

	return m, nil
}

func newMatcherError(f string, args ...interface{}) errMatcherError {
	return errMatcherError{err: fmt.Errorf(f, args...)}
}

type errMatcherError struct {
	err error
}

func (e errMatcherError) MatcherError() bool { return true }
func (e errMatcherError) Error() string      { return e.err.Error() }

type matcherError interface {
	MatcherError() bool
}

func IsMatcherError(err error) bool {
	te, ok := err.(matcherError)
	return ok && te.MatcherError()
}

func validMatcherKeys() []string {
	return []string{"kind", "k", "name", "n"}
}
