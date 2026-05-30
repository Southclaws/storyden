package remove

import (
	"fmt"
	"sort"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

type RemoveCommand *cobra.Command

func New(store *config.Store) RemoveCommand {
	command := &cobra.Command{
		Use:   "remove [context]",
		Short: "Remove a Storyden auth context",
		Long: `# Remove Authentication Context

Remove a saved Storyden authentication context.

Run without arguments to choose a context interactively:
~~~bash
sd auth remove
~~~

Remove a context directly:
~~~bash
sd auth remove localhost-8000
~~~
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := store.Load()
			if err != nil {
				return err
			}

			names := contextNames(cfg)
			if len(names) == 0 {
				return fmt.Errorf("no auth contexts found")
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

			if err := removeContext(store, cfg, selected); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", tui.Accent.Render("Removed context:"), selected)

			if cfg.CurrentContext != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", tui.Accent.Render("Current context:"), cfg.CurrentContext)
			}

			return nil
		},
	}

	help.SetupMarkdownHelp(command)

	return RemoveCommand(command)
}

func removeContext(store *config.Store, cfg *config.Config, name string) error {
	delete(cfg.Contexts, name)

	if cfg.CurrentContext == name {
		cfg.CurrentContext = nextCurrentContext(cfg)
	}

	if err := store.DeleteAuth(name); err != nil {
		return err
	}

	return store.Save(cfg)
}

func nextCurrentContext(cfg *config.Config) string {
	names := contextNames(cfg)
	if len(names) == 0 {
		return ""
	}

	return names[0]
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
				Title(tui.Title.Render("Choose a Storyden context to remove")).
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
