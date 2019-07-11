package filter

import (
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

type includeFilter struct {
	matcher Matcher
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
	}
	return true
}

// TODO: right now this just does a copy
func validSelectors(selectors []Selector) []Selector {
	return append([]Selector{}, selectors...)
}
