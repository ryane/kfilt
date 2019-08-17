package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"github.com/ryane/kfilt/pkg/resource"
)

// TODO: add label matcher tests
func TestNewMatcher(t *testing.T) {
	var noError = func(err error) bool { return err == nil }
	tests := []struct {
		query         string
		expected      filter.Matcher
		expectedError func(err error) bool
	}{
		// kind
		{
			"kind=secret",
			filter.Matcher{Kind: "secret"},
			noError,
		},
		// kind abbreviation
		{
			"k=secret",
			filter.Matcher{Kind: "secret"},
			func(error) bool { return true },
		},
		// last kind winws
		{
			"kind=Secret,kind=ServiceAccount",
			filter.Matcher{Kind: "ServiceAccount"},
			noError,
		},
		// name
		{
			"name=test",
			filter.Matcher{Name: "test"},
			noError,
		},
		// name abbreviation
		{
			"n=test",
			filter.Matcher{Name: "test"},
			noError,
		},
		// no matcher
		{
			"",
			filter.Matcher{},
			filter.IsMatcherParseError,
		},
		// bad key
		{
			"wtf=happened",
			filter.Matcher{},
			filter.IsMatcherParseError,
		},
	}

	for _, test := range tests {
		result, err := filter.NewMatcher(test.query)

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

// TODO: add label matcher tests
func TestMatcher(t *testing.T) {
	tests := []struct {
		matcher  filter.Matcher
		resource resource.Resource
		expected bool
	}{
		// empty matcher should match
		{filter.Matcher{}, role(), true},
		// kind matchers
		{
			filter.Matcher{
				Kind: "Role",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Kind: "role",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Kind: "ServiceAccount",
			},
			role(),
			false,
		},
		{
			filter.Matcher{
				Version: "v1", Kind: "ServiceAccount",
			},
			serviceAccount(),
			true,
		},
		{
			filter.Matcher{
				Group: "", Version: "v1", Kind: "ServiceAccount",
			},
			serviceAccount(),
			true,
		},
		{
			filter.Matcher{
				Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Group: "rbac.authorization.k8s.io", Version: "v1beta1", Kind: "Role",
			},
			role(),
			false,
		},
		{
			filter.Matcher{
				Group:   "rbac.authorization.k8s.io",
				Version: "v1",
				Kind:    "Role",
				Name:    "test-role",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Group:   "rbac.authorization.k8s.io",
				Version: "v1",
				Kind:    "Role",
				Name:    "test-sa",
			},
			role(),
			false,
		},
		{
			filter.Matcher{
				Group:     "rbac.authorization.k8s.io",
				Version:   "v1",
				Kind:      "Role",
				Name:      "test-role",
				Namespace: "",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Group:     "rbac.authorization.k8s.io",
				Version:   "v1",
				Kind:      "Role",
				Name:      "test-role",
				Namespace: "default",
			},
			role(),
			true,
		},
		{
			filter.Matcher{
				Version:   "v1",
				Kind:      "ServiceAccount",
				Name:      "test-sa",
				Namespace: "monitoring",
			},
			serviceAccount(),
			true,
		},
		{
			filter.Matcher{
				Version: "v1",
				Kind:    "ServiceAccount",
				Name:    "test-sa",
			},
			serviceAccount(),
			true,
		},
		{
			filter.Matcher{
				Version:   "v1",
				Kind:      "ServiceAccount",
				Name:      "test-sa",
				Namespace: "default",
			},
			serviceAccount(),
			false,
		},
	}

	for _, test := range tests {
		if result := test.matcher.Match(test.resource); result != test.expected {
			t.Errorf("expected %v for %v, got %v", test.expected, test.matcher, result)
			t.FailNow()
		}
	}
}

func role() resource.Resource {
	return resource.New(
		map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "Role",
			"metadata": map[string]interface{}{
				"name": "test-role",
			},
		},
	)
}

func serviceAccount() resource.Resource {
	return resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      "test-sa",
				"namespace": "monitoring",
			},
		},
	)
}
