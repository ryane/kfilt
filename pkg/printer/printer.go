package printer

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ryane/kfilt/pkg/resource"
	yaml "gopkg.in/yaml.v2"
)

type Printer interface {
	Print([]resource.Resource) error
}

type consolePrinter struct{}

func New() Printer {
	return &consolePrinter{}
}

func (p *consolePrinter) Print(resources []resource.Resource) error {
	for _, r := range resources {
		if err := p.printResource(r); err != nil {
			return err
		}
	}
	return nil
}

func (p *consolePrinter) printResource(r resource.Resource) error {
	data, err := yaml.Marshal(r.Object)
	if err != nil {
		return errors.Wrap(err, "failed to marshal yaml")
	}

	fmt.Println("---")
	fmt.Println(string(data))

	return nil
}
