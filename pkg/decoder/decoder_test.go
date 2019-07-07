package decoder_test

import (
	"os"
	"testing"

	"github.com/ryane/kfilt/pkg/decoder"
)

func TestDecoder(t *testing.T) {
	// load test data
	file := "./test.yaml"
	in, err := os.Open(file)
	if err != nil {
		t.Errorf("error loading test data: %v", err)
		t.FailNow()
	}
	defer in.Close()

	// decode
	d := decoder.New()
	results, err := d.Decode(in)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}

	expectedCount := 5
	if len(results) != expectedCount {
		t.Errorf("expected %d results, got %d", expectedCount, len(results))
		t.FailNow()
	}

	expectNames := []string{"test", "test2", "example-config", "cluster-specification", "handler"}
	for i, u := range results {
		name := u.GetName()
		if name != expectNames[i] {
			t.Errorf("expected %s, got %s", expectNames[i], name)
			t.FailNow()
		}
	}

	expectKinds := []string{"ServiceAccount", "ServiceAccount", "ConfigMap", "ClusterSpec", "stdio"}
	for i, u := range results {
		kind := u.GetKind()
		if kind != expectKinds[i] {
			t.Errorf("expected %s, got %s", expectKinds[i], kind)
			t.FailNow()
		}
	}
}
