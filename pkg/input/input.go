package input

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Read(f string) (io.ReadCloser, error) {
	if f == "" {
		return ioutil.NopCloser(os.Stdin), nil
	}

	if strings.HasPrefix(f, "http://") || strings.HasPrefix(f, "https://") {
		url, err := url.Parse(f)
		if err != nil {
			return nil, err
		}

		return httpGet(url)
	}

	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func httpGet(url *url.URL) (io.ReadCloser, error) {
	u := url.String()
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting %q: %s (%d)", u, resp.Status, resp.StatusCode)
	}

	return resp.Body, nil
}
