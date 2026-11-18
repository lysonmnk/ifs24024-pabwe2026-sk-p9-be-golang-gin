package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OpenAIClient interface {
	Generate(prompt string) (string, error)
}

type openAIClient struct {
	endpoint   string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(endpoint, apiKey, model string) OpenAIClient {
	return &openAIClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *openAIClient) Generate(prompt string) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("OPENAI_KEY tidak dikonfigurasi")
	}

	apiURL := strings.TrimRight(c.endpoint, "/") + "/chat/completions"

	body := OpenAIRequest{
		Model: c.model,
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "Anda adalah asisten pembuat draf blog profesional. Kembalikan response HANYA dalam format JSON dengan kunci: title, category, excerpt, content.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	jsonBytes, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai returned status: %d", resp.StatusCode)
	}

	var raw OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	if len(raw.Choices) == 0 {
		return "", errors.New("tidak ada respon dari model AI")
	}

	return raw.Choices[0].Message.Content, nil
}
