package decoder_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ryane/kfilt/pkg/decoder"
)

func TestDecoder(t *testing.T) {
	// list.yaml contains the same resources as test.yaml but is formatted as a
	// List
	files := []string{"./test.yaml", "./list.yaml"}
	for _, file := range files {
		// load test data
		in, err := os.Open(file)
		if err != nil {
			t.Errorf("error loading test data from %s: %v", file, err)
			t.FailNow()
		}
		defer in.Close()

		// decode
		d := decoder.New()
		results, err := d.Decode(in)

		if err != nil {
			t.Errorf("unexpected error decoding %s: %v", file, err)
			t.FailNow()
		}

		expectedCount := 5
		if len(results) != expectedCount {
			t.Errorf("expected %d results from %s, got %d", expectedCount, file, len(results))
			t.FailNow()
		}

		expectNames := []string{"test", "test2", "example-config", "cluster-specification", "handler"}
		for i, res := range results {
			name := res.GetName()
			if name != expectNames[i] {
				t.Errorf("expected %s from %s, got %s", expectNames[i], file, name)
				t.FailNow()
			}
		}

		expectKinds := []string{"ServiceAccount", "ServiceAccount", "ConfigMap", "ClusterSpec", "stdio"}
		for i, res := range results {
			kind := res.GetKind()
			if kind != expectKinds[i] {
				t.Errorf("expected %s from %s, got %s", expectKinds[i], file, kind)
				t.FailNow()
			}
		}
	}
}

func TestDecoderErrors(t *testing.T) {
	tests := []struct {
		filename      string
		expectedError string
	}{
		{"./bad.yaml", "failed to decode input"},
		{"./bad_list.yaml", "failed to explode list"},
	}
	for _, test := range tests {
		// load test data
		in, err := os.Open(test.filename)
		if err != nil {
			t.Errorf("error loading test data from %s: %v", test.filename, err)
			t.FailNow()
		}
		defer in.Close()

		d := decoder.New()
		_, err = d.Decode(in)
		if err == nil {
			t.Errorf("expected error %q, got nil", test.expectedError)
			t.FailNow()
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf(
				"error does not contain %q: %v",
				test.expectedError,
				err,
			)
			t.Fail()
		}
	}
}
