package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func NewStreams() Streams {
	return Streams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

type CommandError struct {
	Err error
}

func (e CommandError) Error() string {
	return e.Err.Error()
}

func (e CommandError) Unwrap() error {
	return e.Err
}

func IsCommandError(err error) bool {
	var commandErr CommandError
	return errors.As(err, &commandErr)
}

func Execute(ctx context.Context, root *cobra.Command) error {
	args := os.Args[1:]
	if isCarapaceRequest(args) {
		root.SetArgs(normalizeCarapaceArgs(root, args))
		return root.ExecuteContext(ctx)
	}

	if isHelpRequest(args) {
		return root.ExecuteContext(ctx)
	}

	if err := fang.Execute(ctx, root); err != nil {
		return CommandError{Err: err}
	}

	return nil
}

func isRawCobraRequest(args []string) bool {
	return isHelpRequest(args) || isCarapaceRequest(args)
}

func normalizeCarapaceArgs(root *cobra.Command, args []string) []string {
	normalized := append([]string(nil), args...)
	if len(normalized) == 3 && normalized[2] == root.Name() {
		normalized = append(normalized, "")
	}
	return normalized
}

func isHelpRequest(args []string) bool {
	if len(args) == 0 {
		return true
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" || strings.EqualFold(arg, "help") {
			return true
		}
	}

	return false
}

func isCarapaceRequest(args []string) bool {
	return len(args) > 0 && args[0] == "_carapace"
}
