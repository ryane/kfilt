package filter

import (
	"sort"

	"github.com/ryane/kfilt/pkg/resource"
)

// Filter contains slices of inclusion and exclusion matchers and other
// configuration
type Filter struct {
	include []Matcher
	exclude []Matcher
	limit   int
}

// New creates a new Filter
func New() *Filter {
	return &Filter{
		include: []Matcher{},
		exclude: []Matcher{},
		limit:   0,
	}
}

// Filter returns a filtered slice of Resources
func (f *Filter) Filter(resources []resource.Resource) ([]resource.Resource, error) {
	var err error
	filtered := append([]resource.Resource{}, resources...)

	ordermap := map[string]int{}
	for i, res := range resources {
		ordermap[res.ID()] = i
	}

	// excludes
	for _, matcher := range f.exclude {
		filtered, err = exclude(filtered, matcher)
		if err != nil {
			return filtered, err
		}
	}

	// includes
	if len(f.include) > 0 {
		includeMap := make(map[string]interface{})
		included := []resource.Resource{}
		for _, matcher := range f.include {
			results, err := filter(filtered, matcher)
			if err != nil {
				return filtered, err
			}
			for _, match := range results {
				matchID := match.ID()
				if _, ok := includeMap[matchID]; !ok {
					includeMap[matchID] = struct{}{}
					included = append(included, match)
				}
			}
		}
		filtered = included
	}

	sort.Slice(filtered, func(i, j int) bool {
		return ordermap[filtered[i].ID()] < ordermap[filtered[j].ID()]
	})

	// limit output
	filtered, err = limit(filtered, f.limit)
	if err != nil {
		return filtered, err
	}

	return filtered, nil
}

func (f *Filter) Included() []Matcher {
	return f.include
}

func (f *Filter) Excluded() []Matcher {
	return f.exclude
}

// AddInclude adds an inclusion matcher
func (f *Filter) AddInclude(s Matcher) {
	f.include = append(f.include, s)
}

// AddExclude adds an inclusion matcher
func (f *Filter) AddExclude(s Matcher) {
	f.exclude = append(f.exclude, s)
}

func (f *Filter) Limit(limit int) {
	f.limit = limit
}

func filter(resources []resource.Resource, matcher Matcher) ([]resource.Resource, error) {
	filtered := []resource.Resource{}
	for _, r := range resources {
		ok, err := matcher.Match(r)
		if err != nil {
			return filtered, err
		}
		if ok {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func exclude(resources []resource.Resource, matcher Matcher) ([]resource.Resource, error) {
	filtered := []resource.Resource{}
	for _, r := range resources {

		ok, err := matcher.Match(r)
		if err != nil {
			return filtered, err
		}
		if !ok {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

func limit(resources []resource.Resource, limit int) ([]resource.Resource, error) {
	var filtered []resource.Resource
	if limit > 0 && limit <= len(resources) {
		filtered = resources[:limit]
	} else {
		filtered = resources
	}

	return filtered, nil
}
