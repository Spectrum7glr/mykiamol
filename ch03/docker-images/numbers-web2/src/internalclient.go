package numbersweb2

import (
	"bytes"
	"encoding/json"
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

func (client *httpClient) makeRequest(method, url string, payload io.Reader, headers map[string]string) (string, error) {

	req, err := client.buildRequest(method, url, payload, headers)
	if err != nil {
		return "", err
	}

	response, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func (client *httpClient) buildRequest(method, url string, payload io.Reader, headers map[string]string) (*http.Request, error) {

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

func (client *httpClient) Get(url string, headers map[string]string) (string, error) {
	return client.makeRequest(http.MethodGet, url, nil, headers)
}

func (client *httpClient) Post(url string, payload any, headers map[string]string) (string, error) {

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return client.makeRequest(http.MethodPost, url, bytes.NewReader(payloadData), headers)
}
