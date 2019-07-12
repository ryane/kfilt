package filter

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Filter struct {
	Include []Selector
	Exclude []Selector
}

func New() *Filter {
	return &Filter{
		Include: []Selector{},
		Exclude: []Selector{},
	}
}

func (f *Filter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	filtered := append([]unstructured.Unstructured{}, unstructureds...)

	// excludes
	for _, selector := range f.Exclude {
		filtered = exclude(filtered, selector)
	}

	// includes
	if len(f.Include) > 0 {
		included := []unstructured.Unstructured{}
		for _, selector := range f.Include {
			included = append(included, filter(filtered, selector)...)
		}
		filtered = included
	}

	return filtered
}

func (f *Filter) AddInclude(s Selector) {
	f.Include = append(f.Include, s)
}

func (f *Filter) AddExclude(s Selector) {
	f.Exclude = append(f.Exclude, s)
}

func filter(unstructureds []unstructured.Unstructured, matcher Selector) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, u := range unstructureds {
		if matcher.Match(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}

func exclude(unstructureds []unstructured.Unstructured, matcher Selector) []unstructured.Unstructured {
	filtered := []unstructured.Unstructured{}
	for _, u := range unstructureds {
		if !matcher.Match(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}
