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
	if isHelpRequest(os.Args[1:]) {
		return root.ExecuteContext(ctx)
	}

	if err := fang.Execute(ctx, root, fang.WithoutCompletions()); err != nil {
		return CommandError{Err: err}
	}

	return nil
}

func isHelpRequest(args []string) bool {
	if len(args) == 0 {
		return true
	}

	if args[0] == "_carapace" {
		return false
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" || strings.EqualFold(arg, "help") {
			return true
		}
	}

	return false
}
