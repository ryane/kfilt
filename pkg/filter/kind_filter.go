package filter

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type kindFilter struct {
	matcher Matcher
}

func (f *kindFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return filter(unstructureds, f.matcher)
	}
	return unstructureds
}

func KindFilter(kinds ...string) Filter {
	return &kindFilter{KindMatcher(kinds)}
}

type excludeKindFilter struct {
	matcher Matcher
}

func (f *excludeKindFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return excludeFilter(unstructureds, f.matcher)
	}
	return unstructureds
}

func ExcludeKindFilter(kinds ...string) Filter {
	return &excludeKindFilter{KindMatcher(kinds)}
}

type kindMatcher struct {
	kinds []string
}

func KindMatcher(kinds []string) Matcher {
	return &kindMatcher{validKinds(kinds)}
}

func (f *kindMatcher) Valid() bool {
	return len(f.kinds) > 0
}

func (f *kindMatcher) Match(u unstructured.Unstructured) bool {
	for _, kind := range f.kinds {
		if strings.EqualFold(kind, u.GetKind()) {
			return true
		}
	}

	return false
}

func validKinds(kinds []string) []string {
	valid := []string{}
	for _, kind := range kinds {
		trimKind := strings.TrimSpace(kind)
		if len(trimKind) > 0 {
			valid = append(valid, trimKind)
		}
	}
	return valid
}
