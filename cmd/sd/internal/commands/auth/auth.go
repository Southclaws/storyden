package auth

import (
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/login"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/remove"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/switcher"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type AuthCommand *cobra.Command

func New(
	loginCommand login.LoginCommand,
	removeCommand remove.RemoveCommand,
	switchCommand switcher.SwitchCommand,
) AuthCommand {
	command := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Storyden instances",
	}

	command.AddCommand((loginCommand))
	command.AddCommand((removeCommand))
	command.AddCommand((switchCommand))

	help.SetupMarkdownHelp(command)

	return AuthCommand(command)
}
