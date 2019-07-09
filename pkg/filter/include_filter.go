package filter

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type includeFilter struct {
	matcher Matcher
}

type Selector struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string
}

func (f *includeFilter) Filter(unstructureds []unstructured.Unstructured) []unstructured.Unstructured {
	if f.matcher.Valid() {
		return filter(unstructureds, f.matcher)
	}
	return unstructureds
}

func IncludeFilter(selectors ...Selector) Filter {
	return &includeFilter{SelectorMatcher(selectors)}
}

type selectorMatcher struct {
	selectors []Selector
}

func SelectorMatcher(selectors []Selector) Matcher {
	return &selectorMatcher{validSelectors(selectors)}
}

func (f *selectorMatcher) Valid() bool {
	return len(f.selectors) > 0
}

func (f *selectorMatcher) Match(u unstructured.Unstructured) bool {
	for _, s := range f.selectors {
		gvk := u.GroupVersionKind()
		if s.Group != "" && s.Group != gvk.Group {
			return false
		}
		if s.Version != "" && s.Version != gvk.Version {
			return false
		}
		if s.Kind != "" && s.Kind != gvk.Kind {
			return false
		}
		if s.Name != "" && s.Name != u.GetName() {
			return false
		}
		if s.Namespace != "" {
			ns := u.GetNamespace()
			if s.Namespace == "default" {
				if ns != "" && s.Namespace != ns {
					return false
				}
			} else if s.Namespace != ns {
				return false
			}
		}
	}
	return true
}

// TODO: right now this just does a copy
func validSelectors(selectors []Selector) []Selector {
	return append([]Selector{}, selectors...)
}
