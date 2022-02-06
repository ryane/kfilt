package filter_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"github.com/ryane/kfilt/pkg/resource"
)

type excludeMatchers []filter.Matcher
type includeMatchers []filter.Matcher
type expectIDs []string

func TestFilter(t *testing.T) {
	var noError = func(err error) bool { return err == nil }
	tests := []struct {
		name    string
		exclude excludeMatchers
		include includeMatchers
		expectIDs
		expectedError func(err error) bool
	}{
		{
			"no filters, return all",
			excludeMatchers{},
			includeMatchers{},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"exclude service accounts",
			excludeMatchers{
				{
					Kind: "ServiceAccount",
				},
			},
			includeMatchers{},
			expectIDs{
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"exclude service accounts and pods",
			excludeMatchers{
				{
					Kind: "ServiceAccount",
				},
				{
					Kind: "pod",
				},
			},
			includeMatchers{},
			expectIDs{
				"extensions/v1beta1:deployment:test-ns:test-deployment",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"exclude deployments named \"app\"",
			excludeMatchers{
				{
					Kind: "deployment",
					Name: "app",
				},
			},
			includeMatchers{},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"include service accounts",
			excludeMatchers{},
			includeMatchers{
				{
					Kind: "ServiceAccount",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
			},
			noError,
		},
		{
			"include service accounts and pods",
			excludeMatchers{},
			includeMatchers{
				{
					Kind: "ServiceAccount",
				},
				{
					Kind: "pod",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
			},
			noError,
		},
		{
			"include service accounts and pods, but drop test-sa-2",
			excludeMatchers{
				{
					Name: "test-sa-2",
				},
			},
			includeMatchers{
				{
					Kind: "ServiceAccount",
				},
				{
					Kind: "pod",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:pod:test-ns:test-pod",
			},
			noError,
		},
		{
			"don't include duplicate resources",
			excludeMatchers{},
			includeMatchers{
				{
					Kind: "ServiceAccount",
				},
				{
					Name: "test-sa-2",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
			},
			noError,
		},
		{
			"label key selector",
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
			},
			noError,
		},
		{
			"preserve order of returned resources",
			excludeMatchers{},
			includeMatchers{
				{
					Kind: "Pod",
				},
				{
					Kind: "ServiceAccount",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
			},
			noError,
		},
		{
			"label key/value selector",
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app=test",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
			},
			noError,
		},
		{
			"label key/value selector",
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app=test2",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa-2",
			},
			noError,
		},
		{
			"label != selector",
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app!=test",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa-2",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"exclude by label selector",
			excludeMatchers{
				{
					LabelSelector: "app==test",
				},
			},
			includeMatchers{},
			expectIDs{
				"/v1:serviceaccount::test-sa-2",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			noError,
		},
		{
			"exclude by label != selector",
			excludeMatchers{
				{
					LabelSelector: "app!=test",
				},
			},
			includeMatchers{},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
			},
			noError,
		},
		{
			"bad matcher filter include error",
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app===test",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa",
				"/v1:serviceaccount::test-sa-2",
				"/v1:pod:test-ns:test-pod",
				"extensions/v1beta1:deployment:test-ns:test-deployment",
				"extensions/v1beta1:deployment:app:app",
				"/v1:configmap:app:app",
			},
			filter.IsMatcherParseError,
		},
		{
			"bad matcher filter exclude error",
			excludeMatchers{
				{
					LabelSelector: "app===test",
				},
			},
			includeMatchers{},
			expectIDs{},
			filter.IsMatcherParseError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := filter.New()
			for _, m := range test.include {
				f.AddInclude(m)
			}
			for _, m := range test.exclude {
				f.AddExclude(m)
			}

			results, err := f.Filter(input)
			if !test.expectedError(err) {
				t.Errorf("unexpected error for %s: %v", test.name, err)
				t.FailNow()
			}

			if len(results) != len(test.expectIDs) {
				t.Errorf("%s: expected %d results, got %d\nincludes: %+v, excludes: %+v\nresults: %v", test.name, len(test.expectIDs), len(results), f.Include, f.Exclude, resourceIDs(results))
				t.FailNow()
			}

			for i, res := range results {
				id := res.ID()
				if id != test.expectIDs[i] {
					t.Errorf("%s: expected %s, got %s\nincludes: %v, excludes: %v", test.name, test.expectIDs[i], id, f.Include, f.Exclude)
					t.FailNow()
				}
			}
		})
	}
}

func resourceIDs(resources []resource.Resource) []string {
	ids := make([]string, len(resources))
	for i, res := range resources {
		ids[i] = res.ID()
	}
	return ids
}
