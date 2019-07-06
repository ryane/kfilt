package filter

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type nameFilter struct {
	matcher Matcher
}

func (f *nameFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return filter(unstructureds, f.matcher)
	}
	return unstructureds
}

func NameFilter(names ...string) Filter {
	return &nameFilter{NameMatcher(names)}
}

type excludeNameFilter struct {
	matcher Matcher
}

func (f *excludeNameFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return excludeFilter(unstructureds, f.matcher)
	}
	return unstructureds
}

func ExcludeNameFilter(names ...string) Filter {
	return &excludeNameFilter{NameMatcher(names)}
}

type nameMatcher struct {
	names []string
}

func NameMatcher(names []string) Matcher {
	return &nameMatcher{validNames(names)}
}

func (f *nameMatcher) Valid() bool {
	return len(f.names) > 0
}

func (f *nameMatcher) Match(u unstructured.Unstructured) bool {
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
