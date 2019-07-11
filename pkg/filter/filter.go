package filter

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Filter interface {
	Filter([]unstructured.Unstructured) []unstructured.Unstructured
}

type chainedFilter struct {
	filters []Filter
}

type Matcher interface {
	Match(unstructured.Unstructured) bool
	Valid() bool
}

func New(filters ...Filter) Filter {
	return &chainedFilter{filters}
}

func (f *chainedFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	filtered := unstructureds

	for _, filter := range f.filters {
		filtered = filter.Filter(filtered)
	}

	return filtered
}

func filter(unstructureds []unstructured.Unstructured, matcher Matcher) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, u := range unstructureds {
		if matcher.Match(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func exclude(unstructureds []unstructured.Unstructured, matcher Matcher) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, u := range unstructureds {
		if !matcher.Match(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}
