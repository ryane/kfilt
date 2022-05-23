package filter

import (
	"sort"

	"github.com/ryane/kfilt/pkg/resource"
)

// Filter contains slices of inclusion and exclusion matchers and other configuration
type Filter struct {
	Include []Matcher
	Exclude []Matcher
	count   int
}

// New creates a new Filter
func New() *Filter {
	return &Filter{
		Include: []Matcher{},
		Exclude: []Matcher{},
		count:   0,
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
	for _, matcher := range f.Exclude {
		filtered, err = exclude(filtered, matcher)
		if err != nil {
			return filtered, err
		}
	}

	// includes
	if len(f.Include) > 0 {
		includeMap := make(map[string]interface{})
		included := []resource.Resource{}
		for _, matcher := range f.Include {
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
	filtered, err = count(filtered, f.count)
	if err != nil {
		return filtered, err
	}

	return filtered, nil
}

// AddInclude adds an inclusion matcher
func (f *Filter) AddInclude(s Matcher) {
	f.Include = append(f.Include, s)
}

// AddExclude adds an inclusion matcher
func (f *Filter) AddExclude(s Matcher) {
	f.Exclude = append(f.Exclude, s)
}

func (f *Filter) SetCount(count int) {
	f.count = count
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

func count(resources []resource.Resource, count int) ([]resource.Resource, error) {
	filtered := []resource.Resource{}
	if count > 0 && count <= len(resources) {
		filtered = resources[:count]
	} else {
		filtered = resources
	}

	return filtered, nil
}
