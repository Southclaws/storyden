package token

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token/rotate"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type TokenCommand *cobra.Command

func New(
	rotateCommand rotate.RotateCommand,
) TokenCommand {
	command := &cobra.Command{
		Use:   "token",
		Short: "Manage external plugin RPC tokens",
		Long: `# Plugin Tokens

Manage static RPC tokens for external plugin installations.

Use token rotation only when replacing an exposed credential or performing incident response.
`,
	}

	command.AddCommand(rotateCommand)

	help.SetupMarkdownHelp(command)

	return TokenCommand(command)
}
