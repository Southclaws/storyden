package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteUsesFangForCommandErrors(t *testing.T) {
	var stderr bytes.Buffer

	root := &cobra.Command{
		Use:           "sd",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetErr(&stderr)
	root.SetArgs([]string{"missing"})

	err := Execute(context.Background(), root)
	if err == nil {
		t.Fatal("expected command error")
	}
	if !IsCommandError(err) {
		t.Fatalf("expected command error wrapper, got %T", err)
	}

	output := stderr.String()
	if !strings.Contains(strings.ToLower(output), "unknown command") {
		t.Fatalf("expected fang error output, got %q", output)
	}
}

func TestIsHelpRequest(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "no args", args: nil, want: true},
		{name: "long help flag", args: []string{"node", "update", "--help"}, want: true},
		{name: "short help flag", args: []string{"node", "-h"}, want: true},
		{name: "help command", args: []string{"help", "node"}, want: true},
		{name: "normal command", args: []string{"node", "list"}, want: false},
		{name: "carapace callback", args: []string{"_carapace", "bash"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHelpRequest(tt.args); got != tt.want {
				t.Fatalf("isHelpRequest(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
