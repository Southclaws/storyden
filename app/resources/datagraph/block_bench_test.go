package datagraph

import (
	"fmt"
	"strings"
	"testing"
)

func makeParagraphs(n int) string {
	var sb strings.Builder
	for i := range n {
		fmt.Fprintf(&sb, "<p>Paragraph number %d with some realistic body text to exercise the matching algorithm properly.</p>\n", i+1)
	}
	return sb.String()
}

func makeListItems(n int) string {
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	for i := range n {
		fmt.Fprintf(&sb, "<li>List item %d with some content for benchmarking purposes here.</li>\n", i+1)
	}
	sb.WriteString("</ul>\n")
	return sb.String()
}

func makeCustomBlocks(n int) string {
	var sb strings.Builder
	for i := range n {
		// Simulate opaque TipTap node views with data attributes – the backend
		// must not try to interpret them.
		fmt.Fprintf(&sb, `<div data-type="linkCard" data-url="https://example.com/%d" data-title="Title %d"><p>Link card %d</p></div>`+"\n", i, i, i)
	}
	return sb.String()
}

var smallArticle = `
<h1>Introduction</h1>
<p>This is a short article about benchmarking the block ID assignment system.</p>
<p>It has a few paragraphs and a list.</p>
<ul>
  <li>First item</li>
  <li>Second item</li>
  <li>Third item</li>
</ul>
<p>And a conclusion paragraph at the end of the document.</p>
`

var largeArticle = func() string {
	var sb strings.Builder
	sb.WriteString("<h1>Large Article</h1>\n")
	sb.WriteString(makeParagraphs(80))
	sb.WriteString(makeListItems(40))
	sb.WriteString(makeParagraphs(30))
	return sb.String()
}()

var manyParagraphs = makeParagraphs(200)

var customDivArticle = func() string {
	var sb strings.Builder
	sb.WriteString("<h2>Rich content</h2>\n")
	sb.WriteString(makeCustomBlocks(20))
	sb.WriteString(makeParagraphs(30))
	return sb.String()
}()

func benchmarkNewRichText(b *testing.B, html string) {
	b.Helper()
	b.ReportAllocs()
	for b.Loop() {
		_, err := NewRichText(html)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkWithStableBlockIDs(b *testing.B, html string) {
	b.Helper()
	prevContent, err := NewRichText(html)
	if err != nil {
		b.Fatal(err)
	}
	prev, err := NewRichTextWithNewBlocks(prevContent)
	if err != nil {
		b.Fatal(err)
	}

	// Simulate a minor edit: append a space to the last paragraph.
	edited := strings.Replace(html, "</p>\n", " </p>\n", 1)
	next, err := NewRichText(edited)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, err := NewRichTextWithChangedBlocks(prev.Content, next)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// adversarialEdit simulates a bulk find-replace where every paragraph changes
// slightly (no exact-HTML fast path hits, all blocks reach Levenshtein).
// This exercises the O(N*M*L²) worst case of the LCS DP.
func adversarialEdit(html string) string {
	return strings.ReplaceAll(html, "Paragraph number", "Paragraph no.")
}

func benchmarkWithStableBlockIDsAdversarial(b *testing.B, html string) {
	b.Helper()
	prevContent, err := NewRichText(html)
	if err != nil {
		b.Fatal(err)
	}
	prev, err := NewRichTextWithNewBlocks(prevContent)
	if err != nil {
		b.Fatal(err)
	}

	edited := adversarialEdit(html)
	next, err := NewRichText(edited)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, err := NewRichTextWithChangedBlocks(prev.Content, next)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewRichText_Small(b *testing.B)  { benchmarkNewRichText(b, smallArticle) }
func BenchmarkNewRichText_Large(b *testing.B)  { benchmarkNewRichText(b, largeArticle) }
func BenchmarkNewRichText_ManyP(b *testing.B)  { benchmarkNewRichText(b, manyParagraphs) }
func BenchmarkNewRichText_Custom(b *testing.B) { benchmarkNewRichText(b, customDivArticle) }

func BenchmarkWithStableBlockIDs_Small(b *testing.B) { benchmarkWithStableBlockIDs(b, smallArticle) }
func BenchmarkWithStableBlockIDs_Large(b *testing.B) { benchmarkWithStableBlockIDs(b, largeArticle) }
func BenchmarkWithStableBlockIDs_ManyP(b *testing.B) { benchmarkWithStableBlockIDs(b, manyParagraphs) }
func BenchmarkWithStableBlockIDs_Custom(b *testing.B) {
	benchmarkWithStableBlockIDs(b, customDivArticle)
}

func BenchmarkWithStableBlockIDs_AdversarialLarge(b *testing.B) {
	benchmarkWithStableBlockIDsAdversarial(b, largeArticle)
}

func BenchmarkWithStableBlockIDs_AdversarialManyP(b *testing.B) {
	benchmarkWithStableBlockIDsAdversarial(b, manyParagraphs)
}
