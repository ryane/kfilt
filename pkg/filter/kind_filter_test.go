package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
)

func TestKindFilterNil(t *testing.T) {
	matcher := filter.KindMatcher(nil)

	for _, u := range input {
		if !matcher.Match(u) {
			t.Errorf("expected match for %s", u.GetKind())
			t.FailNow()
		}
	}
}

func TestKindFilterEmptyVals(t *testing.T) {
	matcher := filter.KindMatcher([]string{" ", ""})

	for _, u := range input {
		if !matcher.Match(u) {
			t.Errorf("expected match for %s", u.GetKind())
			t.FailNow()
		}
	}
}

func TestKindFilter(t *testing.T) {
	tests := []struct {
		kinds       []string
		expectNames []string
	}{
		{[]string{"ServiceAccount"}, []string{"test-sa", "test-sa-2"}},
		{[]string{"Deployment"}, []string{"test-deployment"}},
		{[]string{"Pod"}, []string{"test-pod"}},
		{[]string{"ServiceAccount", "Deployment"}, []string{"test-sa", "test-sa-2", "test-deployment"}},
	}

	for _, test := range tests {
		f := filter.New(filter.KindMatcher(test.kinds))

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
