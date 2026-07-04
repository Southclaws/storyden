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

	// Allow all data attributes that TipTap may add.
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

// ContentWithBlocks is Content whose Storyden-owned block IDs are meaningful.
type ContentWithBlocks struct {
	Content
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
	if r.html == nil {
		return true
	}

	if strings.TrimSpace(r.plain) != "" {
		return false
	}

	if len(r.links) > 0 {
		return false
	}

	if len(r.media) > 0 {
		return false
	}

	if len(r.sdrs) > 0 {
		return false
	}

	return true
}

// NewRichText sanitises and parses HTML into Content without assigning block IDs.
//
// The returned Content stores a normalised <body> tree, plaintext, summary,
// outbound links, media URLs, and Storyden references.
func NewRichText(raw string) (Content, error) {
	return NewRichTextFromReader(strings.NewReader(raw))
}

func NewRichTextFromMarkdown(md string) (Content, error) {
	html := blackfriday.Run([]byte(md), blackfriday.WithExtensions(
		blackfriday.NoEmptyLineBeforeBlock,
	))

	return NewRichTextFromReader(strings.NewReader(string(html)))
}

func NewRichTextFromReader(r io.Reader) (Content, error) {
	return parseRichTextFromReader(r, false)
}

// NewRichTextWithBlocks parses stored HTML and preserves existing Storyden block
// IDs, but never creates new IDs. This is the read-path constructor for content
// types that support addressable blocks.
func NewRichTextWithBlocks(raw string) (ContentWithBlocks, error) {
	c, err := parseRichTextFromReader(strings.NewReader(raw), true)
	if err != nil {
		return ContentWithBlocks{}, err
	}

	return ContentWithBlocks{Content: c}, nil
}

// NewRichTextWithNewBlocks prepares newly created Content for storage by
// assigning valid Storyden block IDs to addressable block elements.
func NewRichTextWithNewBlocks(c Content) (ContentWithBlocks, error) {
	stable, err := c.withBlockIDs()
	if err != nil {
		return ContentWithBlocks{}, err
	}

	return ContentWithBlocks{Content: stable}, nil
}

// NewRichTextWithChangedBlocks prepares updated Content for storage by
// preserving block IDs from previous where possible and assigning IDs to new
// addressable blocks.
func NewRichTextWithChangedBlocks(previous Content, next Content) (ContentWithBlocks, error) {
	stable, err := next.withPreviousState(&previous)
	if err != nil {
		return ContentWithBlocks{}, err
	}

	return ContentWithBlocks{Content: stable}, nil
}

func parseRichTextFromReader(r io.Reader, preserveBlockIDs bool) (Content, error) {
	baseURL, err := url.Parse("ignore:")
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

	result, err := readability.New().Parse(bytes.NewReader(sanitised), baseURL.String())
	if err != nil {
		return Content{}, fault.Wrap(err)
	}

	short := getSummary(result)

	bodyTree, links, media, refs := extractReferences(htmlTree, baseURL)
	normaliseIDAttributes(bodyTree, preserveBlockIDs)

	c := Content{
		html:  bodyTree,
		short: short,
		plain: result.TextContent,
		links: links,
		media: media,
		sdrs:  refs,
	}

	return c, nil
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
