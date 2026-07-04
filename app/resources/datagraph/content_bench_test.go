package datagraph

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkNewRichText(b *testing.B) {
	cases := []struct {
		name  string
		input string
	}{
		{"small/plain", buildHTMLParagraphs(benchParasSmall)},
		{"medium/plain", buildHTMLParagraphs(benchParasMedium)},
		{"large/plain", buildHTMLParagraphs(benchParasLarge)},
		{"medium/with_links", buildHTMLWithLinks(benchParasMedium, benchLinksMedium)},
		{"large/with_links", buildHTMLWithLinks(benchParasLarge, benchLinksLarge)},
		{"medium/with_sdrs", buildHTMLWithSDRSeparators(benchParasMedium, 10)},
		{"large/with_sdrs", buildHTMLWithSDRSeparators(benchParasLarge, 30)},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetBytes(int64(len(tc.input)))
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewRichText(tc.input)
			}
		})
	}
}

func BenchmarkNewRichTextFromMarkdown(b *testing.B) {
	cases := []struct {
		name  string
		input string
	}{
		{"small/plain", buildMarkdownParagraphs(benchParasSmall)},
		{"medium/plain", buildMarkdownParagraphs(benchParasMedium)},
		{"large/plain", buildMarkdownParagraphs(benchParasLarge)},
		{"medium/with_links", buildMarkdownWithLinks(benchParasMedium, benchLinksMedium)},
		{"large/with_links", buildMarkdownWithLinks(benchParasLarge, benchLinksLarge)},
		{"medium/with_sdr_links", buildMarkdownWithSDRLinks(benchParasMedium, 10)},
		{"large/with_sdr_links", buildMarkdownWithSDRLinks(benchParasLarge, 30)},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetBytes(int64(len(tc.input)))
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewRichTextFromMarkdown(tc.input)
			}
		})
	}
}

func BenchmarkSplitCardParts(b *testing.B) {
	cases := []struct {
		name    string
		fixture Content
	}{
		{"small/no_splits", fixtureSmallContent},
		{"medium/no_splits", fixtureMediumContent},
		{"large/no_splits", fixtureLargeContent},
		{"medium/10_splits", fixtureSDRContent},
		{"medium/5_splits", mustBenchmarkContent(buildHTMLWithSDRSeparators(benchParasMedium, 5))},
		{"large/30_splits", mustBenchmarkContent(buildHTMLWithSDRSeparators(benchParasLarge, 30))},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tc.fixture.SplitCardParts()
			}
		})
	}
}

func BenchmarkContentHTML(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureSmallContent.HTML()
		}
	})
	b.Run("medium", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureMediumContent.HTML()
		}
	})
	b.Run("large", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureLargeContent.HTML()
		}
	})
}

func BenchmarkContentSplit(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureSmallContent.Split()
		}
	})
	b.Run("medium", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureMediumContent.Split()
		}
	})
	b.Run("large", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = fixtureLargeContent.Split()
		}
	})
}

// BenchmarkFullAgentPipeline simulates the hot path where an agent generates
// markdown with embedded SDR references, which is then parsed and split into
// card parts for rendering.
func BenchmarkFullAgentPipeline(b *testing.B) {
	cases := []struct {
		name  string
		input string
	}{
		{"small/3_sdrs", buildMarkdownWithSDRLinks(benchParasSmall, 3)},
		{"medium/10_sdrs", buildMarkdownWithSDRLinks(benchParasMedium, 10)},
		{"large/30_sdrs", buildMarkdownWithSDRLinks(benchParasLarge, 30)},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetBytes(int64(len(tc.input)))
			b.ReportAllocs()
			for b.Loop() {
				c, err := NewRichTextFromMarkdown(tc.input)
				if err != nil {
					b.Fatal(err)
				}
				_ = c.SplitCardParts()
			}
		})
	}
}

// Benchmark corpus sizes — chosen to surface real-world scaling behaviour.
// "small" ≈ a short reply, "medium" ≈ a typical forum post,
// "large" ≈ a long wiki page or newsletter article.
const (
	benchParasSmall  = 5
	benchParasMedium = 30
	benchParasLarge  = 150

	benchLinksSmall  = 3
	benchLinksMedium = 15
	benchLinksLarge  = 60
)

