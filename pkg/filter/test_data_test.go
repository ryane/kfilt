package filter_test

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var input = []unstructured.Unstructured{
	{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa-2",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":      "test-pod",
				"namespace": "test-ns",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1beta1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "test-deployment",
				"namespace": "test-ns",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1beta1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "app",
				"namespace": "app",
			},
		},
	},
	{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "app",
				"namespace": "app",
			},
		},
	},
}
