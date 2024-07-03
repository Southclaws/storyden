package url

import (
	"bytes"
	"context"
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/cixtor/readability"
	"golang.org/x/net/html"
)

func (s *webScraper) postprocess(ctx context.Context, addr string, r io.Reader) (*WebContent, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := metatable(doc)
	rc, text, err := getArticleContent(bytes.NewReader(buf), addr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	wc := &WebContent{
		Title:       title(t),
		Description: description(t),
		Text:        text,
		Image:       t["og:image"],
		Content:     rc,
	}

	return wc, nil
}

func getArticleContent(r io.Reader, pageURL string) (content.Rich, string, error) {
	result, err := readability.New().Parse(r, pageURL)
	if err != nil {
		return content.Rich{}, "", fault.Wrap(err)
	}

	rc, err := content.NewRichTextFromHTML(result.Node)
	if err != nil {
		return content.Rich{}, result.TextContent, nil
	}

	return rc, result.TextContent, nil
}

func metatable(doc *goquery.Document) map[string]string {
	return dt.Reduce(doc.Find("head > meta").Nodes, func(wc map[string]string, n *html.Node) map[string]string {
		k, v := ogtable(n.Attr)
		if k != "" && v != "" {
			wc[k] = v
		}

		return wc
	}, map[string]string{})
}

func ogtable(attrs []html.Attribute) (k string, v string) {
	for _, a := range attrs {
		switch a.Key {
		case "property":
			k = a.Val
		case "name":
			k = a.Val
		case "content":
			v = a.Val
		}
	}

	return
}

func title(t map[string]string) string {
	if t["og:title"] != "" {
		return t["og:title"]
	}
	if t["title"] != "" {
		return t["title"]
	}
	if t["og:site_name"] != "" {
		return t["og:site_name"]
	}
	if t["og:url"] != "" {
		return t["og:url"]
	}
	if t["title"] != "" {
		return t["title"]
	}

	return ""
}

func description(t map[string]string) string {
	if t["og:description"] != "" {
		return t["og:description"]
	}
	if t["description"] != "" {
		return t["description"]
	}

	return ""
}
