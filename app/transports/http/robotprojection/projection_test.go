package robotprojection

import (
	"encoding/json"
	"testing"

	"google.golang.org/genai"
)

func TestFunctionResponseToUIPartDoesNotCopyOutputIntoInput(t *testing.T) {
	part, err := FunctionResponseToUIPart(&genai.FunctionResponse{
		ID:   "call_1",
		Name: "content_search",
		Response: map[string]any{
			"items":   []any{},
			"results": 0,
		},
	})
	if err != nil {
		t.Fatalf("FunctionResponseToUIPart() error = %v", err)
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var got struct {
		Input  map[string]any `json:"input"`
		Output map[string]any `json:"output"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(got.Input) != 0 {
		t.Fatalf("input = %#v, want empty map", got.Input)
	}
	if got.Output["results"] != float64(0) {
		t.Fatalf("output.results = %#v, want 0", got.Output["results"])
	}
}
