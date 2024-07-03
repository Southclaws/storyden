package content

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"strings"
	"unicode"

	"github.com/Southclaws/fault"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var policy = bluemonday.UGCPolicy()

// MaxSummaryLength is the maximum length of the short summary text
const MaxSummaryLength = 128

// EmptyState is the default value for empty rich text. It could be an empty
// string but on the read path, it'll get turned into this either way.
const EmptyState = `<body></body>`

type Rich struct {
	html  *html.Node
	short string
	links []string
}

func (r Rich) HTML() string {
	if r.html == nil {
		return EmptyState
	}

	w := &bytes.Buffer{}

	err := html.Render(w, r.html)
	if err != nil {
		panic(err)
	}

	return w.String()
}

func (r Rich) Short() string {
	return r.short
}

func (r Rich) Links() []string {
	return r.links
}

// NewRichText will pull out any meaningful structured information from markdown
// document this includes a summary of the text and all link URLs for hydrating.
func NewRichText(raw string) (Rich, error) {
	sanitised := policy.Sanitize(raw)
	htmlTree, err := html.Parse(strings.NewReader(sanitised))
	if err != nil {
		return Rich{}, fault.Wrap(err)
	}

	return NewRichTextFromHTML(htmlTree)
}

func NewRichTextFromHTML(htmlTree *html.Node) (Rich, error) {
	bodyTree := &html.Node{}
	textonly := strings.Builder{}
	links := []string{}

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
						links = append(links, parsed.String())
					}
				}

			case atom.P:
				if n.Type == html.TextNode && len(n.Data) > 0 {

					oneline := strings.ReplaceAll(n.Data, "\n", " ")
					textonly.Write([]byte(oneline))
					textonly.WriteByte(' ')
					return
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

	paragraphs := []rune(strings.TrimSpace(textonly.String()))
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

	return Rich{
		html:  bodyTree,
		short: short,
		links: links,
	}, nil
}
