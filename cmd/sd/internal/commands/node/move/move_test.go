package move

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestMoveCommandDoesNotAcceptIndex(t *testing.T) {
	r := require.New(t)

	command := (*cobra.Command)(New(nil))
	command.SetArgs([]string{"node-slug", "--index", "1"})

	err := command.Execute()
	r.ErrorContains(err, "unknown flag: --index")
}
