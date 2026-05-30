package filter

import (
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

// NodeOptions captures all client-side node filters. Zero values disable each
// filter so callers can supply a partial subset.
type NodeOptions struct {
	LinkDomains     []string // any-match
	LinkURLContains string
	LinkScheme      string // "http" or "https"
	NoLink          bool   // require node.Link == nil
	HasLink         bool   // require node.Link != nil
	RootOnly        bool   // require node.Parent == nil
	OwnerHandle     string
	NameContains    string
}

// Empty reports whether opts has no active filters; useful to skip the
// post-fetch loop entirely on the common case.
func (o NodeOptions) Empty() bool {
	return len(o.LinkDomains) == 0 &&
		o.LinkURLContains == "" &&
		o.LinkScheme == "" &&
		!o.NoLink &&
		!o.HasLink &&
		!o.RootOnly &&
		o.OwnerHandle == "" &&
		o.NameContains == ""
}

// MatchNode reports whether n passes every active predicate in opts. AND
// semantics across fields; OR within LinkDomains (any-domain match).
func MatchNode(n openapi.NodeWithChildren, opts NodeOptions) bool {
	if opts.RootOnly && n.Parent != nil {
		return false
	}
	if opts.NoLink && n.Link != nil {
		return false
	}
	if opts.HasLink && n.Link == nil {
		return false
	}
	if opts.OwnerHandle != "" && string(n.Owner.Handle) != opts.OwnerHandle {
		return false
	}
	if opts.NameContains != "" && !NameContains(string(n.Name), opts.NameContains) {
		return false
	}
	if len(opts.LinkDomains) > 0 {
		if n.Link == nil {
			return false
		}
		matched := false
		for _, d := range opts.LinkDomains {
			if DomainMatches(string(n.Link.Domain), d) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if opts.LinkURLContains != "" {
		if n.Link == nil || !URLContains(string(n.Link.Url), opts.LinkURLContains) {
			return false
		}
	}
	if opts.LinkScheme != "" {
		if n.Link == nil || URLScheme(string(n.Link.Url)) != opts.LinkScheme {
			return false
		}
	}
	return true
}

// FilterNodes returns nodes from in that pass MatchNode for opts.
func FilterNodes(in []openapi.NodeWithChildren, opts NodeOptions) []openapi.NodeWithChildren {
	if opts.Empty() {
		return in
	}
	out := in[:0:0]
	for _, n := range in {
		if MatchNode(n, opts) {
			out = append(out, n)
		}
	}
	return out
}
