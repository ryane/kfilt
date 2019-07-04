package filter

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type nameMatcher struct {
	names []string
}

func NameMatcher(names []string) Matcher {
	return &nameMatcher{validNames(names)}
}

func (f *nameMatcher) Match(u unstructured.Unstructured) bool {
	// no names specified so we just return a match
	if len(f.names) == 0 {
		return true
	}

	for _, name := range f.names {
		if strings.EqualFold(name, u.GetName()) {
			return true
		}
	}

	return false
}

func validNames(names []string) []string {
	valid := []string{}
	for _, name := range names {
		trimName := strings.TrimSpace(name)
		if len(trimName) > 0 {
			valid = append(valid, trimName)
		}
	}
	return valid
}
