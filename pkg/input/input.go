package input

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Read(f string) ([]byte, error) {
	if f == "" {
		return ioutil.ReadAll(os.Stdin)
	}

	if strings.HasPrefix(f, "http://") || strings.HasPrefix(f, "https://") {
		url, err := url.Parse(f)
		if err != nil {
			return nil, err
		}

		return httpGet(url)
	}

	return ioutil.ReadFile(f)
}

func httpGet(url *url.URL) ([]byte, error) {
	u := url.String()
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting %q: %s (%d)", u, resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
