package help

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestIsTerminalReturnsFalseForNonFileWriters(t *testing.T) {
	a := assert.New(t)

	a.False(IsTerminal(&bytes.Buffer{}))
}

func TestGeneratedHelpUsesSummaryAndExamples(t *testing.T) {
	var out bytes.Buffer
	command := &cobra.Command{Use: "sample", Short: "Inspect a sample.", Example: "sample --output json"}
	command.SetOut(&out)
	SetupMarkdownHelp(command)
	assert.NoError(t, command.Help())
	assert.Contains(t, out.String(), "Inspect a sample.")
	assert.Contains(t, out.String(), "## Examples\n")
	assert.Contains(t, out.String(), "sample --output json")
}

func TestPipedHelpIsOneMarkdownDocument(t *testing.T) {
	var out bytes.Buffer
	root := &cobra.Command{Use: "sample"}
	root.PersistentFlags().String("context", "", "Identity for this operation.")
	child := &cobra.Command{
		Use: "read ID", Short: "Read a resource.",
		Long:    "# Read a resource\n\nUse **JSON** to preserve fields.\n",
		Example: "sample read page --output json",
	}
	child.Flags().String("output", "plain", "Output format.")
	root.AddCommand(child)
	root.SetOut(&out)
	SetupMarkdownHelp(root)
	root.SetArgs([]string{"read", "--help"})
	assert.NoError(t, root.Execute())
	help := out.String()
	assert.Contains(t, help, "# Read a resource\n\nUse **JSON** to preserve fields.")
	assert.Contains(t, help, "## Examples\n\n~~~sh\nsample read page --output json\n~~~")
	assert.Contains(t, help, "## Usage\n\n~~~text\nsample read ID [flags]\n~~~")
	assert.Contains(t, help, "## Flags\n\n~~~text\n")
	assert.Contains(t, help, "## Global flags\n\n~~~text\n")
	assert.Equal(t, 1, strings.Count(help, "--context"))
	assert.NotContains(t, help, "\x1b[")
}
