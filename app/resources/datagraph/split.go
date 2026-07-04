package datagraph

import (
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// rough upper bound sentence size for most languages.
const roughMaxSentenceSize = 350

func (c Content) Split() []string {
	if c.IsEmpty() {
		return []string{}
	}

	state := chunkState{max: roughMaxSentenceSize}
	walkChunks(c.html, &state)

	return state.chunks
}

type chunkState struct {
	chunks         []string
	max            int
	heading        string
	pendingHeading bool
}

func (s *chunkState) addChunk(text string, preserveNewlines bool) {
	normalized := normalizeChunkText(text, preserveNewlines)
	if normalized == "" {
		return
	}
	if isNoiseChunk(normalized) {
		return
	}

	if s.pendingHeading && s.heading != "" {
		normalized = s.heading + "\n" + normalized
		s.pendingHeading = false
	}

	s.chunks = append(s.chunks, splitChunk(normalized, s.max, preserveNewlines)...)
}

func isNoiseChunk(s string) bool {
	runes := []rune(strings.TrimSpace(s))
	if len(runes) == 0 {
		return true
	}
	if len(runes) <= 24 && strings.ContainsRune(s, '©') {
		return true
	}
	return false
}

func walkChunks(n *html.Node, state *chunkState) {
	if n == nil || shouldIgnoreSubtree(n) {
		return
	}

	if n.Type == html.TextNode {
		// keep raw text that appears in block/container elements even without
		// paragraph tags.
		if hasIgnoredAncestor(n) {
			return
		}
		if isRawTextContainer(n.Parent) {
			state.addChunk(n.Data, true)
		}
		return
	}

	if n.Type == html.ElementNode {
		switch n.DataAtom {
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			heading := normalizeChunkText(textFromNode(n, false), false)
			if heading != "" {
				state.heading = heading
				state.pendingHeading = true
			}
			return
		case atom.P, atom.Blockquote, atom.Li:
			state.addChunk(textFromNode(n, false), false)
			return
		case atom.Pre:
			state.addChunk(textFromNode(n, true), true)
			return
		case atom.Tr:
			state.addChunk(tableRowToText(n), false)
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkChunks(c, state)
	}
}

func shouldIgnoreSubtree(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}

	return isIgnoredTag(n)
}

func isRawTextContainer(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}

	switch n.DataAtom {
	case atom.Body, atom.Main, atom.Article, atom.Section, atom.Div:
		return true
	default:
		return false
	}
}

func hasIgnoredAncestor(n *html.Node) bool {
	for curr := n.Parent; curr != nil; curr = curr.Parent {
		if curr.Type == html.ElementNode && isIgnoredTag(curr) {
			return true
		}
	}
	return false
}

func isIgnoredTag(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}

	switch n.DataAtom {
	case atom.Nav, atom.Footer, atom.Script, atom.Style, atom.Noscript:
		return true
	}

	switch strings.ToLower(n.Data) {
	case "nav", "footer", "script", "style", "noscript":
		return true
	default:
		return false
	}
}

func tableRowToText(n *html.Node) string {
	cells := []string{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.DataAtom == atom.Th || c.DataAtom == atom.Td) {
			cell := normalizeChunkText(textFromNode(c, false), false)
			if cell != "" {
				cells = append(cells, cell)
			}
		}
	}

	return strings.Join(cells, " | ")
}

func normalizeChunkText(s string, preserveNewlines bool) string {
	if preserveNewlines {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\r", "\n")
		return strings.TrimSpace(s)
	}

	return strings.TrimSpace(spaces.ReplaceAllString(s, " "))
}

func splitChunk(in string, max int, preserveNewlines bool) []string {
	if in == "" {
		return nil
	}

	runes := []rune(in)
	if len(runes) <= max {
		return []string{in}
	}

	var chunks []string
	for len(runes) > 0 {
		if len(runes) <= max {
			chunk := strings.TrimSpace(string(runes))
			if chunk != "" {
				chunks = append(chunks, chunk)
			}
			break
		}

		upper := max - 1
		lower := upper / 2
		boundary := -1
		spaceFallback := -1

		for i := upper; i > lower; i-- {
			switch runes[i] {
			case '.', ';', '!', '?', '\n':
				boundary = i
				i = -1
			case ' ':
				if spaceFallback == -1 {
					spaceFallback = i
				}
			}
		}

		if boundary == -1 {
			if spaceFallback != -1 {
				boundary = spaceFallback
			} else {
				boundary = upper
			}
		}

		left := strings.TrimSpace(string(runes[:boundary+1]))
		if left != "" {
			chunks = append(chunks, left)
		}
		runes = []rune(strings.TrimSpace(string(runes[boundary+1:])))
	}

	if !preserveNewlines {
		for i := range chunks {
			chunks[i] = normalizeChunkText(chunks[i], false)
		}
	}

	return chunks
}

func needsSpace(left rune, right rune) bool {
	if left == 0 || right == 0 {
		return false
	}
	if unicode.IsSpace(left) || unicode.IsSpace(right) {
		return false
	}
	if isNoSpaceScript(left) && isNoSpaceScript(right) {
		return false
	}
	if strings.ContainsRune("([{\"'`", right) {
		return true
	}
	if strings.ContainsRune(")]},.!?:;\"'`", right) {
		return false
	}
	if strings.ContainsRune("([{\"'`", left) {
		return false
	}
	return true
}

func isNoSpaceScript(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul)
}
