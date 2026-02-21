package datagraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/Southclaws/fault"
	"github.com/cixtor/readability"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/samber/lo"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// RefScheme is used as a scheme for URIs to reference resources.
// These can be used in content to refer to profiles, posts, nodes, etc.
const RefScheme = "sdr"

var policy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()

	p.AllowURLSchemes(
		"mailto",
		"http",
		"https",
		RefScheme,
	)

	p.AllowDataAttributes()

	return p
}()

var spaces = regexp.MustCompile(`\s+`)

// MaxSummaryLength is the maximum length of the short summary text
const MaxSummaryLength = 128

// EmptyState is the default value for empty rich text. It could be an empty
// string but on the read path, it'll get turned into this either way.
const EmptyState = `<body></body>`

type Content struct {
	html  *html.Node
	short string
	plain string
	links []string
	media []string
	sdrs  RefList
}

func (c Content) MarshalJSON() ([]byte, error) {
	s := c.HTML()

	return json.Marshal(s)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	parsed, err := NewRichText(s)
	if err != nil {
		return err
	}

	*c = parsed

	return nil
}

func (r Content) HTML() string {
	if r.html == nil || r.html.FirstChild == nil {
		return EmptyState
	}

	w := &bytes.Buffer{}

	err := html.Render(w, r.html)
	if err != nil {
		panic(err)
	}

	return w.String()
}

func (r Content) HTMLTree() *html.Node {
	return r.html
}

func (r Content) Short() string {
	return r.short
}

func (r Content) Plaintext() string {
	return r.plain
}

func (r Content) Links() []string {
	return r.links
}

func (r Content) Media() []string {
	return r.media
}

func (r Content) References() RefList {
	return r.sdrs
}

func (r Content) IsEmpty() bool {
	return r.html == nil || r.plain == "" || r.short == ""
}

type options struct {
	baseURL string
}
type option func(*options)

func WithBaseURL(url string) option {
	return func(o *options) {
		o.baseURL = url
	}
}

// NewRichText will pull out any meaningful structured information from markdown
// document this includes a summary of the text and all link URLs for hydrating.
func NewRichText(raw string) (Content, error) {
	return NewRichTextFromReader(strings.NewReader(raw))
}

// NewRichText will pull out any meaningful structured information from markdown
// document this includes a summary of the text and all link URLs for hydrating.
func NewRichTextWithOptions(raw string, opts ...option) (Content, error) {
	return NewRichTextFromReader(strings.NewReader(raw), opts...)
}

func NewRichTextFromMarkdown(md string) (Content, error) {
	html := blackfriday.Run([]byte(md), blackfriday.WithExtensions(
		blackfriday.NoEmptyLineBeforeBlock,
	))

	return NewRichTextFromReader(strings.NewReader(string(html)))
}

func NewRichTextFromReader(r io.Reader, opts ...option) (Content, error) {
	o := options{baseURL: "ignore:"}
	for _, opt := range opts {
		opt(&o)
	}

	baseURL, err := url.Parse(o.baseURL)
	if err != nil {
		return Content{}, fault.Wrap(err)
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return Content{}, fault.Wrap(err)
	}

	sanitised := policy.SanitizeBytes(buf)

	htmlTree, err := html.Parse(bytes.NewReader(sanitised))
	if err != nil {
		return Content{}, fault.Wrap(err)
	}

	result, err := readability.New().Parse(bytes.NewReader(sanitised), o.baseURL)
	if err != nil {
		return Content{}, fault.Wrap(err)
	}

	short := getSummary(result)

	bodyTree, links, media, refs := extractReferences(htmlTree, baseURL)

	return Content{
		html:  bodyTree,
		short: short,
		plain: result.TextContent,
		links: links,
		media: media,
		sdrs:  refs,
	}, nil
}

