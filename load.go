package JSONParse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func loadDoc(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return loadFromHTTP(path)
	} else {
		return ioutil.ReadFile(path)
	}
}

func loadFromHTTP(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else {
		return nil, fmt.Errorf("Not able to access document: %q", path)
	}
}
