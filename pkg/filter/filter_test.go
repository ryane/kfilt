package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		names       []string
		kinds       []string
		expectNames []string
	}{
		{[]string{"test-sa"}, []string{"Deployment"}, []string{}},
		{[]string{"test-sa"}, []string{}, []string{"test-sa"}},
		{[]string{"test-sa"}, []string{""}, []string{"test-sa"}},
		{[]string{}, []string{"ServiceAccount"}, []string{"test-sa", "test-sa-2"}},
		{[]string{"test-pod", "test-deployment"}, []string{"ServiceAccount"}, []string{}},
	}

	for _, test := range tests {
		f := filter.New(
			filter.KindMatcher(test.kinds),
			filter.NameMatcher(test.names),
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
