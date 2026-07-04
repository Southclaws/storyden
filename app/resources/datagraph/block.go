package datagraph

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Block is an addressable structural unit inside a Content document.
type Block struct {
	ID   string
	Type string
	Text string
	HTML string
}

const (
	blockIDAttributeName = "id"
	blockIDPrefix        = "sdb_"
)

// blockAtoms is the set of HTML element types that receive stable block IDs.
var blockAtoms = map[atom.Atom]bool{
	atom.H1: true, atom.H2: true, atom.H3: true,
	atom.H4: true, atom.H5: true, atom.H6: true,
	atom.P: true, atom.Blockquote: true, atom.Li: true,
	atom.Pre: true, atom.Figure: true, atom.Table: true,
	atom.Img: true,
	atom.Div: true, atom.Section: true, atom.Article: true,
	atom.Ul: true, atom.Ol: true,
}

func isBlockAtom(a atom.Atom) bool {
	return blockAtoms[a]
}

func isBlockNode(n *html.Node) bool {
	return n.Type == html.ElementNode && isBlockAtom(n.DataAtom)
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func setAttr(n *html.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}

func setBlockID(n *html.Node, id string) {
	setAttr(n, blockIDAttributeName, id)
}

func normaliseIDAttributes(n *html.Node, preserveBlockIDs bool) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		attrs := n.Attr[:0]
		for _, attr := range n.Attr {
			if attr.Key != blockIDAttributeName {
				attrs = append(attrs, attr)
				continue
			}
			if preserveBlockIDs && isValidBlockID(attr.Val) {
				attrs = append(attrs, attr)
			}
		}
		n.Attr = attrs
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		normaliseIDAttributes(c, preserveBlockIDs)
	}
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	_ = html.Render(&buf, n)
	return buf.String()
}

// renderNodeWithoutBlockID renders the node's outer HTML with Storyden block
// IDs removed recursively from the subtree. This gives container blocks
// (blockquote, table, li with children) a stable normHTML that is independent
// of whether their children have already been stamped with IDs.
func renderNodeWithoutBlockID(n *html.Node) string {
	type saved struct {
		node  *html.Node
		attrs []html.Attribute
	}
	var restore []saved

	var strip func(*html.Node)
	strip = func(cur *html.Node) {
		if cur.Type == html.ElementNode {
			originalAttrs := append([]html.Attribute(nil), cur.Attr...)
			newAttrs := make([]html.Attribute, 0, len(cur.Attr))
			changed := false
			for _, a := range cur.Attr {
				if a.Key == blockIDAttributeName && isValidBlockID(a.Val) {
					changed = true
					continue
				}
				newAttrs = append(newAttrs, a)
			}
			if changed {
				restore = append(restore, saved{cur, originalAttrs})
				cur.Attr = newAttrs
			}
		}
		for c := cur.FirstChild; c != nil; c = c.NextSibling {
			strip(c)
		}
	}
	strip(n)

	s := renderNode(n)

	for _, item := range restore {
		item.node.Attr = item.attrs
	}
	return s
}

func findBody(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}
	if n.DataAtom == atom.Body {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findBody(c); found != nil {
			return found
		}
	}
	return nil
}
