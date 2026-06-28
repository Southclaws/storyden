package pluginbuilder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	storydenDocsBaseURL = "https://www.storyden.org"
	storydenDocsMaxSize = 128_000
)

type StorydenDocsInput struct {
	Path     string `json:"path" jsonschema:"Storyden documentation path or URL. Use /llms.txt, /docs/..., or https://www.storyden.org/docs/..."`
	MaxBytes int    `json:"max_bytes,omitempty" jsonschema:"Maximum response bytes to return"`
}

type StorydenDocsResult struct {
	URL       string `json:"url"`
	Content   string `json:"content"`
	Truncated bool   `json:"truncated"`
}

func (a *Agent) addStorydenDocsTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name: "plugin_storyden_docs",
		Description: `Fetch constrained Storyden documentation as Markdown or text.

Use this for Storyden product, plugin, SDK, API, and manifest documentation that
is not discoverable from Go package symbols alone. This is not a general web
browser: it only reads https://www.storyden.org/llms.txt and subpaths under
https://www.storyden.org/docs/.

Use "/docs/introduction/members/permissions" when choosing manifest access
permissions for plugins that call Storyden host HTTP APIs.

Pages on storyden.org support Markdown responses. This tool automatically
requests the Markdown route for docs pages and sends
"Accept: text/markdown, text/plain" so the response is agent-friendly.

Start with "/llms.txt" to discover documentation entry points, then fetch
specific "/docs/..." paths as needed.`,
	}, func(ctx adktool.Context, args StorydenDocsInput) (StorydenDocsResult, error) {
		return FetchStorydenDocs(ctx, args)
	}))
}

func FetchStorydenDocs(ctx context.Context, in StorydenDocsInput) (StorydenDocsResult, error) {
	u, err := storydenDocsURL(in.Path)
	if err != nil {
		return StorydenDocsResult{}, err
	}

	limit := in.MaxBytes
	if limit <= 0 || limit > storydenDocsMaxSize {
		limit = storydenDocsMaxSize
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return StorydenDocsResult{}, err
	}
	req.Header.Set("Accept", "text/markdown, text/plain;q=0.9, text/*;q=0.8")
	req.Header.Set("User-Agent", "Storyden-Plugin-Builder/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return StorydenDocsResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return StorydenDocsResult{}, fmt.Errorf("storyden docs request failed: status %d", resp.StatusCode)
	}

	reader := io.LimitReader(resp.Body, int64(limit)+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return StorydenDocsResult{}, err
	}

	truncated := len(data) > limit
	if truncated {
		data = data[:limit]
	}

	return StorydenDocsResult{
		URL:       u.String(),
		Content:   string(data),
		Truncated: truncated,
	}, nil
}

func storydenDocsURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("path is required")
	}

	var u *url.URL
	var err error
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err = url.Parse(raw)
		if err != nil {
			return nil, err
		}
	} else {
		if !strings.HasPrefix(raw, "/") {
			raw = "/" + raw
		}
		base, err := url.Parse(storydenDocsBaseURL)
		if err != nil {
			return nil, err
		}
		u = base.ResolveReference(&url.URL{Path: raw})
	}

	if u.Scheme != "https" {
		return nil, fmt.Errorf("storyden docs URL must use https: %s", u.String())
	}
	if u.Host != "www.storyden.org" && u.Host != "storyden.org" {
		return nil, fmt.Errorf("storyden docs URL host is not allowed: %s", u.Host)
	}
	u.Host = "www.storyden.org"
	u.RawQuery = ""
	u.Fragment = ""
	u.Path = path.Clean("/" + strings.TrimPrefix(u.Path, "/"))

	if u.Path == "/llms.txt" {
		return u, nil
	}
	if strings.HasPrefix(u.Path, "/docs/") || u.Path == "/docs" {
		u.Path = storydenDocsMarkdownPath(u.Path)
		return u, nil
	}

	return nil, fmt.Errorf("storyden docs path is not allowed: %s; use /llms.txt or /docs/...", u.Path)
}

func storydenDocsMarkdownPath(p string) string {
	if p == "/docs" || strings.HasSuffix(p, ".md") {
		return p
	}
	return p + ".md"
}
