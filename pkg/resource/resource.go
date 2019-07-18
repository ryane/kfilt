package resource

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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

func (r *Resource) ToList() ([]Resource, error) {
	if !r.IsList() {
		return []Resource{*r}, nil
	}

	resources := []Resource{}
	err := r.EachListItem(func(o runtime.Object) error {
		u := o.(*unstructured.Unstructured)
		if len(u.Object) > 0 {
			resources = append(resources, New(u.Object))
		}
		return nil
	})

	return resources, err
}
