package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// FetchModels returns the model IDs available on an OpenAI-compatible endpoint
// by calling GET {baseURL}/models.
func FetchModels(baseURL, apiKey string) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := strings.TrimSuffix(baseURL, "/") + "/models"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("failed to fetch models: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(result.Data))
	for _, m := range result.Data {
		if m.ID != "" {
			ids = append(ids, m.ID)
		}
	}
	return ids, nil
}
