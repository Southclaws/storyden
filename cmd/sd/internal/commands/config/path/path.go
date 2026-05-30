package path

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type PathCommand *cobra.Command

func New(store *config.Store) PathCommand {
	command := &cobra.Command{
		Use:   "path",
		Short: "Print the sd config file path",
		Long: `# Show Config File Path

Print the path to your sd configuration file (stores auth contexts and tokens).

## Examples

View path:
~~~bash
sd config path
~~~

View config contents:
~~~bash
cat $(sd config path)
~~~

Edit config:
~~~bash
code $(sd config path)
~~~

Backup config:
~~~bash
cp $(sd config path) ~/backup.yaml
~~~
`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), store.Path())
			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return PathCommand(command)
}
