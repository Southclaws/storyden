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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHelpRequest(tt.args); got != tt.want {
				t.Fatalf("isHelpRequest(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestIsRawCobraRequest(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "no args", args: nil, want: true},
		{name: "help", args: []string{"node", "--help"}, want: true},
		{name: "carapace script", args: []string{"_carapace", "nushell"}, want: true},
		{name: "carapace callback", args: []string{"_carapace", "nushell", "sd", "node"}, want: true},
		{name: "normal command", args: []string{"node", "list"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRawCobraRequest(tt.args); got != tt.want {
				t.Fatalf("isRawCobraRequest(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestNormalizeCarapaceArgs(t *testing.T) {
	root := &cobra.Command{Use: "sd"}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "root command callback gets empty current token",
			args: []string{"_carapace", "nushell", "sd"},
			want: []string{"_carapace", "nushell", "sd", ""},
		},
		{
			name: "script generation is unchanged",
			args: []string{"_carapace", "nushell"},
			want: []string{"_carapace", "nushell"},
		},
		{
			name: "partial command callback is unchanged",
			args: []string{"_carapace", "nushell", "sd", "n"},
			want: []string{"_carapace", "nushell", "sd", "n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeCarapaceArgs(root, tt.args)
			if strings.Join(got, "\x00") != strings.Join(tt.want, "\x00") {
				t.Fatalf("normalizeCarapaceArgs(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
