package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/r3labs/sse/v2"
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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
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

type Perplexity struct {
	endpoint    string
	apiKey      string
	model       string
	httpClient  *http.Client
	httpTimeout time.Duration
}

func newPerplexity(cfg config.Config) (*Perplexity, error) {
	s := &Perplexity{
		apiKey:      cfg.OpenAIKey,
		endpoint:    DefaultEndpoint,
		model:       Llama_3_1SonarSmall_128kChat,
		httpClient:  &http.Client{},
		httpTimeout: DefautTimeout,
	}
	return s, nil
}

func (s *Perplexity) Prompt(ctx context.Context, input string) (*Result, error) {
	r := &CompletionResponse{}

	reqBody := CompletionRequest{
		Messages: []Message{{Role: "user", Content: input}},
		Model:    s.model,
	}

	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(s.httpTimeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("unauthorized: check your API key")
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w - body response=%s", err, string(body))
	}

	return &Result{
		Answer: r.Choices[0].Message.Content,
	}, nil
}

func (s *Perplexity) PromptStream(ctx context.Context, input string) (chan string, chan error) {
	outch := make(chan string)
	errch := make(chan error)

	client := sse.NewClient(DefaultEndpoint)

	eventch := make(chan *sse.Event)

	go func() {
		err := client.SubscribeChan("completions", eventch)
		if err != nil {
			errch <- fault.Wrap(err, fctx.With(ctx))
			return
		}

		client.Unsubscribe(eventch)

		for e := range eventch {
			var cr CompletionResponse

			if err := json.Unmarshal(e.Data, &cr); err != nil {
				errch <- fault.Wrap(err, fctx.With(ctx))
				return
			}

			fmt.Println(cr.Citations)

			outch <- cr.Choices[0].Delta.Content

			if cr.Choices[0].FinishReason == "stop" {
				client.Unsubscribe(eventch)
				break
			}
		}

		close(outch)
		close(errch)
		close(eventch)
	}()

	return outch, errch
}

func (o *Perplexity) EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error) {
	panic("not implemented")
}
