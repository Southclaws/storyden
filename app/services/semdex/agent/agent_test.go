package agent_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	agentpkg "google.golang.org/adk/agent"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/semdex/agent"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestAgent(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	cfg := &config.Config{
		LanguageModelProvider: "openai",
		OpenAIKey:             key,
	}

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		nodeManager *node_mutate.Manager,
		agent agent.Agent,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			acc, err := aw.Create(root, "test-semdex", account_writer.WithAdmin(true))
			r.NoError(err)

			content, err := datagraph.NewRichText("<body>This is a test node for the Storyden Semdex agent.</body>")
			r.NoError(err)

			// The node we want to find.
			node, err := nodeManager.Create(root, acc.ID, "Test Node", node_mutate.Partial{
				Content:    opt.New(content),
				Name:       opt.New("Test Node"),
				Visibility: opt.New(visibility.VisibilityPublished),
			})
			r.NoError(err)

			prompt := fmt.Sprintf("Call the getLibraryPage tool with slug '%s' and describe the page.", node.Mark.Slug())
			input := &genai.Content{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{Text: prompt},
				},
			}

			stream := agent.Run(root, acc.ID.String(), "agent-smoke-session", input, agentpkg.RunConfig{})
			for event, err := range stream {
				if err != nil {
					if strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "lookup") {
						t.Skipf("network unavailable: %v", err)
					}
					t.Fatalf("agent stream error: %v", err)
				}
				if event == nil {
					continue
				}

				pp.Println(event)
			}
		}))
	}))
}
