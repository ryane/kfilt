package main

import (
	"log"

	"github.com/ryane/kfilt/pkg/filter"
	"github.com/ryane/kfilt/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/yaml"
)

//noinspection GoUnusedGlobalVariable
//nolint: golint
type plugin struct {
	rf       *resmap.Factory
	Excludes []filter.Matcher `json:"excludes"`
	Includes []filter.Matcher `json:"includes"`
}

//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.rf = rf

	err := yaml.Unmarshal(c, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	return nil, nil
}

func (p *plugin) Transform(m resmap.ResMap) error {
	// setup kfilter
	kfilt := filter.New()

	for _, exclude := range p.Excludes {
		kfilt.AddExclude(exclude)
	}

	for _, include := range p.Includes {
		kfilt.AddInclude(include)
	}

	// copy kustomize resources into kfilt resources
	resources := make([]resource.Resource, m.Size())
	for i, res := range m.Resources() {
		log.Printf("resid: %s", res.CurId())
		resources[i] = resource.New(res.Map())
	}

	// create new resmap with filtered resources
	kresmap := resmap.New()
	kresourceFactory := p.rf.RF()
	for _, res := range kfilt.Filter(resources) {
		kres := kresourceFactory.FromMap(res.Object)
		if err := kresmap.Append(kres); err != nil {
			return err
		}
	}

	// reset the original resmap and add matching
	m.Clear()
	for _, res := range kresmap.Resources() {
		if err := m.Append(res); err != nil {
			return err
		}
	}

	return nil
}
