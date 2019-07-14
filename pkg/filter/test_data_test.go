package filter_test

import "github.com/ryane/kfilt/pkg/resource"

var input = []resource.Resource{
	resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa",
			},
		},
	),
	resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": "test-sa-2",
			},
		},
	),
	resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":      "test-pod",
				"namespace": "test-ns",
			},
		},
	),
	resource.New(
		map[string]interface{}{
			"apiVersion": "extensions/v1beta1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "test-deployment",
				"namespace": "test-ns",
			},
		},
	),
	resource.New(
		map[string]interface{}{
			"apiVersion": "extensions/v1beta1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "app",
				"namespace": "app",
			},
		},
	),
	resource.New(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "app",
				"namespace": "app",
			},
		},
	),
}
