package printer

import (
	"fmt"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Printer interface {
	Print([]unstructured.Unstructured) error
}

type consolePrinter struct{}

func New() Printer {
	return &consolePrinter{}
}

func (p *consolePrinter) Print(unstructureds []unstructured.Unstructured) error {
	for _, u := range unstructureds {
		if err := p.printUnstructured(u); err != nil {
			return err
		}
	}
	return nil
}

func (p *consolePrinter) printUnstructured(u unstructured.Unstructured) error {
	data, err := yaml.Marshal(u.Object)
	if err != nil {
		return errors.Wrap(err, "failed to marshal yaml")
	}

	fmt.Println("---")
	fmt.Println(string(data))

	return nil
}
