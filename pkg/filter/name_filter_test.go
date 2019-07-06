package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNameMatcher(t *testing.T) {
	var u = unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind": "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa",
			},
		},
	}

	tests := []struct {
		names    []string
		expected bool
	}{
		{[]string{"test-sa"}, true},
		{[]string{"TEST-SA"}, true},
		{[]string{""}, false},
		{[]string{" "}, false},
		{[]string{}, false},
		{[]string{"test-pod"}, false},
		{[]string{"test-deployment", "test-sa"}, true},
	}

	for _, test := range tests {
		matcher := filter.NameMatcher(test.names)
		if result := matcher.Match(u); result != test.expected {
			t.Errorf("expected %v for %v, got %v", test.expected, test.names, result)
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
		{[]string{"TEST-DEPLOYMENT"}, []string{"test-deployment"}},
		{[]string{"test-deployment", "test-pod"}, []string{"test-pod", "test-deployment"}},
	}

	for _, test := range tests {
		f := filter.NameFilter(test.names...)

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

func TestExcludeNameFilter(t *testing.T) {
	tests := []struct {
		names       []string
		expectNames []string
	}{
		{[]string{}, []string{"test-sa", "test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"", " "}, []string{"test-sa", "test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"test-sa"}, []string{"test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"TEST-sa"}, []string{"test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"test-deployment"}, []string{"test-sa", "test-sa-2", "test-pod"}},
		{[]string{"test-pod"}, []string{"test-sa", "test-sa-2", "test-deployment"}},
		{[]string{"test-sa", "test-deployment", "test-sa-2"}, []string{"test-pod"}},
	}

	for _, test := range tests {
		f := filter.ExcludeNameFilter(test.names...)

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
