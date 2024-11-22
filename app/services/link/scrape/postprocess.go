package scrape

import (
	"bytes"
	"context"
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"golang.org/x/net/html"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (s *webScraper) postprocess(ctx context.Context, addr url.URL, r io.Reader) (*WebContent, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := metatable(doc)
	rc, err := getArticleContent(bytes.NewReader(buf), addr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	text := rc.Short()

	withBaseURL := func(urlOrPath string) string {
		if urlOrPath == "" {
			return ""
		}

		u, err := url.Parse(urlOrPath)
		if err != nil {
			return ""
		}

		if u.IsAbs() {
			return u.String()
		}

		return addr.ResolveReference(u).String()
	}

	wc := &WebContent{
		Title:       title(t),
		Description: description(t),
		Text:        text,
		Favicon:     withBaseURL(favicon(doc)),
		Image:       withBaseURL(t["og:image"]),
		Content:     rc,
	}

	return wc, nil
}

func getArticleContent(r io.Reader, pageURL url.URL) (datagraph.Content, error) {
	rc, err := datagraph.NewRichTextFromReader(r)
	if err != nil {
		return datagraph.Content{}, nil
	}

	return rc, nil
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

func favicon(doc *goquery.Document) string {
	if href, ok := doc.Find("link[rel='icon']").Attr("href"); ok {
		return href
	}

	if href, ok := doc.Find("link[rel='shortcut icon']").Attr("href"); ok {
		return href
	}

	if href, ok := doc.Find("link[rel='apple-touch-icon']").Attr("href"); ok {
		return href
	}

	if href, ok := doc.Find("link[rel='apple-touch-icon-precomposed']").Attr("href"); ok {
		return href
	}

	return "/favicon.ico"
}
