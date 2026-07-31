package switcher

import (
	"context"
	"fmt"
	"sort"

	"charm.land/huh/v2"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

// noContextsError is the single canonical message for "no contexts
// configured", shared with remove so the two commands no longer disagree on
// wording for the same condition.
const noContextsError = "no auth contexts found; run sd auth login first"

func New(store *config.Store) cligen.AuthSwitchHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.AuthSwitchParams) error {
		cfg, err := store.Load()
		if err != nil {
			return err
		}

		names := contextNames(cfg)
		if len(names) == 0 {
			return fmt.Errorf(noContextsError)
		}

		selected := ""
		if p.ContextName != "" {
			selected = p.ContextName
			if _, ok := cfg.Contexts[selected]; !ok {
				return fmt.Errorf("unknown context %q", selected)
			}
		} else {
			selected = cfg.CurrentContext
			if err := selectContext(ctx, io, cfg, names, &selected); err != nil {
				return err
			}
		}

		cfg.SetCurrentContext(selected)
		if err := store.Save(cfg); err != nil {
			return err
		}

		fmt.Fprintf(io.Out, "%s %s\n", tui.Accent.Render("Current context:"), selected)

		return nil
	}
}

func selectContext(ctx context.Context, io cligen.IO, cfg *config.Config, names []string, selected *string) error {
	options := make([]huh.Option[string], 0, len(names))
	for _, name := range names {
		label := name
		if name == cfg.CurrentContext {
			label += " " + tui.Muted.Render("(current)")
		}

		if entry, ok := cfg.Contexts[name]; ok && entry.APIURL != "" {
			label += tui.Muted.Render("  " + entry.APIURL)
		}

		options = append(options, huh.NewOption(label, name))
	}

	return tui.NewForm(
		io.In,
		io.Err,
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(tui.Title.Render("Choose a Storyden context")).
				Options(options...).
				Value(selected),
		),
	).RunWithContext(ctx)
}

func contextNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Contexts))
	for name := range cfg.Contexts {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
