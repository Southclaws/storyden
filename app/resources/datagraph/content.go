package datagraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/Southclaws/fault"
	"github.com/cixtor/readability"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
	"go.uber.org/zap"
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

func (r Content) References() RefList {
	return r.sdrs
}

type options struct {
	baseURL string
}
type option func(*options)

// NewRichText will pull out any meaningful structured information from markdown
// document this includes a summary of the text and all link URLs for hydrating.
func NewRichText(raw string) (Content, error) {
	return NewRichTextFromReader(strings.NewReader(raw))
}

func NewRichTextFromReader(r io.Reader, opts ...option) (Content, error) {
	o := options{baseURL: "ignore:"}
	for _, opt := range opts {
		opt(&o)
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

	bodyTree, links, refs := extractReferences(htmlTree)

	return Content{
		html:  bodyTree,
		short: short,
		plain: result.TextContent,
		links: links,
		sdrs:  refs,
	}, nil
}

func extractReferences(htmlTree *html.Node) (*html.Node, []string, RefList) {
	bodyTree := &html.Node{}
	links := []string{}
	sdrs := []url.URL{}

	if htmlTree.DataAtom == atom.Body {
		bodyTree = htmlTree
	}

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Parent != nil {
			switch n.Parent.DataAtom {
			case atom.A:
				href, hasHref := lo.Find(n.Parent.Attr, func(a html.Attribute) bool {
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
			zap.L().Warn("invalid SDR in content", zap.Error(err), zap.String("ref", v.Opaque))
			continue
		}
		refs = append(refs, r)
	}

	return bodyTree, links, refs
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
