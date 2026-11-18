package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Service interface {
	GeneratePost(req *GenerateRequest) (*GenerateResponse, error)
}

type aiService struct {
	searchClient SearchClient
	openaiClient OpenAIClient
}

func NewService(searchClient SearchClient, openaiClient OpenAIClient) Service {
	return &aiService{
		searchClient: searchClient,
		openaiClient: openaiClient,
	}
}

func (s *aiService) GeneratePost(req *GenerateRequest) (*GenerateResponse, error) {
	var sources []SourceItem

	if req.IncludeWebSearch && s.searchClient != nil {
		foundSources, err := s.searchClient.Search(req.Topic)
		if err == nil && len(foundSources) > 0 {
			sources = foundSources
		}
	}

	prompt := fmt.Sprintf("Buat draf artikel blog dengan topik: '%s' dan kategori: '%s'.", req.Topic, req.Category)
	if len(sources) > 0 {
		prompt += "\nGunakan referensi berikut sebagai konteks dan cantumkan dalam pembahasan:\n"
		for i, src := range sources {
			prompt += fmt.Sprintf("[%d] %s: %s (%s)\n", i+1, src.Title, src.Snippet, src.URL)
		}
	}

	aiContent, err := s.openaiClient.Generate(prompt)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Excerpt  string `json:"excerpt"`
		Content  string `json:"content"`
	}

	cleanJSON := strings.TrimSpace(aiContent)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	if err := json.Unmarshal([]byte(cleanJSON), &parsed); err != nil {
		parsed.Title = fmt.Sprintf("Mengenal Lebih Dalam %s", req.Topic)
		parsed.Category = req.Category
		parsed.Excerpt = fmt.Sprintf("Ulasan komprehensif mengenai %s.", req.Topic)
		parsed.Content = aiContent
	}

	if parsed.Category == "" {
		parsed.Category = req.Category
	}

	return &GenerateResponse{
		Title:    parsed.Title,
		Category: parsed.Category,
		Excerpt:  parsed.Excerpt,
		Content:  parsed.Content,
		Sources:  sources,
	}, nil
}
