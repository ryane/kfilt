package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
)

type includeNames []string
type includeKinds []string
type excludeNames []string
type excludeKinds []string
type expectNames []string

func TestFilter(t *testing.T) {
	tests := []struct {
		includeKinds includeKinds
		includeNames includeNames
		excludeKinds excludeKinds
		excludeNames excludeNames
		expectNames  []string
	}{
		{
			includeKinds{"Deployment", "Pod"},
			includeNames{},
			excludeKinds{},
			excludeNames{},
			expectNames{"test-pod", "test-deployment"},
		},
		{
			includeKinds{"Deployment"},
			includeNames{"test-sa"},
			excludeKinds{},
			excludeNames{},
			expectNames{},
		},
		{
			includeKinds{},
			includeNames{"test-sa"},
			excludeKinds{},
			excludeNames{},
			expectNames{"test-sa"},
		},
		{
			includeKinds{},
			includeNames{"test-sa", "test-sa-2"},
			excludeKinds{},
			excludeNames{},
			expectNames{"test-sa", "test-sa-2"},
		},
		{
			includeKinds{""},
			includeNames{"test-sa"},
			excludeKinds{},
			excludeNames{},
			expectNames{"test-sa"},
		},
		{
			includeKinds{"ServiceAccount"},
			includeNames{},
			excludeKinds{},
			excludeNames{},
			expectNames{"test-sa", "test-sa-2"},
		},
		{
			includeKinds{"ServiceAccount"},
			includeNames{"test-pod", "test-deployment"},
			excludeKinds{},
			excludeNames{},
			expectNames{},
		},
		{
			includeKinds{"ServiceAccount"},
			includeNames{},
			excludeKinds{},
			excludeNames{"test-sa"},
			expectNames{"test-sa-2"},
		},
		{
			includeKinds{},
			includeNames{},
			excludeKinds{"ServiceAccount"},
			excludeNames{},
			expectNames{"test-pod", "test-deployment"},
		},
		{
			includeKinds{},
			includeNames{},
			excludeKinds{"ServiceAccount", "Deployment"},
			excludeNames{},
			expectNames{"test-pod"},
		},
		{
			includeKinds{"ServiceAccount", "Deployment"},
			includeNames{},
			excludeKinds{"ServiceAccount"},
			excludeNames{},
			expectNames{"test-deployment"},
		},
		{
			includeKinds{},
			includeNames{"test-sa", "test-sa-2"},
			excludeKinds{"ServiceAccount"},
			excludeNames{},
			expectNames{},
		},
	}

	for _, test := range tests {
		f := filter.New(
			filter.ExcludeNameFilter(test.excludeNames...),
			filter.ExcludeKindFilter(test.excludeKinds...),
			filter.NameFilter(test.includeNames...),
			filter.KindFilter(test.includeKinds...),
		)

		results := f.Filter(input)
		if len(results) != len(test.expectNames) {
			t.Errorf("expected %d results, got %d", len(test.expectNames), len(results))
			t.FailNow()
		}

		for i, u := range results {
			name := u.GetName()
			if name != test.expectNames[i] {
				t.Errorf("expected %s, got %s", test.expectNames[i], name)
				t.FailNow()
			}
		}
	}
}
