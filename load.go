package JSONParse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// loads a document, schema or schema external reference
// document can be loaded from file or over http
// the prefix of the path determines whether to access file
// system or internet.
//
// todo: support format of file://
func loadDoc(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return loadFromHTTP(path)
	} else {
		return ioutil.ReadFile(path)
	}
}

// uses http.Get to retrieve the document over http protocol
func loadFromHTTP(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		Error.Fatalln(err.Error())
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else {
		Error.Fatalln("Not able to access document: %q", path)
		return nil, fmt.Errorf("Not able to access document: %q", path)
	}
}
