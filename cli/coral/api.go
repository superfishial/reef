package coral

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func makeRequest(token string, baseURL string, paths ...string) (req *http.Request, err error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}

	req, err = http.NewRequest("GET", url.JoinPath(paths...).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

func Get(token string, baseURL string, filename string) (io.ReadCloser, error) {
	req, err := makeRequest(token, baseURL, "/v1", "/coral", filename)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get file with status %s", resp.Status)
	}
	return resp.Body, nil
}
