package input_test

import (
	"io/ioutil"
	"os"
	"strings"
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
		expectedError  string
	}{
		{"", "", ""},
		{tmpfile, filedata, ""},
		{"https://httpbin.org/status/200", "", ""},
		{"https://httpbin.org/status/404", "", "404 Not Found"},
		{"/tmp/fake-file-that-doesnt-exist.fake", "", "no such file or directory"},
		{"https://broken\\url", "", "invalid character"},
		{"http://fake_host", "", "dial tcp"},
	}

	for _, test := range tests {
		r, err := input.Read(test.filename)

		if test.expectedError != "" {
			if !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf(
					"error does not contain %q: %v",
					test.expectedError,
					err,
				)
				t.FailNow()
			}
			continue
		}
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		defer r.Close()
		data, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

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
