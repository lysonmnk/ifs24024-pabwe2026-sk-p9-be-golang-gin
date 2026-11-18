package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SearchClient interface {
	Search(query string) ([]SourceItem, error)
}

type searchClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewSearchClient(endpoint string) SearchClient {
	return &searchClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *searchClient) Search(query string) ([]SourceItem, error) {
	if c.endpoint == "" {
		return nil, nil
	}

	searchURL := fmt.Sprintf("%s?q=%s&format=json", c.endpoint, url.QueryEscape(query))
	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search engine returned status: %d", resp.StatusCode)
	}

	var raw SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var sources []SourceItem
	if raw.AbstractText != "" {
		sources = append(sources, SourceItem{
			Title:   raw.Heading,
			URL:     raw.AbstractURL,
			Snippet: raw.AbstractText,
		})
	}

	for _, rt := range raw.RelatedTopics {
		if rt.Text != "" && len(sources) < 5 {
			sources = append(sources, SourceItem{
				Title:   rt.Text,
				URL:     rt.FirstURL,
				Snippet: rt.Text,
			})
		}
	}

	if len(sources) == 0 {
		sources = append(sources, SourceItem{
			Title:   query,
			URL:     searchURL,
			Snippet: fmt.Sprintf("Informasi referensi untuk topik %s", query),
		})
	}

	return sources, nil
}
