package mockllm

import (
	"context"
	"fmt"
	"iter"
	"os"
	"strings"
	"time"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
	"gopkg.in/yaml.v3"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

var Provider = llm_provider.ProviderMock

type Mock struct{}

func (Mock) Provider() model_ref.Provider { return Provider }

func (Mock) RequiresAPIKey() bool { return false }

func (Mock) Configure(llm_provider.Config) {}

func (Mock) ListModels(ctx context.Context) ([]model_ref.Info, error) {
	return nil, nil
}

func (Mock) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	return newMockModel(ref.Model.String())
}

func (Mock) ValidateModel(ctx context.Context, ref model_ref.ModelRef) error {
	_, err := newMockModel(ref.Model.String())
	return err
}

type mockScript struct {
	Steps []mockStep `yaml:"steps"`
}

type mockStep struct {
	Match   mockMatcher  `yaml:"match"`
	Respond mockResponse `yaml:"respond"`
}

type mockMatcher struct {
	Contains   string `yaml:"contains"`
	ToolResult string `yaml:"tool_result"`
	Any        bool   `yaml:"any"`
}

type mockResponse struct {
	Text      string         `yaml:"text"`
	ToolCalls []mockToolCall `yaml:"tool_calls"`
	Err       string         `yaml:"error"`
	Finish    string         `yaml:"finish"`
	DelayMS   int            `yaml:"delay_ms"`
}

type mockToolCall struct {
	ID   string         `yaml:"id"`
	Name string         `yaml:"name"`
	Args map[string]any `yaml:"args"`
}

type mockModel struct {
	script *mockScript
}

func newMockModel(scriptPath string) (model.LLM, error) {
	raw, err := os.ReadFile(scriptPath)
	if err != nil {
		return nil, err
	}

	script, err := parseMockScript(raw)
	if err != nil {
		return nil, err
	}

	return &mockModel{script: script}, nil
}

func (m *mockModel) Name() string { return "mock" }

func (m *mockModel) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		step, ok := m.matchRequest(req)
		if !ok {
			yield(nil, fmt.Errorf("mockllm: no step matched (contents=%d)", len(req.Contents)))
			return
		}

		if step.Respond.DelayMS > 0 {
			timer := time.NewTimer(time.Duration(step.Respond.DelayMS) * time.Millisecond)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			case <-timer.C:
			}
		}

		if step.Respond.Err != "" {
			yield(nil, fault.New(step.Respond.Err, fmsg.With(step.Respond.Err)))
			return
		}

		yield(buildMockResponse(step.Respond), nil)
	}
}

func (m *mockModel) matchRequest(req *model.LLMRequest) (mockStep, bool) {
	last := lastMockContent(req.Contents)
	for _, step := range m.script.Steps {
		if matchesMockContent(step.Match, last) {
			return step, true
		}
	}
	return mockStep{}, false
}

func matchesMockContent(mat mockMatcher, c *genai.Content) bool {
	switch {
	case mat.Contains != "":
		needle := strings.ToLower(mat.Contains)
		for _, part := range c.Parts {
			if part != nil && strings.Contains(strings.ToLower(part.Text), needle) {
				return true
			}
		}
		return false

	case mat.ToolResult != "":
		for _, part := range c.Parts {
			if part != nil && part.FunctionResponse != nil &&
				part.FunctionResponse.Name == mat.ToolResult {
				return true
			}
		}
		return false

	case mat.Any:
		return true
	}
	return false
}

func lastMockContent(contents []*genai.Content) *genai.Content {
	for i := len(contents) - 1; i >= 0; i-- {
		if contents[i] != nil {
			return contents[i]
		}
	}
	return &genai.Content{}
}

func buildMockResponse(r mockResponse) *model.LLMResponse {
	content := &genai.Content{Role: genai.RoleModel}

	if r.Text != "" {
		content.Parts = append(content.Parts, &genai.Part{Text: r.Text})
	}

	for _, tc := range r.ToolCalls {
		args := tc.Args
		if args == nil {
			args = map[string]any{}
		}
		content.Parts = append(content.Parts, &genai.Part{
			FunctionCall: &genai.FunctionCall{
				ID:   tc.ID,
				Name: tc.Name,
				Args: args,
			},
		})
	}

	finish := genai.FinishReasonStop
	switch r.Finish {
	case "max_tokens":
		finish = genai.FinishReasonMaxTokens
	case "safety":
		finish = genai.FinishReasonSafety
	}

	return &model.LLMResponse{
		Content:      content,
		FinishReason: finish,
		TurnComplete: true,
	}
}

func parseMockScript(data []byte) (*mockScript, error) {
	var s mockScript
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if len(s.Steps) == 0 {
		return nil, fmt.Errorf("mock script has no steps")
	}
	return &s, nil
}
