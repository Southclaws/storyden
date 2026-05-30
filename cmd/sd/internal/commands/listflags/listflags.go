// Package listflags provides a reusable flag set for list-style commands so
// node list, thread list, and node children share the same surface: --page,
// --limit, --all, --format, --output.
package listflags

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const (
	FormatAuto  = "auto"
	FormatPlain = "plain"
	FormatJSON  = "json"
	FormatJSONL = "jsonl"

	OutputDefault = "default"
	OutputWide    = "wide"
)

// Flags is the shared list flag set. Embed it in a command's local state and
// call Bind during command setup.
type Flags struct {
	Page   int
	Limit  int
	All    bool
	Format string
	Output string
}

// Bind registers the shared flags on cmd with sensible defaults.
func (f *Flags) Bind(cmd *cobra.Command) {
	cmd.Flags().IntVar(&f.Page, "page", 1, "Page to request")
	cmd.Flags().IntVar(&f.Limit, "limit", 0, "Stop after N matches (0 = no limit)")
	cmd.Flags().BoolVar(&f.All, "all", false, "Fetch every page, streaming output as it goes")
	cmd.Flags().StringVar(&f.Format, "format", FormatAuto, "Output format: auto, plain, json, jsonl")
	cmd.Flags().StringVarP(&f.Output, "output", "o", OutputDefault, "Column profile: default, wide")
}

// Validate checks the user-supplied flag values up-front so commands can fail
// fast before making API calls.
func (f *Flags) Validate() error {
	if f.Page < 1 {
		return fmt.Errorf("--page must be greater than zero")
	}
	if f.Limit < 0 {
		return fmt.Errorf("--limit must be zero or positive")
	}
	switch f.Format {
	case FormatAuto, FormatPlain, FormatJSON, FormatJSONL:
	default:
		return fmt.Errorf("--format must be one of: auto, plain, json, jsonl")
	}
	switch f.Output {
	case OutputDefault, OutputWide:
	default:
		return fmt.Errorf("--output must be one of: default, wide")
	}
	return nil
}

// ResolveFormat reduces FormatAuto to a concrete format. Auto always resolves
// to plain — agents run via subprocess (non-TTY) and need a predictable shape
// unless they explicitly ask for JSON. The writer arg is accepted for future
// TTY-vs-pipe heuristics; today it's intentionally unused so behaviour is the
// same in every environment.
func (f *Flags) ResolveFormat(_ io.Writer) string {
	if f.Format == FormatAuto {
		return FormatPlain
	}
	return f.Format
}

// Wide reports whether the wide column profile is selected.
func (f *Flags) Wide() bool {
	return f.Output == OutputWide
}
