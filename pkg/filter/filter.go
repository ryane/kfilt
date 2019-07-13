package filter

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Filter struct {
	Include []Matcher
	Exclude []Matcher
}

func New() *Filter {
	return &Filter{
		Include: []Matcher{},
		Exclude: []Matcher{},
	}
}

func (f *Filter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	filtered := append([]unstructured.Unstructured{}, unstructureds...)

	// excludes
	for _, matcher := range f.Exclude {
		filtered = exclude(filtered, matcher)
	}

	// includes
	if len(f.Include) > 0 {
		included := []unstructured.Unstructured{}
		for _, matcher := range f.Include {
			included = append(included, filter(filtered, matcher)...)
		}
		filtered = included
	}

	return filtered
}

func (f *Filter) AddInclude(s Matcher) {
	f.Include = append(f.Include, s)
}

func (f *Filter) AddExclude(s Matcher) {
	f.Exclude = append(f.Exclude, s)
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
