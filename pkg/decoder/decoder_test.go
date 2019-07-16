package decoder_test

import (
	"os"
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