// Representative paragraph body (~200 chars). Varied enough that readability
// and the sanitiser don't get pathological cache hits.
const benchPara = "Community platforms thrive when members contribute knowledge freely. " +
	"Documentation, discussion threads, and curated node collections all play distinct roles. " +
	"A healthy signal-to-noise ratio depends on good moderation tooling and clear taxonomy."

// sdrIDs are valid xid strings used as SDR reference targets in benchmarks.
var sdrIDs = []string{
	"crk0gvqfunp7891n7ah0",
	"cn2h3gfljatbqvjqctdg",
	"cto7n8ifunp55p1bujv0",
	"cto7nm2funp55p1bujvg",
	"crk0gvqfunp7891n7ag0",
}

func buildHTMLParagraphs(n int) string {
	var b strings.Builder
	b.WriteString("<body>")
	for i := range n {
		fmt.Fprintf(&b, "<p>%s (paragraph %d)</p>\n", benchPara, i+1)
	}
	b.WriteString("</body>")
	return b.String()
}

func buildHTMLWithLinks(paras, links int) string {
	var b strings.Builder
	b.WriteString("<body>")
	for i := range paras {
		fmt.Fprintf(&b, "<p>%s (paragraph %d)</p>\n", benchPara, i+1)
		if i < links {
			fmt.Fprintf(&b,
				"<p>See <a href=\"https://example.com/resource/%d\">resource %d</a> for more.</p>\n",
				i, i)
		}
	}
	b.WriteString("</body>")
	return b.String()
}

func buildHTMLWithSDRSeparators(paras, sdrs int) string {
	var b strings.Builder
	b.WriteString("<body>")

	sdrEvery := max(1, paras/max(1, sdrs))
	sdrCount := 0
	for i := range paras {
		fmt.Fprintf(&b, "<p>%s (paragraph %d)</p>\n", benchPara, i+1)
		if sdrCount < sdrs && (i+1)%sdrEvery == 0 {
			id := sdrIDs[sdrCount%len(sdrIDs)]
			fmt.Fprintf(&b,
				"<p><a href=\"sdr:node/%s\">Reference %d</a></p>\n",
				id, sdrCount+1)
			sdrCount++
		}
	}
	b.WriteString("</body>")
	return b.String()
}

func buildMarkdownParagraphs(n int) string {
	var b strings.Builder
	for i := range n {
		fmt.Fprintf(&b, "%s (paragraph %d)\n\n", benchPara, i+1)
	}
	return b.String()
}

func buildMarkdownWithLinks(paras, links int) string {
	var b strings.Builder
	for i := range paras {
		fmt.Fprintf(&b, "%s (paragraph %d)\n\n", benchPara, i+1)
		if i < links {
			fmt.Fprintf(&b, "See [resource %d](https://example.com/resource/%d) for more.\n\n", i, i)
		}
	}
	return b.String()
}

func buildMarkdownWithSDRLinks(paras, sdrs int) string {
	var b strings.Builder
	sdrEvery := max(1, paras/max(1, sdrs))
	sdrCount := 0
	for i := range paras {
		fmt.Fprintf(&b, "%s (paragraph %d)\n\n", benchPara, i+1)
		if sdrCount < sdrs && (i+1)%sdrEvery == 0 {
			id := sdrIDs[sdrCount%len(sdrIDs)]
			fmt.Fprintf(&b, "[Reference %d](sdr:node/%s)\n\n", sdrCount+1, id)
			sdrCount++
		}
	}
	return b.String()
}

var (
	fixtureSmallContent  = mustBenchmarkContent(buildHTMLParagraphs(benchParasSmall))
	fixtureMediumContent = mustBenchmarkContent(buildHTMLParagraphs(benchParasMedium))
	fixtureLargeContent  = mustBenchmarkContent(buildHTMLParagraphs(benchParasLarge))
	fixtureSDRContent    = mustBenchmarkContent(buildHTMLWithSDRSeparators(benchParasMedium, 10))
)

func mustBenchmarkContent(html string) Content {
	c, err := NewRichText(html)
	if err != nil {
		panic(err)
	}
	return c
}
