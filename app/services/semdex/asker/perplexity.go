package asker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	DefaultEndpoint = "https://api.perplexity.ai/chat/completions"
	DefautTimeout   = 10 * time.Second
)

const (
	Llama_3_1SonarSmall_128kChat   = "llama-3.1-sonar-small-128k-chat"
	Llama_3_1SonarLarge_128kChat   = "llama-3.1-sonar-large-128k-chat"
	Llama_3_1SonarSmall_128kOnline = "llama-3.1-sonar-small-128k-online"
	Llama_3_1SonarLarge_128kOnline = "llama-3.1-sonar-large-128k-online"
	Llama_3_1_8bInstruct           = "llama-3.1-8b-instruct"
	Llama_3_1_70bInstruct          = "llama-3.1-70b-instruct"
)

var citationMarkerRegex = regexp.MustCompile(`\[\d+\]`)

type Perplexity struct {
	endpoint    string
	apiKey      string
	model       string
	httpClient  *http.Client
	httpTimeout time.Duration
	searcher    semdex.Searcher
}

func newPerplexityAsker(cfg config.Config, searcher semdex.Searcher) (*Perplexity, error) {
	s := &Perplexity{
		apiKey:      cfg.PerplexityAPIKey,
		endpoint:    DefaultEndpoint,
		model:       Llama_3_1SonarSmall_128kOnline,
		httpClient:  &http.Client{},
		httpTimeout: DefautTimeout,
		searcher:    searcher,
	}
	return s, nil
}

func (a *Perplexity) Ask(ctx context.Context, q string, parent opt.Optional[xid.ID]) (semdex.AskResponseIterator, error) {
	t, err := buildContextPrompt(ctx, a.searcher, q)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resp, err := func() (*http.Response, error) {
		reqBody := CompletionRequest{
			Stream:   true,
			Messages: []Message{{Role: "user", Content: t}},
			Model:    a.model,
		}

		requestBody, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+a.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		return resp, nil
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	dec := ssestream.NewDecoder(resp)

	iter := func(yield func(string, error) bool) {
		defer resp.Body.Close()

		for dec.Next() {
			event := dec.Event()
			var cr CompletionResponse

			if err := json.Unmarshal(event.Data, &cr); err != nil {
				yield("", fmt.Errorf("failed to unmarshal SSE event: %w", err))
				return
			}

			if len(cr.Choices) == 0 {
				yield("", fmt.Errorf("no choices in response"))
				return
			}

			if len(cr.Citations) == 0 {
				yield("", fmt.Errorf("no citations in response"))
				return
			}

			choice := cr.Choices[0]

			chunk := choice.Delta.Content

			// replace [1]/[2]/etc citation markers with empty string
			// cleaned := citationMarkerRegex.ReplaceAllString(chunk, "")

			if !yield(chunk, nil) {
				return
			}

			if choice.FinishReason == "stop" {
				return
			}
		}

		if dec.Err() != nil {
			yield("", fmt.Errorf("failed to read SSE stream: %w", dec.Err()))
		}
	}

	return streamExtractor(iter), nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
	Delta        Message `json:"delta"`
}

type CompletionResponse struct {
	ID        string   `json:"id"`
	Model     string   `json:"model"`
	Created   int      `json:"created"`
	Usage     Usage    `json:"usage"`
	Citations []string `json:"citations"`
	Object    string   `json:"object"`
	Choices   []Choice `json:"choices"`
}
