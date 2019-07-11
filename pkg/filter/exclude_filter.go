package filter

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type excludeFilter struct {
	matcher Matcher
}

func (f *excludeFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return filter(unstructureds, f.matcher)
	}
	return unstructureds
}

func ExcludeFilter(selectors ...Selector) Filter {
	return &excludeFilter{SelectorMatcher(selectors)}
}
