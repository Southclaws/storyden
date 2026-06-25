package presentation

import (
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

const (
	DataRenderCard = "data-render_card"
)

type PartKind string

const (
	PartText       PartKind = "text"
	PartRenderCard PartKind = "render_card"
)

type Part struct {
	Kind PartKind
	Text string
	Ref  *datagraph.Ref
}

type RenderCardData struct {
	Ref  string `json:"ref"`
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

var paragraphSeparator = regexp.MustCompile(`\n[ \t]*\n+`)

func Parse(input string) []Part {
	if input == "" {
		return nil
	}

	return parseSDRPresentationMarkdown(input)
}

func parseSDRPresentationMarkdown(input string) []Part {
	matches := paragraphSeparator.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return parseSDRPresentationParagraph(input)
	}

	parts := []Part{}
	position := 0
	for _, match := range matches {
		paragraph := input[position:match[0]]
		separator := input[match[0]:match[1]]
		parts = append(parts, parseSDRPresentationParagraph(paragraph)...)
		if separator != "" {
			parts = append(parts, Part{Kind: PartText, Text: separator})
		}
		position = match[1]
	}
	parts = append(parts, parseSDRPresentationParagraph(input[position:])...)

	return mergeTextParts(parts)
}

func parseSDRPresentationParagraph(paragraph string) []Part {
	if paragraph == "" {
		return nil
	}

	ref, ok := standaloneSDRReference(paragraph)
	if !ok {
		return []Part{{Kind: PartText, Text: paragraph}}
	}

	return []Part{{Kind: PartRenderCard, Ref: ref}}
}

func standaloneSDRReference(markdown string) (*datagraph.Ref, bool) {
	content, err := datagraph.NewRichTextFromMarkdown(markdown)
	if err != nil {
		return nil, false
	}

	refs := content.References()
	if len(refs) != 1 {
		return nil, false
	}

	ref, ok := standaloneSDRReferenceFromHTML(content.HTML())
	if !ok {
		return nil, false
	}

	if ref.ID != refs[0].ID || ref.Kind != refs[0].Kind {
		return nil, false
	}

	return ref, true
}

func standaloneSDRReferenceFromHTML(rawHTML string) (*datagraph.Ref, bool) {
	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))

	var ref *datagraph.Ref
	anchorDepth := 0
	anchorCount := 0

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return ref, tokenizer.Err() == io.EOF && anchorCount == 1

		case html.StartTagToken:
			name, hasAttr := tokenizer.TagName()
			tag := string(name)
			switch tag {
			case "html", "head", "body", "p":
				continue
			case "a":
				anchorCount++
				if anchorCount > 1 {
					return nil, false
				}

				href := ""
				for hasAttr {
					var key, value []byte
					key, value, hasAttr = tokenizer.TagAttr()
					if strings.EqualFold(string(key), "href") {
						href = string(value)
					}
				}

				parsed, err := url.Parse(href)
				if err != nil {
					return nil, false
				}

				parsedRef, err := datagraph.NewRefFromSDR(*parsed)
				if err != nil {
					return nil, false
				}

				ref = parsedRef
				anchorDepth++

			default:
				return nil, false
			}

		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			tag := string(name)
			switch tag {
			case "html", "head", "body", "p":
				continue
			case "a":
				if anchorDepth == 0 {
					return nil, false
				}
				anchorDepth--
			default:
				return nil, false
			}

		case html.TextToken:
			if anchorDepth > 0 {
				continue
			}
			if strings.TrimSpace(string(tokenizer.Text())) != "" {
				return nil, false
			}
		}
	}
}

func NewRenderCardData(ref *datagraph.Ref) RenderCardData {
	return RenderCardData{
		Ref:  ref.String(),
		Kind: ref.Kind.String(),
		ID:   ref.ID.String(),
	}
}

func mergeTextParts(parts []Part) []Part {
	if len(parts) < 2 {
		return parts
	}

	merged := []Part{}
	for _, part := range parts {
		if part.Kind == PartText && len(merged) > 0 && merged[len(merged)-1].Kind == PartText {
			merged[len(merged)-1].Text += part.Text
			continue
		}
		merged = append(merged, part)
	}

	return merged
}
