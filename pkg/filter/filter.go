package filter

import (
	"github.com/ryane/kfilt/pkg/resource"
)

// Filter contains slices of inclusion and exclusion matchers
type Filter struct {
	Include []Matcher
	Exclude []Matcher
}

// New creates a new Filter
func New() *Filter {
	return &Filter{
		Include: []Matcher{},
		Exclude: []Matcher{},
	}
}

// Filter returns a filtered slice of Resources
func (f *Filter) Filter(resources []resource.Resource) []resource.Resource {
	filtered := append([]resource.Resource{}, resources...)

	// excludes
	for _, matcher := range f.Exclude {
		filtered = exclude(filtered, matcher)
	}

	// includes
	if len(f.Include) > 0 {
		includeMap := make(map[string]interface{})
		included := []resource.Resource{}
		for _, matcher := range f.Include {
			for _, match := range filter(filtered, matcher) {
				matchID := match.ID()
				if _, ok := includeMap[matchID]; !ok {
					includeMap[matchID] = struct{}{}
					included = append(included, match)
				}
			}
		}
		filtered = included
	}

	return filtered
}

// AddInclude adds an inclusion matcher
func (f *Filter) AddInclude(s Matcher) {
	f.Include = append(f.Include, s)
}

// AddExclude adds an inclusion matcher
func (f *Filter) AddExclude(s Matcher) {
	f.Exclude = append(f.Exclude, s)
}

func filter(resources []resource.Resource, matcher Matcher) []resource.Resource {
	filtered := []resource.Resource{}
	for _, r := range resources {
		if matcher.Match(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func exclude(resources []resource.Resource, matcher Matcher) []resource.Resource {
	filtered := []resource.Resource{}
	for _, r := range resources {
		if !matcher.Match(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