func extractReferences(htmlTree *html.Node, baseURL *url.URL) (*html.Node, []string, []string, RefList) {
	bodyTree := &html.Node{}
	links := []string{}
	media := []string{}
	sdrs := []url.URL{}

	if htmlTree.DataAtom == atom.Body {
		bodyTree = htmlTree
	}

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Parent != nil {
			switch n.DataAtom {
			case atom.A:
				href, hasHref := lo.Find(n.Attr, func(a html.Attribute) bool {
					return strings.ToLower(a.Key) == "href"
				})

				if hasHref {
					if parsed, err := url.Parse(href.Val); err == nil {
						switch parsed.Scheme {
						case "http", "https":
							links = append(links, parsed.String())
						case RefScheme:
							sdrs = append(sdrs, *parsed)
						}
					}
				}

			case atom.Img:
				src, hasSrc := lo.Find(n.Attr, func(a html.Attribute) bool {
					return strings.ToLower(a.Key) == "src"
				})

				if hasSrc {
					if parsed, err := url.Parse(src.Val); err == nil {
						switch parsed.Scheme {
						case "":
							abs := baseURL.ResolveReference(parsed).String()
							media = append(media, abs)
						case "http", "https":
							media = append(media, parsed.String())
						case RefScheme:
							sdrs = append(sdrs, *parsed)
						}
					}
				}
			}
		}

		// NOTE: We don't need the <html>/<head> tags so skip them for storage.
		if n.DataAtom == atom.Body {
			bodyTree = n
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(htmlTree)

	var refs RefList
	for _, v := range sdrs {
		r, err := NewRefFromSDR(v)
		if err != nil {
			slog.Warn("invalid SDR in content", slog.String("error", err.Error()), slog.String("ref", v.Opaque))

			continue
		}
		refs = append(refs, r)
	}

	return bodyTree, links, media, refs
}

func getSummary(article readability.Article) string {
	trimmed := strings.TrimSpace(article.TextContent)
	collapsed := spaces.ReplaceAllString(trimmed, " ")

	paragraphs := []rune(collapsed)
	end := int(math.Min(float64(len(paragraphs)-1), MaxSummaryLength))

	var short string
	if len(paragraphs) > MaxSummaryLength {
		for ; end > MaxSummaryLength/2; end-- {
			if unicode.IsPunct(paragraphs[end]) || unicode.IsSpace(paragraphs[end]) {
				break
			}
		}

		// If stopped on a punctuation (like a comma) continue to walk backwards
		// until a letter is found. Since this function finally places an
		// elipsis at the end, a string ending like `hello, john` would output
		// as: `hello,...` which looks weird, so this makes sure the elipsis is
		// placed against a letter. If it fails, as with the above loop, it just
		// uses the max short body length cut in half as a fallback.
		if !unicode.IsLetter(paragraphs[end] - 1) {
			for ; end > MaxSummaryLength/2; end-- {
				if unicode.IsLetter(paragraphs[end]) {
					// shift forwards again so we don't chop off the last char.
					end += 1
					break
				}
			}
		}

		short = fmt.Sprint(string(paragraphs[:end]), "...")
	} else {
		short = string(paragraphs)
	}

	return short
}

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

func textFromNode(n *html.Node, preserveNewlines bool) string {
	var buf strings.Builder
	var last rune
	var collect func(*html.Node)
	collect = func(curr *html.Node) {
		switch curr.Type {
		case html.TextNode:
			data := curr.Data
			if data == "" {
				return
			}
			if preserveNewlines {
				buf.WriteString(data)
				return
			}

			normalized := spaces.ReplaceAllString(data, " ")
			normalized = strings.TrimSpace(normalized)
			if normalized == "" {
				return
			}

			first := []rune(normalized)[0]
			if buf.Len() > 0 && needsSpace(last, first) {
				buf.WriteByte(' ')
			}
			buf.WriteString(normalized)
			last = []rune(normalized)[len([]rune(normalized))-1]
		case html.ElementNode:
			if curr.DataAtom == atom.Br && preserveNewlines {
				buf.WriteByte('\n')
			}
			for c := curr.FirstChild; c != nil; c = c.NextSibling {
				collect(c)
			}
		default:
			for c := curr.FirstChild; c != nil; c = c.NextSibling {
				collect(c)
			}
		}
	}

	collect(n)

	return buf.String()
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
