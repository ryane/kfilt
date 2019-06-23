package filter_test

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var input = []unstructured.Unstructured{
	{
		Object: map[string]interface{}{
			"kind": "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"kind": "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa-2",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"kind": "Pod",
			"metadata": map[string]interface{}{
				"name": "test-pod",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"kind": "Deployment",
			"metadata": map[string]interface{}{
				"name": "test-deployment",
			},
		},
	},
}
