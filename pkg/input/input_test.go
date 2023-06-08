package input_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ryane/kfilt/pkg/input"
)

func waitForServer(URL string, timeout time.Duration) error {
	type ok struct{}

	done := make(chan ok)
	go func() {
		for {
			_, err := http.Get(URL)
			if err == nil {
				done <- ok{}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("server did not reply after %v", timeout)
	}
}

func TestRead(t *testing.T) {
	filedata := "---"
	tmpfile, err := createTempFile(filedata)
	if err != nil {
		t.Error(err, "unable to create temp file")
		t.FailNow()
	}
	defer os.Remove(tmpfile)

	http.HandleFunc("/404", http.NotFound)
	http.HandleFunc("/200", func(http.ResponseWriter, *http.Request) {})
	go func() {
		_ = http.ListenAndServe(":8822", nil)
	}()

	if err := waitForServer("http://localhost:8822", time.Second*3); err != nil {
		t.Error(err, "test server not responding")
		t.FailNow()
	}

	tests := []struct {
		filename       string
		expectedString string
		expectedError  string
	}{
		{"", "", ""},
		{tmpfile, filedata, ""},
		{"http://localhost:8822/200", "", ""},
		{"http://localhost:8822/404", "", "404 Not Found"},
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
		data, err := io.ReadAll(r)
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
	tmpfile, err := os.CreateTemp("", "input_test")
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
