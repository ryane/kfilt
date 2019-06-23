package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
)

func TestNameFilterNil(t *testing.T) {
	matcher := filter.NameMatcher(nil)

	for _, u := range input {
		if !matcher.Match(u) {
			t.Errorf("expected match for %s", u.GetName())
			t.FailNow()
		}
	}
}

func TestNameFilterEmptyVals(t *testing.T) {
	matcher := filter.NameMatcher([]string{" ", ""})

	for _, u := range input {
		if !matcher.Match(u) {
			t.Errorf("expected match for %s", u.GetName())
			t.FailNow()
		}
	}
}

func TestNameFilter(t *testing.T) {
	tests := []struct {
		names       []string
		expectNames []string
	}{
		{[]string{"test-sa"}, []string{"test-sa"}},
		{[]string{"test-sa", "test-sa-2"}, []string{"test-sa", "test-sa-2"}},
		{[]string{"test-deployment"}, []string{"test-deployment"}},
		{[]string{"test-deployment", "test-pod"}, []string{"test-pod", "test-deployment"}},
	}

	for _, test := range tests {
		f := filter.New(filter.NameMatcher(test.names))

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
