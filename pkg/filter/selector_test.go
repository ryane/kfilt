package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewSelector(t *testing.T) {
	var noError = func(err error) bool { return err == nil }
	tests := []struct {
		query         string
		expected      filter.Selector
		expectedError func(err error) bool
	}{
		// kind
		{
			"kind=secret",
			filter.Selector{Kind: "secret"},
			noError,
		},
		// kind abbreviation
		{
			"k=secret",
			filter.Selector{Kind: "secret"},
			func(error) bool { return true },
		},
		// last kind winws
		{
			"kind=Secret,kind=ServiceAccount",
			filter.Selector{Kind: "ServiceAccount"},
			noError,
		},
		// name
		{
			"name=test",
			filter.Selector{Name: "test"},
			noError,
		},
		// name abbreviation
		{
			"n=test",
			filter.Selector{Name: "test"},
			noError,
		},
		// no matcher
		{
			"",
			filter.Selector{},
			filter.IsMatcherError,
		},
		// bad key
		{
			"wtf=happened",
			filter.Selector{},
			filter.IsMatcherError,
		},
	}

	for _, test := range tests {
		result, err := filter.NewSelector(test.query)

		if !test.expectedError(err) {
			t.Errorf("unexpected error %v for %q", err, test.query)
			t.FailNow()
		}

		if result != test.expected {
			t.Errorf("expected %#v for %s, got %v", test.expected, test.query, result)
			t.FailNow()
		}
	}
}

func TestSelector(t *testing.T) {
	tests := []struct {
		selectors []filter.Selector // TODO: convert this to a single selector
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
		for _, selector := range test.selectors {
			if result := selector.Match(test.resource); result != test.expected {
				t.Errorf("expected %v for %v, got %v", test.expected, test.selectors, result)
				t.FailNow()
			}
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
