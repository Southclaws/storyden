package filter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func node(name, ownerHandle string, link *openapi.LinkReference, parent *openapi.Node) openapi.NodeWithChildren {
	return openapi.NodeWithChildren{
		Name:   openapi.NodeName(name),
		Owner:  openapi.ProfileReference{Handle: openapi.AccountHandle(ownerHandle)},
		Link:   link,
		Parent: parent,
	}
}

func TestMatchNodeEmptyOptionsAllPass(t *testing.T) {
	r := require.New(t)
	r.True(MatchNode(node("a", "b", nil, nil), NodeOptions{}))
}

func TestMatchNodeRootOnly(t *testing.T) {
	r := require.New(t)
	root := node("root", "", nil, nil)
	child := node("child", "", nil, &openapi.Node{})
	r.True(MatchNode(root, NodeOptions{RootOnly: true}))
	r.False(MatchNode(child, NodeOptions{RootOnly: true}))
}

func TestMatchNodeNoLinkVsHasLink(t *testing.T) {
	r := require.New(t)
	withLink := node("a", "", &openapi.LinkReference{Domain: "example.com", Url: "https://example.com"}, nil)
	noLink := node("b", "", nil, nil)

	r.True(MatchNode(withLink, NodeOptions{HasLink: true}))
	r.False(MatchNode(noLink, NodeOptions{HasLink: true}))
	r.False(MatchNode(withLink, NodeOptions{NoLink: true}))
	r.True(MatchNode(noLink, NodeOptions{NoLink: true}))
}

func TestMatchNodeLinkDomainAnyMatch(t *testing.T) {
	r := require.New(t)
	a := node("a", "", &openapi.LinkReference{Domain: "media.tenor.com"}, nil)
	b := node("b", "", &openapi.LinkReference{Domain: "youtube.com"}, nil)
	c := node("c", "", &openapi.LinkReference{Domain: "example.com"}, nil)

	opts := NodeOptions{LinkDomains: []string{"tenor.com", "youtu.be"}}
	r.True(MatchNode(a, opts))
	r.False(MatchNode(b, opts))
	r.False(MatchNode(c, opts))
}

func TestMatchNodeOwnerHandle(t *testing.T) {
	r := require.New(t)
	a := node("a", "alice", nil, nil)
	b := node("b", "bob", nil, nil)
	r.True(MatchNode(a, NodeOptions{OwnerHandle: "alice"}))
	r.False(MatchNode(b, NodeOptions{OwnerHandle: "alice"}))
}

func TestMatchNodeURLContainsAndScheme(t *testing.T) {
	r := require.New(t)
	a := node("a", "", &openapi.LinkReference{Domain: "youtube.com", Url: "https://youtube.com/watch?v=1"}, nil)
	b := node("b", "", &openapi.LinkReference{Domain: "youtube.com", Url: "http://youtube.com/x"}, nil)
	r.True(MatchNode(a, NodeOptions{LinkURLContains: "/watch"}))
	r.False(MatchNode(b, NodeOptions{LinkURLContains: "/watch"}))
	r.True(MatchNode(a, NodeOptions{LinkScheme: "https"}))
	r.False(MatchNode(b, NodeOptions{LinkScheme: "https"}))
}

func TestFilterNodesPreservesOrder(t *testing.T) {
	r := require.New(t)
	nodes := []openapi.NodeWithChildren{
		node("a", "alice", nil, nil),
		node("b", "bob", nil, nil),
		node("c", "alice", nil, nil),
	}
	out := FilterNodes(nodes, NodeOptions{OwnerHandle: "alice"})
	r.Len(out, 2)
	r.Equal("a", string(out[0].Name))
	r.Equal("c", string(out[1].Name))
}
