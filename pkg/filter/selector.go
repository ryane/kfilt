package filter

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Selector struct {
	Group     string
	Version   string
	Kind      string
	Name      string
	Namespace string
}

func (s *Selector) Match(u unstructured.Unstructured) bool {
	gvk := u.GroupVersionKind()

	if s.Group != "" && !strings.EqualFold(s.Group, gvk.Group) {
		return false
	}
	if s.Version != "" && !strings.EqualFold(s.Version, gvk.Version) {
		return false
	}
	if s.Kind != "" && !strings.EqualFold(s.Kind, gvk.Kind) {
		return false
	}
	if s.Name != "" && !strings.EqualFold(s.Name, u.GetName()) {
		return false
	}
	if s.Namespace != "" {
		ns := u.GetNamespace()
		if strings.ToLower(s.Namespace) == "default" {
			if ns != "" && !strings.EqualFold(s.Namespace, ns) {
				return false
			}
		} else if !strings.EqualFold(s.Namespace, ns) {
			return false
		}
	}

	return true
}
