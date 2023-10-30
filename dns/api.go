package dns

import (
	"encoding/json"
	"io"
	"net/http"
)

func request(url string, method string, payload io.Reader, headers map[string]string) (map[string]any, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]any)
	return parsed, json.Unmarshal(body, &parsed)
}
