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
	tests := []struct {
		exclude excludeMatchers
		include includeMatchers
		expectIDs
	}{
		// no filters, return all
		{
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
		},
		// exclude service accounts
		{
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
		},
		// exclude service accounts and pods
		{
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
		},
		// exclude deployments named "app"
		{
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
		},
		// include service accounts
		{
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
		},
		// include service accounts and pods
		{
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
		},
		// include service accounts and pods, but drop test-sa-2
		{
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
		},
		// don't include duplicate resources
		{
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
		},
		// label key selector
		{
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
		},
		// label key/value selector
		{
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
		},
		// label key/value selector
		{
			excludeMatchers{},
			includeMatchers{
				{
					LabelSelector: "app=test2",
				},
			},
			expectIDs{
				"/v1:serviceaccount::test-sa-2",
			},
		},
		// label != selector
		{
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
		},
	}

	for _, test := range tests {
		f := &filter.Filter{test.include, test.exclude}

		results := f.Filter(input)
		if len(results) != len(test.expectIDs) {
			t.Errorf("expected %d results, got %d\nincludes: %+v, excludes: %+v\nresults: %v", len(test.expectIDs), len(results), f.Include, f.Exclude, resourceIDs(results))
			t.FailNow()
		}

		for i, res := range results {
			id := res.ID()
			if id != test.expectIDs[i] {
				t.Errorf("expected %s, got %s\nincludes: %v, excludes: %v", test.expectIDs[i], id, f.Include, f.Exclude)
				t.FailNow()
			}
		}
	}
}

func resourceIDs(resources []resource.Resource) []string {
	ids := make([]string, len(resources))
	for i, res := range resources {
		ids[i] = res.ID()
	}
	return ids
}
