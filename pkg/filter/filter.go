package filter

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Filter interface {
	Filter([]unstructured.Unstructured) []unstructured.Unstructured
}

type defaultFilter struct {
	matchers []Matcher
}

type Matcher interface {
	Match(unstructured.Unstructured) bool
}

func New(matchers ...Matcher) Filter {
	return &defaultFilter{matchers}
}

func (f *defaultFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	filtered := unstructureds

	for _, matcher := range f.matchers {
		filtered = filter(filtered, matcher)
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
