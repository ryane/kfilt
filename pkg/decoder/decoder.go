package decoder

import (
	"io"

	"github.com/pkg/errors"
	"github.com/ryane/kfilt/pkg/resource"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type Decoder interface {
	Decode(io.Reader) ([]resource.Resource, error)
}

type kubernetesDecoder struct{}

func New() Decoder {
	return &kubernetesDecoder{}
}

func (k *kubernetesDecoder) Decode(in io.Reader) ([]resource.Resource, error) {
	var (
		result []resource.Resource
		err    error
	)

	decoder := k8syaml.NewYAMLOrJSONDecoder(in, 1024)

	for err == nil {
		var out resource.Resource
		err = decoder.Decode(&out)
		if err == nil && len(out.Object) > 0 {
			result = append(result, out)
		}
	}
	if err != io.EOF {
		return nil, errors.Wrap(err, "failed to decode input")
	}

	return result, nil
}
