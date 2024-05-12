package simplecall

import (
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	Client http.Client
}

func InitHttpClient(timeout time.Duration) *httpClient {
	return &httpClient{
		Client: http.Client{
			Timeout: timeout,
		},
	}
}

func (client *httpClient) BuildRequest(method, url string, payload io.Reader, headers map[string]string) (*http.Request, error) {

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
