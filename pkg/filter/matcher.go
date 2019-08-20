package filter

import (
	"strings"

	"github.com/ryane/kfilt/pkg/resource"
	"k8s.io/apimachinery/pkg/labels"
)

// Matcher represents match criteria
type Matcher struct {
	Group         string
	Version       string
	Kind          string
	Name          string
	Namespace     string
	LabelSelector string
}

// Match returns true if a Resource matches the criteria
func (s *Matcher) Match(r resource.Resource) (bool, error) {
	gvk := r.GroupVersionKind()

	if s.Group != "" && !strings.EqualFold(s.Group, gvk.Group) {
		return false, nil
	}
	if s.Version != "" && !strings.EqualFold(s.Version, gvk.Version) {
		return false, nil
	}
	if s.Kind != "" && !strings.EqualFold(s.Kind, gvk.Kind) {
		return false, nil
	}
	if s.Name != "" && !strings.EqualFold(s.Name, r.GetName()) {
		return false, nil
	}
	if s.Namespace != "" {
		ns := r.GetNamespace()
		if strings.ToLower(s.Namespace) == "default" {
			if ns != "" && !strings.EqualFold(s.Namespace, ns) {
				return false, nil
			}
		} else if !strings.EqualFold(s.Namespace, ns) {
			return false, nil
		}
	}

	if s.LabelSelector != "" {
		selector, err := labels.Parse(s.LabelSelector)
		if err != nil {
			return false, newMatcherParseError("invalid label selector: %v", err)
		}
		labelSet := labels.Set{}
		for name, val := range r.GetLabels() {
			labelSet[name] = val
		}

		if !selector.Matches(labelSet) {
			return false, nil
		}
	}

	return true, nil
}

// NewMatcher creates a Matcher
func NewMatcher(q string) (Matcher, error) {
	m := Matcher{}

	var criteria []string
	parts := strings.Split(q, ",")
	for _, part := range parts {
		if part != "" {
			criteria = append(criteria, part)
		}
	}

	if len(criteria) == 0 {
		return m, newMatcherParseError(
			"invalid matcher %q. query is required",
			q,
		)
	}

	for _, criterion := range criteria {
		parts := strings.SplitN(criterion, "=", 2)
		if len(parts) != 2 {
			return m, newMatcherParseError(
				"invalid matcher %q. Should be in the format %q",
				criterion,
				"key=value",
			)
		}

		key, val := strings.ToLower(parts[0]), parts[1]

		switch key {
		case "kind", "k":
			m.Kind = val
		case "name", "n":
			m.Name = val
		case "group", "g":
			m.Group = val
		case "version", "v":
			m.Version = val
		case "namespace", "ns":
			m.Namespace = val
		case "labels", "l":
			m.LabelSelector = val
		default:
			return m, newMatcherParseError(
				"invalid matcher %q. key should be one of %v",
				criterion,
				validMatcherKeys(),
			)
		}
	}

	return m, nil
}

func validMatcherKeys() []string {
	return []string{
		"kind", "k",
		"name", "n",
		"group", "g",
		"version", "v",
		"namespace", "ns",
		"labels", "l",
	}
}
