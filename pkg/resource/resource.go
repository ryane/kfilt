package resource

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resource struct {
	unstructured.Unstructured
}

func New(obj map[string]interface{}) Resource {
	return Resource{unstructured.Unstructured{Object: obj}}
}

func (r *Resource) ID() string {
	gvk := r.GroupVersionKind()
	return strings.ToLower(
		gvk.Group + "/" + gvk.Version + ":" + gvk.Kind + ":" + r.GetNamespace() + ":" + r.GetName(),
	)
}
