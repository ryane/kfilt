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

func TestToList(t *testing.T) {
	// an actual list
	list := resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "List",
			"items": []interface{}{
				role().Object,
				serviceAccount().Object,
				map[string]interface{}{},
			},
		},
	)

	resources, err := list.ToList()
	if err != nil {
		t.Errorf("unexpected error expanding list: %v", err)
		t.FailNow()
	}

	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d\nresources: %+v", len(resources), resources)
		t.FailNow()
	}

	if resources[0].GetKind() != "Role" {
		t.Errorf("expected Role, got %s", resources[0].GetKind())
		t.FailNow()
	}

	if resources[1].GetKind() != "ServiceAccount" {
		t.Errorf("expected ServiceAccount, got %s", resources[1].GetKind())
		t.FailNow()
	}

	// a single resource
	role := role()
	resources, err = role.ToList()
	if err != nil {
		t.Errorf("unexpected error expanding list: %v", err)
		t.FailNow()
	}

	if len(resources) != 1 {
		t.Errorf("expected 1 resource, got %d\nresources: %+v", len(resources), resources)
		t.FailNow()
	}

	if resources[0].GetKind() != "Role" {
		t.Errorf("expected Role, got %s", resources[0].GetKind())
		t.FailNow()
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
