package switcher

import (
	"fmt"
	"sort"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

type SwitchCommand *cobra.Command

func New(store *config.Store) SwitchCommand {
	command := &cobra.Command{
		Use:   "switch [context]",
		Short: "Switch the active Storyden auth context",
		Long: `# Switch Authentication Context

Switch between different authenticated Storyden instances.

Contexts are auto-generated from hostnames (e.g., ` + "`community.com`" + ` → ` + "`community-com`" + `).

## Examples

Switch interactively:
~~~bash
sd auth switch
~~~

Switch directly:
~~~bash
sd auth switch my-community-com
~~~

Workflow with multiple instances:
~~~bash
sd auth login https://community1.com
sd auth login https://community2.com

sd auth switch community1-com
sd node list  # From community1

sd auth switch community2-com
sd node list  # From community2
~~~
`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := store.Load()
			if err != nil {
				return err
			}

			names := contextNames(cfg)
			if len(names) == 0 {
				return fmt.Errorf("no auth contexts found; run sd auth login first")
			}

			selected := ""
			if len(args) == 1 {
				selected = args[0]
				if _, ok := cfg.Contexts[selected]; !ok {
					return fmt.Errorf("unknown context %q", selected)
				}
			} else {
				selected = cfg.CurrentContext
				if err := selectContext(cmd, cfg, names, &selected); err != nil {
					return err
				}
			}

			cfg.SetCurrentContext(selected)
			if err := store.Save(cfg); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", tui.Accent.Render("Current context:"), selected)

			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return SwitchCommand(command)
}

func selectContext(cmd *cobra.Command, cfg *config.Config, names []string, selected *string) error {
	options := make([]huh.Option[string], 0, len(names))
	for _, name := range names {
		label := name
		if name == cfg.CurrentContext {
			label += " " + tui.Muted.Render("(current)")
		}

		if ctx, ok := cfg.Contexts[name]; ok && ctx.APIURL != "" {
			label += tui.Muted.Render("  " + ctx.APIURL)
		}

		options = append(options, huh.NewOption(label, name))
	}

	return tui.NewForm(
		cmd.InOrStdin(),
		cmd.ErrOrStderr(),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(tui.Title.Render("Choose a Storyden context")).
				Options(options...).
				Value(selected),
		),
	).RunWithContext(cmd.Context())
}

func contextNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Contexts))
	for name := range cfg.Contexts {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
