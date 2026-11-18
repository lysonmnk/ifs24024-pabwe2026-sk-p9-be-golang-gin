package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSearchClient(t *testing.T) {
	t.Run("Empty endpoint returns nil, nil", func(t *testing.T) {
		client := NewSearchClient("")
		sources, err := client.Search("golang")
		assert.NoError(t, err)
		assert.Nil(t, sources)
	})

	t.Run("Valid endpoint with results", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := SearchResponse{
				Heading:      "Go Language",
				AbstractText: "Go is an open source programming language",
				AbstractURL:  "https://golang.org",
				RelatedTopics: []struct {
					Text     string `json:"Text"`
					FirstURL string `json:"FirstURL"`
				}{
					{Text: "Gin Web Framework", FirstURL: "https://gin-gonic.com"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		client := NewSearchClient(ts.URL)
		sources, err := client.Search("golang")
		assert.NoError(t, err)
		assert.Len(t, sources, 2)
		assert.Equal(t, "Go Language", sources[0].Title)
	})

	t.Run("Empty results fallback", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := SearchResponse{}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		client := NewSearchClient(ts.URL)
		sources, err := client.Search("golang")
		assert.NoError(t, err)
		assert.Len(t, sources, 1)
		assert.Equal(t, "golang", sources[0].Title)
	})

	t.Run("Search server error status and invalid json", func(t *testing.T) {
		tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer tsErr.Close()

		clientErr := NewSearchClient(tsErr.URL)
		_, err := clientErr.Search("golang")
		assert.Error(t, err)

		tsBadJSON := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer tsBadJSON.Close()

		clientBad := NewSearchClient(tsBadJSON.URL)
		_, err = clientBad.Search("golang")
		assert.Error(t, err)

		// Unreachable server error
		clientDead := NewSearchClient("http://127.0.0.1:59999")
		_, err = clientDead.Search("golang")
		assert.Error(t, err)

		// Invalid URL for NewRequest
		clientInvalidURL := NewSearchClient("http://invalid url with spaces")
		_, err = clientInvalidURL.Search("golang")
		assert.Error(t, err)
	})
}

func TestOpenAIClient(t *testing.T) {
	t.Run("Empty apiKey returns error", func(t *testing.T) {
		client := NewOpenAIClient("http://api.openai.com", "", "gpt-4o-mini")
		_, err := client.Generate("test")
		assert.Error(t, err)
		assert.Equal(t, "OPENAI_KEY tidak dikonfigurasi", err.Error())
	})

	t.Run("Success generation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
			resp := OpenAIResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{
					{
						Message: struct {
							Content string `json:"content"`
						}{
							Content: `{"title":"Judul","category":"Cat","excerpt":"Ex","content":"Con"}`,
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		client := NewOpenAIClient(ts.URL, "test-key", "gpt-4o-mini")
		content, err := client.Generate("test prompt")
		assert.NoError(t, err)
		assert.Contains(t, content, "Judul")
	})

	t.Run("Server error and empty choices and bad json", func(t *testing.T) {
		ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts500.Close()

		client500 := NewOpenAIClient(ts500.URL, "k", "m")
		_, err := client500.Generate("test")
		assert.Error(t, err)

		tsEmpty := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(OpenAIResponse{})
		}))
		defer tsEmpty.Close()

		clientEmpty := NewOpenAIClient(tsEmpty.URL, "k", "m")
		_, err = clientEmpty.Generate("test")
		assert.Error(t, err)

		tsBad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bad"))
		}))
		defer tsBad.Close()

		clientBad := NewOpenAIClient(tsBad.URL, "k", "m")
		_, err = clientBad.Generate("test")
		assert.Error(t, err)

		clientDead := NewOpenAIClient("http://127.0.0.1:59999", "k", "m")
		_, err = clientDead.Generate("test")
		assert.Error(t, err)

		clientInvalidURL := NewOpenAIClient("http://invalid url with spaces", "k", "m")
		_, err = clientInvalidURL.Generate("test")
		assert.Error(t, err)
	})
}

