// Package completion adds Carapace actions to the spec-generated command tree.
package completion

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

type command struct {
	Path, Operation  string
	Arguments, Flags []input
}

type input struct {
	Name, Type             string
	Choices                []string
	Variadic, SplitOnComma bool
}

// Setup registers actions without reading credentials or contacting an instance.
func Setup(root *cobra.Command, store *config.Store) {
	carapace.Gen(root)
	for _, def := range definitions() {
		cmd, _, err := root.Find(strings.Fields(def.Path))
		if err != nil {
			continue
		}
		flags := carapace.ActionMap{}
		for _, in := range def.Flags {
			action := inputAction(cmd, store, def, in)
			if in.SplitOnComma {
				action = action.UniqueList(",")
			}
			flags[in.Name] = action
		}
		carapace.Gen(cmd).FlagCompletion(flags)
		var positional []carapace.Action
		for _, in := range def.Arguments {
			action := inputAction(cmd, store, def, in)
			if in.Variadic {
				carapace.Gen(cmd).PositionalAnyCompletion(action.FilterArgs())
			} else {
				positional = append(positional, action)
			}
		}
		carapace.Gen(cmd).PositionalCompletion(positional...)
	}
}

//go:embed bridge.yaml
var bridgeSpec string

func New() cligen.CompletionHandler {
	return func(_ context.Context, cmd *cobra.Command, streams cligen.IO, p cligen.CompletionParams) error {
		if p.Shell == cligen.CompletionShellCarapace {
			_, err := fmt.Fprint(streams.Out, bridgeSpec)
			return err
		}
		snippet, err := carapace.Gen(cmd.Root()).Snippet(string(p.Shell))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(streams.Out, snippet)
		return err
	}
}

func inputAction(cmd *cobra.Command, store *config.Store, def command, in input) carapace.Action {
	if len(in.Choices) > 0 {
		return carapace.ActionValues(in.Choices...)
	}
	if in.Name == "FILE" {
		return carapace.ActionFiles()
	}
	name := strings.ToLower(strings.ReplaceAll(in.Name, "-", "_"))
	switch name {
	case "context", "context_name":
		return contextNames(store)
	case "operation":
		return operationNames()
	case "dir", "directory":
		return carapace.ActionDirectories()
	case "file", "content_file", "message_file":
		return carapace.Batch(carapace.ActionFiles(), carapace.ActionValuesDescribed("-", "Read from stdin")).ToA()
	case "output_file", "manifest":
		return carapace.ActionFiles()
	case "auth_storage":
		return carapace.ActionValues("auto", "file")
	}
	if in.Type == "file" || in.Type == "path" {
		return carapace.ActionFiles()
	}
	if in.Type == "boolean" {
		return carapace.ActionValues("true", "false")
	}
	if def.Path == "api request" {
		return apiInput(cmd, store, name)
	}
	if def.Path == "page visibility" && name == "values" {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			visibility := carapace.ActionValues("draft", "review", "published", "unlisted")
			if cmd.Flag("from-stdin").Value.String() == "true" {
				return visibility
			}
			pages := resourceAction(cmd, store, def, pageSource("slug"))
			if len(c.Args) > 0 {
				return carapace.Batch(pages, visibility).ToA()
			}
			return pages
		})
	}
	if name == "token" && strings.HasPrefix(def.Path, "page properties schema ") {
		return carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 1:
				return carapace.ActionValues("text", "number", "boolean", "timestamp").Suffix(":")
			case 2:
				return carapace.ActionValues("asc", "desc")
			}
			return carapace.ActionValues()
		})
	}
	if src, ok := resourceSource(def.Path, name); ok {
		return resourceAction(cmd, store, def, src)
	}
	return carapace.ActionValues()
}

func contextNames(store *config.Store) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		cfg, err := store.Load()
		if err != nil {
			return carapace.ActionValues()
		}
		var values []string
		for name := range cfg.Contexts {
			description := "Saved context"
			if name == cfg.CurrentContext {
				description = "Default context"
			}
			values = append(values, name, description)
		}
		return carapace.ActionValuesDescribed(values...)
	})
}
