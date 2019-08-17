package filter

import (
	"log"
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

	labelSelector labels.Selector // parsed LabelSelector
}

// Match returns true if a Resource matches the criteria
func (s *Matcher) Match(r resource.Resource) bool {
	gvk := r.GroupVersionKind()

	if s.Group != "" && !strings.EqualFold(s.Group, gvk.Group) {
		return false
	}
	if s.Version != "" && !strings.EqualFold(s.Version, gvk.Version) {
		return false
	}
	if s.Kind != "" && !strings.EqualFold(s.Kind, gvk.Kind) {
		return false
	}
	if s.Name != "" && !strings.EqualFold(s.Name, r.GetName()) {
		return false
	}
	if s.Namespace != "" {
		ns := r.GetNamespace()
		if strings.ToLower(s.Namespace) == "default" {
			if ns != "" && !strings.EqualFold(s.Namespace, ns) {
				return false
			}
		} else if !strings.EqualFold(s.Namespace, ns) {
			return false
		}
	}

	if s.LabelSelector != "" {
		selector, err := labels.Parse(s.LabelSelector)
		// TODO: fix this
		if err != nil {
			log.Fatal(err)
		}
		labelSet := labels.Set{}
		for name, val := range r.GetLabels() {
			labelSet[name] = val
		}

		if !selector.Matches(labelSet) {
			return false
		}
	}

	return true
}

// NewMatcher creates a Matcher
func NewMatcher(q string) (Matcher, error) {
	m := Matcher{}
	criteria := strings.Split(q, ",")

	if len(criteria) == 0 {
		return m, newMatcherParseError("invalid matcher %q. query is required", q)
	}

	for _, criterion := range criteria {
		parts := strings.Split(criterion, "=")
		if len(parts) != 2 {
			return m, newMatcherParseError("invalid matcher %q. Should be in the format %q", criterion, "key=value")
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
		default:
			return m, newMatcherParseError("invalid matcher %q. key should be one of %v", criterion, validMatcherKeys())
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
	}
}
