package filter_test

import (
	"strings"
	"testing"

	"github.com/ryane/kfilt/pkg/filter"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type excludeMatchers []filter.Matcher
type includeMatchers []filter.Matcher
type expectGVKNS []string

func TestFilter(t *testing.T) {
	tests := []struct {
		exclude     excludeMatchers
		include     includeMatchers
		expectNames []string
	}{
		// no filters, return all
		{
			excludeMatchers{},
			includeMatchers{},
			expectGVKNS{
				"/v1:serviceaccount:test-sa",
				"/v1:serviceaccount:test-sa-2",
				"/v1:pod:test-pod",
				"extensions/v1beta1:deployment:test-deployment",
				"extensions/v1beta1:deployment:app",
				"/v1:configmap:app",
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
			expectGVKNS{
				"/v1:pod:test-pod",
				"extensions/v1beta1:deployment:test-deployment",
				"extensions/v1beta1:deployment:app",
				"/v1:configmap:app",
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
			expectGVKNS{
				"extensions/v1beta1:deployment:test-deployment",
				"extensions/v1beta1:deployment:app",
				"/v1:configmap:app",
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
			expectGVKNS{
				"/v1:serviceaccount:test-sa",
				"/v1:serviceaccount:test-sa-2",
				"/v1:pod:test-pod",
				"extensions/v1beta1:deployment:test-deployment",
				"/v1:configmap:app",
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
			expectGVKNS{
				"/v1:serviceaccount:test-sa",
				"/v1:serviceaccount:test-sa-2",
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
			expectGVKNS{
				"/v1:serviceaccount:test-sa",
				"/v1:serviceaccount:test-sa-2",
				"/v1:pod:test-pod",
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
			expectGVKNS{
				"/v1:serviceaccount:test-sa",
				"/v1:pod:test-pod",
			},
		},
	}

	for _, test := range tests {
		f := &filter.Filter{test.include, test.exclude}

		results := f.Filter(input)
		if len(results) != len(test.expectNames) {
			t.Errorf("expected %d results, got %d\nincludes: %v, excludes: %v", len(test.expectNames), len(results), f.Include, f.Exclude)
			t.FailNow()
		}

		for i, u := range results {
			id := gvkn(u)
			if gvkn(u) != test.expectNames[i] {
				t.Errorf("expected %s, got %s\nincludes: %v, excludes: %v", test.expectNames[i], id, f.Include, f.Exclude)
				t.FailNow()
			}
		}
	}
}

func gvkn(u unstructured.Unstructured) string {
	gvk := u.GroupVersionKind()
	return strings.ToLower(
		gvk.Group + "/" + gvk.Version + ":" + gvk.Kind + ":" + u.GetName(),
	)
}
