package ai

type GenerateRequest struct {
	Topic            string `json:"topic" binding:"required,min=3"`
	Category         string `json:"category" binding:"required"`
	IncludeWebSearch bool   `json:"include_web_search"`
}

type SourceItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type GenerateResponse struct {
	Title    string       `json:"title"`
	Category string       `json:"category"`
	Excerpt  string       `json:"excerpt"`
	Content  string       `json:"content"`
	Sources  []SourceItem `json:"sources"`
}

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type SearchResponse struct {
	Heading      string `json:"Heading"`
	AbstractText string `json:"AbstractText"`
	AbstractURL  string `json:"AbstractURL"`
	RelatedTopics []struct {
		Text     string `json:"Text"`
		FirstURL string `json:"FirstURL"`
	} `json:"RelatedTopics"`
}
