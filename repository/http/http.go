package http

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPMirror struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) HTTPMirror {
	return HTTPMirror{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (hm HTTPMirror) GetResource(resource string) (io.Reader, error) {
	resp, err := hm.client.Get(hm.baseURL + resource)
	if err != nil {
		return nil, fmt.Errorf("getting resource %s from %s: %v", resource, hm.baseURL, err)
	}

	return resp.Body, nil
}

var DefaultHTTPMirror = New("https://github.com/lvrach/brave-mirror/raw/master/")
