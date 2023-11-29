package extractor

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"strings"
	"unicode"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// MaxShortBodyLength is the maximum length of the short summary text
const MaxShortBodyLength = 128

type EnrichedProperties struct {
	Short string
	Links []string
}

// Destructure will pull out any meaningful structured information from markdown
// document this includes a summary of the text and all link URLs for hydrating.
func Destructure(markdown string) EnrichedProperties {
	textonly := strings.Builder{}
	p := parser.New()
	tree := p.Parse([]byte(markdown))

	var short string
	links := []string{}

	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		switch node := n.(type) {
		case *ast.Link:
			if parsed, err := url.Parse(string(node.Destination)); err == nil {
				links = append(links, parsed.String())
			}

		case *ast.Text:
			if len(node.Literal) == 0 {
				return
			}

			oneline := bytes.ReplaceAll(node.Literal, []byte("\n"), []byte(" "))
			textonly.Write(oneline)
			textonly.WriteByte(' ')

		default:
			container := n.AsContainer()
			if container == nil {
				return
			}

			children := container.Children
			for _, c := range children {
				walk(c)
			}
		}
	}
	walk(tree)

	paragraphs := []rune(strings.TrimSpace(textonly.String()))
	end := int(math.Min(float64(len(paragraphs)-1), MaxShortBodyLength))

	if len(paragraphs) > MaxShortBodyLength {
		for ; end > MaxShortBodyLength/2; end-- {
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
			for ; end > MaxShortBodyLength/2; end-- {
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

	return EnrichedProperties{
		Short: short,
		Links: links,
	}
}
