package cligen

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"time"
)

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func newIO(cmd *cobra.Command) IO {
	return IO{In: cmd.InOrStdin(), Out: cmd.OutOrStdout(), Err: cmd.ErrOrStderr()}
}

func rangeArgs(min, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min {
			return fmt.Errorf("requires at least %d arg(s), only received %d", min, len(args))
		}
		if max >= 0 && len(args) > max {
			return fmt.Errorf("accepts at most %d arg(s), received %d", max, len(args))
		}
		return nil
	}
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("opencli: invalid duration default " + s + ": " + err.Error())
	}
	return d
}
