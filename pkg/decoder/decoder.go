package decoder

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type Decoder interface {
	Decode([]byte) ([]unstructured.Unstructured, error)
}

type kubernetesDecoder struct{}

func New() Decoder {
	return &kubernetesDecoder{}
}

func (k *kubernetesDecoder) Decode(in []byte) ([]unstructured.Unstructured, error) {
	var (
		result []unstructured.Unstructured
		err    error
	)

	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(in), 1024)

	for err == nil {
		var out unstructured.Unstructured
		err = decoder.Decode(&out)
		if err == nil {
			result = append(result, out)
		}
	}
	if err != io.EOF {
		return nil, errors.Wrap(err, "failed to decode input")
	}

	return result, nil
}
