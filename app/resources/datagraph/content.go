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

	r := []html.Node{}

	// first, walk the tree for the top-most block-content nodes.
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.DataAtom {
			case
				atom.H1,
				atom.H2,
				atom.H3,
				atom.H4,
				atom.H5,
				atom.H6,
				atom.Blockquote,
				atom.Pre,
				atom.P:
				r = append(r, *n)
				// once split, exit out of this branch
				return
			}
		}

		if n.Type == html.TextNode {
			// if the text node is empty, skip it.
			if strings.TrimSpace(n.Data) == "" {
				return
			}

			r = append(r, *n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(c.html)

	// now, iterate these top level nodes and split any that are "too big"
	chunks := chunksFromNodes(r, roughMaxSentenceSize)

	return chunks
}

func chunksFromNodes(ns []html.Node, max int) []string {
	chunks := []string{}

	for _, n := range ns {
		t := textfromnode(&n)
		if len(t) > max {
			// TODO: Split logic
			chunks = append(chunks, splitearly(t, max)...)
		} else {
			chunks = append(chunks, t)
		}
	}

	return chunks
}

func splitearly(in string, max int) []string {
	var chunks []string
	var split func(s string)
	split = func(s string) {
		if len(s) <= max {
			chunks = append(chunks, strings.TrimSpace(s))
			return
		}

		upper := min(len(s), max) - 1
		if upper == -1 {
			// reached end of input stream
			return
		}

		lower := upper / 2
		boundary := upper
		fallback := -1
	outer:
		for ; boundary > lower; boundary-- {
			c := s[boundary]
			switch c {
			// very rudimentary sentence boundaries (latin only at the moment)
			case '.', ';', '!', '?':
				break outer
			// worst case: no boundaries found, use the closest space
			case ' ':
				if fallback == -1 {
					fallback = boundary
				}
			}
		}

		if boundary <= lower {
			if fallback > -1 {
				// worst case: no sent boundaries, split at fallback position.
				boundary = fallback
			} else {
				// worst case: no fallback either (the input string was a solid
				// block of text with no spaces or sentence boundaries.)
				boundary = upper
			}
		}

		left := strings.TrimSpace(s[:boundary])
		right := strings.TrimSpace(s[boundary+1:])
		chunks = append(chunks, left)

		if len(right) > 0 {
			split(right)
		}
	}
	split(in)

	return chunks
}

func textfromnode(n *html.Node) string {
	var collect func(*html.Node, *strings.Builder)
	collect = func(cc *html.Node, buf *strings.Builder) {
		if cc.Type == html.TextNode {
			buf.WriteString(cc.Data)
		}
		for c := cc.FirstChild; c != nil; c = c.NextSibling {
			collect(c, buf)
		}
	}
	buf := &strings.Builder{}
	collect(n, buf)
	return buf.String()
}
