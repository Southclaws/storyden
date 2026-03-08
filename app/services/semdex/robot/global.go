package robot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/adk/agent"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/mcp"
)

func (s *Agent) globalInstructionProvider(ctx context.Context, chatContext *mcp.RobotChatContext) func(ctx agent.ReadonlyContext) (string, error) {
	return func(ctx agent.ReadonlyContext) (string, error) {
		acc, err := session.GetAccount(ctx)
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}
		roles := session.GetRoles(ctx)

		var b strings.Builder

		// from global.md file embed - this is a hard coded Storyden bootstrap.
		b.WriteString(globalInstruction)

		b.WriteString("\n\n## Current Context\n\n")
		b.WriteString(fmt.Sprintf("Current date and time: %s\n\n", time.Now().UTC().Format(time.RFC3339)))

		b.WriteString(fmt.Sprintf("The user is: %s who is %s\n\n", acc.Name, acc.Kind.String()))

		// Add user permissions
		permissions := roles.Permissions().List()
		if len(permissions) > 0 {
			b.WriteString("### Permissions\n\n")
			b.WriteString("The current user has the following permissions:\n")
			for _, perm := range permissions {
				b.WriteString(fmt.Sprintf("- %s\n", perm.String()))
			}
			b.WriteString("\nOnly provide functionality and suggestions that align with these permissions. Do not suggest actions the user cannot perform.\n\n")
		}

		if chatContext != nil {
			if chatContext.DatagraphItem != nil {
				item := chatContext.DatagraphItem
				b.WriteString("The user is currently viewing:\n")
				b.WriteString(fmt.Sprintf("- Type: %s\n", item.Kind))
				b.WriteString(fmt.Sprintf("- ID: %s\n", item.Id))
				b.WriteString(fmt.Sprintf("- Slug: %s\n", item.Slug))
				b.WriteString("\nThis is the primary context of the conversation. When the user refers to \"this\", \"here\", or similar demonstratives, they likely mean this item.\n")
			} else if chatContext.PageType != nil && *chatContext.PageType != "" {
				b.WriteString(fmt.Sprintf("The user is currently on: %s\n", *chatContext.PageType))
			}
		}

		return b.String(), nil
	}
}