func TestAIService(t *testing.T) {
	mockSearch := &mockSearchClient{
		sources: []SourceItem{
			{Title: "Source 1", URL: "https://example.com/1", Snippet: "Snippet 1"},
		},
	}

	t.Run("GeneratePost with web search and json output", func(t *testing.T) {
		mockAI := &mockOpenAIClient{
			result: "```json\n{\"title\":\"Judul AI\",\"category\":\"Teknologi\",\"excerpt\":\"Excerpt AI\",\"content\":\"# Konten AI\"}\n```",
		}
		svc := NewService(mockSearch, mockAI)

		req := &GenerateRequest{
			Topic:            "Microservices Go",
			Category:         "Teknologi",
			IncludeWebSearch: true,
		}
		resp, err := svc.GeneratePost(req)
		assert.NoError(t, err)
		assert.Equal(t, "Judul AI", resp.Title)
		assert.Equal(t, "Teknologi", resp.Category)
		assert.Len(t, resp.Sources, 1)
	})

	t.Run("GeneratePost without web search and raw text fallback", func(t *testing.T) {
		mockAI := &mockOpenAIClient{
			result: "Ini adalah konten tulisan bebas tanpa format JSON.",
		}
		svc := NewService(mockSearch, mockAI)

		req := &GenerateRequest{
			Topic:            "Microservices Go",
			Category:         "Teknologi",
			IncludeWebSearch: false,
		}
		resp, err := svc.GeneratePost(req)
		assert.NoError(t, err)
		assert.Equal(t, "Mengenal Lebih Dalam Microservices Go", resp.Title)
		assert.Equal(t, "Teknologi", resp.Category)
		assert.Len(t, resp.Sources, 0)
	})

	t.Run("GeneratePost with empty category in json", func(t *testing.T) {
		mockAI := &mockOpenAIClient{
			result: `{"title":"Judul","category":"","excerpt":"Ex","content":"Con"}`,
		}
		svc := NewService(nil, mockAI)

		req := &GenerateRequest{
			Topic:            "Microservices",
			Category:         "DefaultCat",
			IncludeWebSearch: true,
		}
		resp, err := svc.GeneratePost(req)
		assert.NoError(t, err)
		assert.Equal(t, "DefaultCat", resp.Category)
	})

	t.Run("GeneratePost error from openai client", func(t *testing.T) {
		mockAI := &mockOpenAIClient{
			err: errors.New("ai error"),
		}
		svc := NewService(mockSearch, mockAI)

		req := &GenerateRequest{
			Topic:    "Topic",
			Category: "Cat",
		}
		_, err := svc.GeneratePost(req)
		assert.Error(t, err)
	})
}

type mockSearchClient struct {
	sources []SourceItem
	err     error
}

func (m *mockSearchClient) Search(q string) ([]SourceItem, error) {
	return m.sources, m.err
}

type mockOpenAIClient struct {
	result string
	err    error
}

func (m *mockOpenAIClient) Generate(p string) (string, error) {
	return m.result, m.err
}

func TestAIHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAI := &mockOpenAIClient{
		result: `{"title":"Judul","category":"Cat","excerpt":"Ex","content":"Con"}`,
	}
	svc := NewService(nil, mockAI)
	handler := NewHandler(svc)

	t.Run("GeneratePost success", func(t *testing.T) {
		body, _ := json.Marshal(GenerateRequest{
			Topic:            "Golang Web",
			Category:         "Teknologi",
			IncludeWebSearch: false,
		})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/posts/ai-generate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.GeneratePost(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GeneratePost invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/posts/ai-generate", bytes.NewBuffer([]byte("bad")))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.GeneratePost(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GeneratePost service error", func(t *testing.T) {
		errMock := &mockOpenAIClient{err: errors.New("quota exceeded")}
		errSvc := NewService(nil, errMock)
		errH := NewHandler(errSvc)

		body, _ := json.Marshal(GenerateRequest{
			Topic:    "Golang Web",
			Category: "Teknologi",
		})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/posts/ai-generate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		errH.GeneratePost(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
