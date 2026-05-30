package filter

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NodeFlags is the cobra-bindable companion to NodeOptions. Embed in a command
// and call Bind during setup. Build returns the NodeOptions for use at run
// time, after running validation.
type NodeFlags struct {
	LinkDomains     []string
	LinkURLContains string
	LinkScheme      string
	NoLink          bool
	HasLink         bool
	RootOnly        bool
	OwnerHandle     string
	NameContains    string
}

// Bind registers all client-side filter flags on cmd.
func (f *NodeFlags) Bind(cmd *cobra.Command) {
	cmd.Flags().StringSliceVar(&f.LinkDomains, "link-domain", nil, "Filter to nodes whose link.domain matches; repeatable")
	cmd.Flags().StringVar(&f.LinkURLContains, "link-url-contains", "", "Filter to nodes whose link.url contains this substring")
	cmd.Flags().StringVar(&f.LinkScheme, "link-scheme", "", "Filter to nodes whose link.url uses this scheme (http, https)")
	cmd.Flags().BoolVar(&f.NoLink, "no-link", false, "Filter to nodes that have no link attached")
	cmd.Flags().BoolVar(&f.HasLink, "has-link", false, "Filter to nodes that have a link attached")
	cmd.Flags().BoolVar(&f.RootOnly, "root-only", false, "Filter to nodes with no parent")
	cmd.Flags().StringVar(&f.OwnerHandle, "owner-handle", "", "Filter to nodes owned by this account handle")
	cmd.Flags().StringVar(&f.NameContains, "name-contains", "", "Filter to nodes whose name contains this substring (case-insensitive)")
}

// Validate returns an error if mutually exclusive flags are set together.
func (f *NodeFlags) Validate() error {
	if f.NoLink && f.HasLink {
		return fmt.Errorf("--no-link and --has-link are mutually exclusive")
	}
	if f.LinkScheme != "" && f.LinkScheme != "http" && f.LinkScheme != "https" {
		return fmt.Errorf("--link-scheme must be http or https")
	}
	return nil
}

// Build assembles a NodeOptions value. Call after Validate.
func (f *NodeFlags) Build() NodeOptions {
	return NodeOptions{
		LinkDomains:     f.LinkDomains,
		LinkURLContains: f.LinkURLContains,
		LinkScheme:      f.LinkScheme,
		NoLink:          f.NoLink,
		HasLink:         f.HasLink,
		RootOnly:        f.RootOnly,
		OwnerHandle:     f.OwnerHandle,
		NameContains:    f.NameContains,
	}
}
