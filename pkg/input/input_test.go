package input_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ryane/kfilt/pkg/input"
)

func TestRead(t *testing.T) {
	filedata := "---"
	tmpfile, err := createTempFile(filedata)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.Remove(tmpfile)

	tests := []struct {
		filename       string
		expectedString string
	}{
		{"", ""},
		{tmpfile, filedata},
		// TODO: fix me
		{"https://knative.dev/robots.txt", "User-agent: *"},
	}

	for _, test := range tests {
		data, _ := input.Read(test.filename)
		if string(data) != test.expectedString {
			t.Errorf("expected %s, got %s", test.expectedString, string(data))
			t.FailNow()
		}
	}
}

func createTempFile(data string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "input_test")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(data)); err != nil {
		return "", err
	}

	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
