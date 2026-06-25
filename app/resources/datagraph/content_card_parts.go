package datagraph

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// SplitCardParts splits content at paragraph boundaries that contain only a
// single sdr: anchor. Contiguous non-boundary nodes are grouped into one part;
// each sdr-only paragraph becomes its own isolated part.
//
// Parts are built by deep-cloning the relevant nodes from the already-parsed
// tree, avoiding any re-parsing or re-sanitization overhead.
func (c Content) SplitCardParts() []Content {
	if c.html == nil {
		return nil
	}

	baseURL, _ := url.Parse("ignore:")
	var parts []Content
	var buffer []*html.Node

	flush := func() {
		if !hasContentNodes(buffer) {
			buffer = nil
			return
		}
		parts = append(parts, buildPartFromNodes(buffer, baseURL))
		buffer = nil
	}

	for child := c.html.FirstChild; child != nil; child = child.NextSibling {
		if isSdrOnlyParagraph(child) {
			flush()
			parts = append(parts, buildPartFromNodes([]*html.Node{child}, baseURL))
		} else {
			buffer = append(buffer, child)
		}
	}
	flush()

	return parts
}

func buildPartFromNodes(nodes []*html.Node, baseURL *url.URL) Content {
	body := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}

	var last *html.Node
	for _, n := range nodes {
		cloned := cloneNodeDeep(n)
		cloned.Parent = body
		if last != nil {
			last.NextSibling = cloned
			cloned.PrevSibling = last
		} else {
			body.FirstChild = cloned
		}
		last = cloned
	}
	body.LastChild = last

	bodyTree, links, media, refs := extractReferences(body, baseURL)

	plain := strings.TrimSpace(spaces.ReplaceAllString(textFromNode(bodyTree, false), " "))

	return Content{
		html:  bodyTree,
		short: shortFromPlain(plain),
		plain: plain,
		links: links,
		media: media,
		sdrs:  refs,
	}
}

func cloneNodeDeep(n *html.Node) *html.Node {
	clone := &html.Node{
		Type:      n.Type,
		DataAtom:  n.DataAtom,
		Data:      n.Data,
		Namespace: n.Namespace,
		Attr:      make([]html.Attribute, len(n.Attr)),
	}
	copy(clone.Attr, n.Attr)

	var last *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := cloneNodeDeep(c)
		child.Parent = clone
		if last != nil {
			last.NextSibling = child
			child.PrevSibling = last
		} else {
			clone.FirstChild = child
		}
		last = child
	}
	clone.LastChild = last

	return clone
}

func shortFromPlain(plain string) string {
	paragraphs := []rune(plain)
	end := int(math.Min(float64(len(paragraphs)-1), MaxSummaryLength))

	if len(paragraphs) <= MaxSummaryLength {
		return plain
	}

	for ; end > MaxSummaryLength/2; end-- {
		if unicode.IsPunct(paragraphs[end]) || unicode.IsSpace(paragraphs[end]) {
			break
		}
	}

	if !unicode.IsLetter(paragraphs[end-1]) {
		for ; end > MaxSummaryLength/2; end-- {
			if unicode.IsLetter(paragraphs[end]) {
				end++
				break
			}
		}
	}

	return fmt.Sprint(string(paragraphs[:end]), "...")
}

func hasContentNodes(nodes []*html.Node) bool {
	for _, n := range nodes {
		if n.Type == html.ElementNode {
			return true
		}
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			return true
		}
	}
	return false
}

func isSdrOnlyParagraph(n *html.Node) bool {
	if n.Type != html.ElementNode || n.DataAtom != atom.P {
		return false
	}

	var anchor *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && strings.TrimSpace(c.Data) == "" {
			continue
		}
		if anchor != nil {
			return false
		}
		anchor = c
	}

	if anchor == nil || anchor.Type != html.ElementNode || anchor.DataAtom != atom.A {
		return false
	}

	for _, attr := range anchor.Attr {
		if strings.ToLower(attr.Key) == "href" {
			u, err := url.Parse(attr.Val)
			if err != nil {
				return false
			}
			return u.Scheme == RefScheme
		}
	}

	return false
}
