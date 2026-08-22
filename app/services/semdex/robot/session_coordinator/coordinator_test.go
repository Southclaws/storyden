package session_coordinator

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"

	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
)

func TestShouldProjectUserInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source       robotservice.RunSource
		content      *genai.Content
		continuation bool
		expected     bool
	}{
		"interactive chat":  {source: robotservice.SourceInteractiveChat, content: genai.NewContentFromText("hello", genai.RoleUser), expected: true},
		"plugin RPC":        {source: robotservice.SourcePluginRPC, content: genai.NewContentFromText("hello", genai.RoleUser), expected: true},
		"scheduled Trail":   {source: robotservice.SourceScheduled, content: genai.NewContentFromText("post the weekly prompt", genai.RoleUser), expected: true},
		"delegation":        {source: robotservice.SourceDelegation, content: genai.NewContentFromText("delegated request", genai.RoleUser), expected: false},
		"delegation result": {source: robotservice.SourceDelegationResult, content: genai.NewContentFromText("delegation result", genai.RoleUser), expected: false},
		"model content":     {source: robotservice.SourceScheduled, content: genai.NewContentFromText("model output", genai.RoleModel), expected: false},
		"tool continuation": {source: robotservice.SourceScheduled, content: genai.NewContentFromText("tool result", genai.RoleUser), continuation: true, expected: false},
		"empty content":     {source: robotservice.SourceScheduled, expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			command := &CommandEnqueueMessage{
				InputID: xid.New(),
				Content: test.content,
				Options: robotservice.RunOptions{Source: test.source},
			}
			assert.Equal(t, test.expected, shouldProjectUserInput(command, test.continuation))
		})
	}
}
