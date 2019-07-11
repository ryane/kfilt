package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSelectorMatcher(t *testing.T) {
	tests := []struct {
		selectors []filter.Selector
		resource  unstructured.Unstructured
		expected  bool
	}{
		// empty selector should match
		{[]filter.Selector{}, role(), true},
		// kind selectors
		{
			[]filter.Selector{
				{Kind: "Role"},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{Kind: "role"},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{Kind: "ServiceAccount"},
			},
			role(),
			false,
		},
		{
			[]filter.Selector{
				{Version: "v1", Kind: "ServiceAccount"},
			},
			serviceAccount(),
			true,
		},
		{
			[]filter.Selector{
				{Group: "", Version: "v1", Kind: "ServiceAccount"},
			},
			serviceAccount(),
			true,
		},
		{
			[]filter.Selector{
				{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Kind: "Role"},
			},
			role(),
			false,
		},
		{
			[]filter.Selector{
				{
					Group:   "rbac.authorization.k8s.io",
					Version: "v1",
					Kind:    "Role",
					Name:    "test-role",
				},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{
					Group:   "rbac.authorization.k8s.io",
					Version: "v1",
					Kind:    "Role",
					Name:    "test-sa",
				},
			},
			role(),
			false,
		},
		{
			[]filter.Selector{
				{
					Group:     "rbac.authorization.k8s.io",
					Version:   "v1",
					Kind:      "Role",
					Name:      "test-role",
					Namespace: "",
				},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{
					Group:     "rbac.authorization.k8s.io",
					Version:   "v1",
					Kind:      "Role",
					Name:      "test-role",
					Namespace: "default",
				},
			},
			role(),
			true,
		},
		{
			[]filter.Selector{
				{
					Version:   "v1",
					Kind:      "ServiceAccount",
					Name:      "test-sa",
					Namespace: "monitoring",
				},
			},
			serviceAccount(),
			true,
		},
		{
			[]filter.Selector{
				{
					Version: "v1",
					Kind:    "ServiceAccount",
					Name:    "test-sa",
				},
			},
			serviceAccount(),
			true,
		},
		{
			[]filter.Selector{
				{
					Version:   "v1",
					Kind:      "ServiceAccount",
					Name:      "test-sa",
					Namespace: "default",
				},
			},
			serviceAccount(),
			false,
		},
	}

	for _, test := range tests {
		matcher := filter.SelectorMatcher(test.selectors)
		if result := matcher.Match(test.resource); result != test.expected {
			t.Errorf("expected %v for %v, got %v", test.expected, test.selectors, result)
			t.FailNow()
		}
	}
}

func role() unstructured.Unstructured {
	return unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "Role",
			"metadata": map[string]interface{}{
				"name": "test-role",
			},
		},
	}
}

func serviceAccount() unstructured.Unstructured {
	return unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      "test-sa",
				"namespace": "monitoring",
			},
		},
	}
}
