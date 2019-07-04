package filter

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type kindMatcher struct {
	kinds []string
}

func KindMatcher(kinds []string) Matcher {
	return &kindMatcher{validKinds(kinds)}
}

func (f *kindMatcher) Match(u unstructured.Unstructured) bool {
	// no kinds specified so we just return a match
	if len(f.kinds) == 0 {
		return true
	}

	for _, kind := range f.kinds {
		if strings.EqualFold(kind, u.GetKind()) {
			return true
		}
	}

	return false
}

func validKinds(kinds []string) []string {
	valid := []string{}
	for _, kind := range kinds {
		trimKind := strings.TrimSpace(kind)
		if len(trimKind) > 0 {
			valid = append(valid, trimKind)
		}
	}
	return valid
}
