package post

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// MaxShortBodyLength is the maximum length of the short summary text
const MaxShortBodyLength = 128

// MakeShortBody produces a short summary of a long piece of markdown content.
func MakeShortBody(long string) string {
	textonly := strings.Builder{}
	p := parser.New()
	tree := p.Parse([]byte(long))

	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		para, ok := n.(*ast.Paragraph)
		if ok {
			for _, c := range para.Children {
				text, ok := c.(*ast.Text)
				if ok && len(text.Literal) > 0 {
					oneline := bytes.ReplaceAll(text.Literal, []byte("\n"), []byte(" "))
					textonly.Write(oneline)
					textonly.WriteByte(' ')
				}
			}
		} else {
			container := n.AsContainer()
			if container != nil {
				children := container.Children
				for _, c := range children {
					walk(c)
				}
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

		return fmt.Sprint(string(paragraphs[:end]), "...")
	}

	return string(paragraphs)
}
