package dns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func Update(t string, name string, content string, token string, ttl int) error {
	// Setting to `1` means 'automatic' on Cloudflare.
	if ttl == 0 {
		ttl = 1
	}

	zoneId, err := getZoneId(name, token)
	if err != nil {
		return err
	}

	recordId, err := getRecordId(zoneId, token, name)
	if err != nil {
		return err
	}

	if recordId == "" {
		err := createRecord(zoneId, token, t, name, content, false, ttl)
		if err != nil {
			return err
		}
	}

	return updateRecord(token, zoneId, recordId, t, name, content, false, ttl)
}

func getZoneId(name string, token string) (string, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	result, err := request("https://api.cloudflare.com/client/v4/zones", "GET", nil, headers)
	if err != nil {
		return "", err
	}

	success := result["success"].(bool)
	if !success {
		errs := result["errors"].([]any)
		message := errs[0].(map[string]any)["message"].(string)
		return "", errors.New("DNS error: " + message)
	}

	for _, zone := range result["result"].([]any) {
		zone := zone.(map[string]any)
		if strings.HasSuffix(name, zone["name"].(string)) {
			return zone["id"].(string), nil
		}
	}

	return "", errors.New("DNS error: Fail to find the zone.")
}

func getRecordId(zoneId string, token string, name string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records", zoneId)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	result, err := request(url, "GET", nil, headers)
	if err != nil {
		return "", err
	}

	success := result["success"].(bool)
	if !success {
		errs := result["errors"].([]any)
		message := errs[0].(map[string]any)["message"].(string)
		return "", errors.New("DNS error: " + message)
	}

	for _, record := range result["result"].([]any) {
		record := record.(map[string]any)
		if record["name"].(string) == name {
			return record["id"].(string), nil
		}
	}

	return "", nil
}

type cloudflarePayload struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
}

func createRecord(zoneId string, token string, t string, name string, content string, proxied bool, ttl int) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records", zoneId)
	payload := cloudflarePayload{Type: t, Name: name, Content: content, Proxied: proxied, TTL: ttl}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(encoded)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	_, err = request(url, "POST", reader, headers)
	return err
}

func updateRecord(token string, zoneId string, recordId string, t string, name string, content string, proxied bool, ttl int) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", zoneId, recordId)
	payload := cloudflarePayload{Type: t, Name: name, Content: content, Proxied: proxied, TTL: ttl}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(encoded)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	_, err = request(url, "PATCH", reader, headers)
	return err
}
