package resource_test

import (
	"testing"

	"github.com/ryane/kfilt/pkg/resource"
)

func TestID(t *testing.T) {
	tests := []struct {
		resource.Resource
		expectedID string
	}{
		{role(), "rbac.authorization.k8s.io/v1:role::test-role"},
		{serviceAccount(), "/v1:serviceaccount:monitoring:test-sa"},
	}

	for _, test := range tests {
		id := test.Resource.ID()
		if id != test.expectedID {
			t.Errorf("expected %q, got %q for %+v", test.expectedID, id, test.Resource)
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
