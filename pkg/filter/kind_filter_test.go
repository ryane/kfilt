package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestKindMatcher(t *testing.T) {
	var u = unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind": "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa",
			},
		},
	}

	tests := []struct {
		kinds    []string
		expected bool
	}{
		{[]string{"ServiceAccount"}, true},
		{[]string{"serviceaccount"}, true},
		{[]string{""}, false},
		{[]string{" "}, false},
		{[]string{}, false},
		{[]string{"pod"}, false},
		{[]string{"Deployment", "ServiceAccount"}, true},
	}

	for _, test := range tests {
		matcher := filter.KindMatcher(test.kinds)
		if result := matcher.Match(u); result != test.expected {
			t.Errorf("expected %v for %v, got %v", test.expected, test.kinds, result)
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
		{[]string{"ServiceAccount", "Deployment"}, []string{"test-sa", "test-sa-2", "test-deployment"}},
		{[]string{"Deployment"}, []string{"test-deployment"}},
		{[]string{"deployment"}, []string{"test-deployment"}},
		{[]string{"Pod"}, []string{"test-pod"}},
		{[]string{"ServiceAccount", "Deployment"}, []string{"test-sa", "test-sa-2", "test-deployment"}},
	}

	for _, test := range tests {
		f := filter.KindFilter(test.kinds...)

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

func TestExcludeKindFilter(t *testing.T) {
	tests := []struct {
		kinds       []string
		expectNames []string
	}{
		{[]string{}, []string{"test-sa", "test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"", " "}, []string{"test-sa", "test-sa-2", "test-pod", "test-deployment"}},
		{[]string{"ServiceAccount"}, []string{"test-pod", "test-deployment"}},
		{[]string{"Deployment"}, []string{"test-sa", "test-sa-2", "test-pod"}},
		{[]string{"deployment"}, []string{"test-sa", "test-sa-2", "test-pod"}},
		{[]string{"Pod"}, []string{"test-sa", "test-sa-2", "test-deployment"}},
		{[]string{"ServiceAccount", "Deployment"}, []string{"test-pod"}},
	}

	for _, test := range tests {
		f := filter.ExcludeKindFilter(test.kinds...)

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
