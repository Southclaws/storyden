package datagraph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// mustWithIDs runs NewRichTextWithChangedBlocks, failing the test on error.
func mustWithIDs(t *testing.T, next Content, prev *Content) Content {
	t.Helper()
	require.NotNil(t, prev)
	c, err := NewRichTextWithChangedBlocks(*prev, next)
	require.NoError(t, err)
	return c.Content
}

// idOf returns the ID of the first block whose text contains substr.
func idOf(t *testing.T, c Content, substr string) string {
	t.Helper()
	for _, b := range c.Blocks() {
		if strings.Contains(b.Text, substr) {
			return b.ID
		}
	}
	t.Fatalf("no block found containing %q", substr)
	return ""
}

// idOfType returns the ID of the first block with typ whose text contains substr.
func idOfType(t *testing.T, c Content, typ, substr string) string {
	t.Helper()
	for _, b := range c.Blocks() {
		if b.Type == typ && strings.Contains(b.Text, substr) {
			return b.ID
		}
	}
	t.Fatalf("no %s block found containing %q", typ, substr)
	return ""
}

// idsOfType returns all IDs for blocks with typ whose text contains substr.
func idsOfType(c Content, typ, substr string) []string {
	var ids []string
	for _, b := range c.Blocks() {
		if b.Type == typ && strings.Contains(b.Text, substr) {
			ids = append(ids, b.ID)
		}
	}
	return ids
}

// assertSameID asserts that the block containing substr has the same ID in prev and next.
func assertSameID(t *testing.T, prev, next Content, substr string) {
	t.Helper()
	assert.Equal(t, idOf(t, prev, substr), idOf(t, next, substr), "block containing %q should keep its ID", substr)
}

// assertSameIDOfType asserts that the typed block containing substr keeps its ID.
func assertSameIDOfType(t *testing.T, prev, next Content, typ, substr string) {
	t.Helper()
	assert.Equal(t, idOfType(t, prev, typ, substr), idOfType(t, next, typ, substr), "%s block containing %q should keep its ID", typ, substr)
}

// assertDiffID asserts that the block containing substr has a different ID in prev and next.
func assertDiffID(t *testing.T, prev, next Content, substr string) {
	t.Helper()
	assert.NotEqual(t, idOf(t, prev, substr), idOf(t, next, substr), "block containing %q should get a new ID", substr)
}

// assertAllHaveIDs asserts that every block in c has a valid block ID.
func assertAllHaveIDs(t *testing.T, c Content) {
	t.Helper()
	for _, b := range c.Blocks() {
		assert.True(t, isValidBlockID(b.ID), "block type=%q text=%q has invalid ID %q", b.Type, b.Text, b.ID)
	}
}

// assertUniqueIDs asserts that no two blocks in c share an ID.
func assertUniqueIDs(t *testing.T, c Content) {
	t.Helper()
	seen := map[string]bool{}
	for _, b := range c.Blocks() {
		assert.False(t, seen[b.ID], "duplicate block ID %q (type=%q text=%q)", b.ID, b.Type, b.Text)
		seen[b.ID] = true
	}
}

// WithBlockIDs tests

// TestBlock_FreshContentGetsIDs verifies that plain HTML without IDs gets IDs assigned.
func TestBlock_FreshContentGetsIDs(t *testing.T) {
	c := mustContentWithNewBlocks(t, "<p>Hello world</p><p>Goodbye world</p>")
	assertAllHaveIDs(t, c)
	assertUniqueIDs(t, c)
	require.Len(t, c.Blocks(), 2)
}

// TestBlock_WithBlockIDsIdempotent verifies that content already carrying valid IDs
// passes through the read-path constructor unchanged.
func TestBlock_WithBlockIDsIdempotent(t *testing.T) {
	first := mustContentWithNewBlocks(t, "<p>Hello world</p>")
	firstID := idOf(t, first, "Hello")

	second := mustContentWithBlocks(t, first.HTML())
	assert.Equal(t, firstID, idOf(t, second, "Hello"))
}

func TestBlock_NewRichTextWithBlocksDoesNotAssignMissingIDs(t *testing.T) {
	knownID := "sdb_cv3j6fld0aaaab44kk50"
	raw := `<p id="` + knownID + `">Has ID already</p><p>Needs an ID</p>`
	c := mustContentWithBlocks(t, raw)

	blocks := c.Blocks()
	require.Len(t, blocks, 2)

	assert.Equal(t, knownID, blocks[0].ID, "existing valid ID must be preserved")
	assert.Empty(t, blocks[1].ID, "read path must not assign missing IDs")
}

func TestBlock_NewRichTextWithNewBlocksIgnoresSubmittedIDs(t *testing.T) {
	knownID := "sdb_cv3j6fld0aaaab44kk50"
	c := mustContentWithNewBlocks(t, `<p id="`+knownID+`">Submitted ID is not trusted.</p>`)

	require.Len(t, c.Blocks(), 1)
	assert.True(t, isValidBlockID(c.Blocks()[0].ID))
	assert.NotEqual(t, knownID, c.Blocks()[0].ID)
}

func TestBlock_NewRichTextRepairsDuplicateBlockID(t *testing.T) {
	knownID := "sdb_cv3j6fld0aaaab44kk50"
	c := mustContentWithNewBlocks(t, `<p id="`+knownID+`">First block.</p><p id="`+knownID+`">Second block.</p>`)

	blocks := c.Blocks()
	require.Len(t, blocks, 2)
	assertAllHaveIDs(t, c)
	assertUniqueIDs(t, c)
	assert.NotEqual(t, blocks[0].ID, blocks[1].ID)
}

func TestBlock_ContainerElementsGetIDs(t *testing.T) {
	c := mustContentWithNewBlocks(t, `<div>Intro container</div><section>Section container</section><article>Article container</article><ul><li>List item</li></ul><ol><li>Ordered item</li></ol>`)

	blocks := c.Blocks()
	require.Len(t, blocks, 7)
	assertAllHaveIDs(t, c)
	assertUniqueIDs(t, c)
	assert.Equal(t, []string{"div", "section", "article", "ul", "li", "ol", "li"}, blockTypes(blocks))
}

func TestBlock_ImagesGetIDs(t *testing.T) {
	c := mustContentWithNewBlocks(t, `<p>Image follows.</p><img src="https://example.com/image.jpg" alt="Example image">`)

	blocks := c.Blocks()
	require.Len(t, blocks, 2)
	assertAllHaveIDs(t, c)
	assertUniqueIDs(t, c)
	assert.Equal(t, []string{"p", "img"}, blockTypes(blocks))
	assert.Contains(t, c.HTML(), `<img src="https://example.com/image.jpg" alt="Example image" id="sdb_`)
}

func TestBlock_NewBlocksDoesNotMutateInput(t *testing.T) {
	input := mustContent(t, `<blockquote><p>Nested paragraph stays unassigned.</p></blockquote><p>Second paragraph stays unassigned.</p>`)
	inputHTML := input.HTML()

	stable, err := NewRichTextWithNewBlocks(input)
	require.NoError(t, err)

	assert.Equal(t, inputHTML, input.HTML())
	assert.NotContains(t, input.HTML(), ` id=`)
	assertAllHaveIDs(t, stable.Content)
	assertUniqueIDs(t, stable.Content)
}

// WithStableBlockIDs tests (update path)

// TestBlock_AppendParagraph verifies that appending a paragraph preserves existing IDs.
func TestBlock_AppendParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>First paragraph with some words.</p><p>Second paragraph with some words.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>First paragraph with some words.</p><p>Second paragraph with some words.</p><p>Newly appended paragraph at the end.</p>"), &prev)

	assertSameID(t, prev, next, "First paragraph")
	assertSameID(t, prev, next, "Second paragraph")
	assert.True(t, isValidBlockID(idOf(t, next, "Newly appended")))
}

// TestBlock_PrependParagraph verifies that prepending a paragraph preserves existing IDs.
func TestBlock_PrependParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>First paragraph with some words.</p><p>Second paragraph with some words.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>Brand new paragraph at the very beginning.</p><p>First paragraph with some words.</p><p>Second paragraph with some words.</p>"), &prev)

	assertSameID(t, prev, next, "First paragraph")
	assertSameID(t, prev, next, "Second paragraph")
	assert.True(t, isValidBlockID(idOf(t, next, "Brand new")))
}

// TestBlock_InsertMiddleParagraph verifies that inserting in the middle preserves flanking IDs.
func TestBlock_InsertMiddleParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>Opening paragraph with enough words here.</p><p>Closing paragraph with enough words here.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>Opening paragraph with enough words here.</p><p>Middle paragraph inserted between the two.</p><p>Closing paragraph with enough words here.</p>"), &prev)

	assertSameID(t, prev, next, "Opening paragraph")
	assertSameID(t, prev, next, "Closing paragraph")
	assert.True(t, isValidBlockID(idOf(t, next, "Middle paragraph")))
}

// TestBlock_DeleteFirstParagraph verifies that deleting the first paragraph leaves the rest unchanged.
func TestBlock_DeleteFirstParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>This one will be deleted from the document.</p><p>This one stays in the document.</p><p>This one also stays in the document.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>This one stays in the document.</p><p>This one also stays in the document.</p>"), &prev)

	assertSameID(t, prev, next, "This one stays")
	assertSameID(t, prev, next, "This one also stays")
}

// TestBlock_DeleteLastParagraph verifies that deleting the last paragraph leaves the rest unchanged.
func TestBlock_DeleteLastParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>First paragraph remains in the document here.</p><p>Second paragraph remains in the document here.</p><p>Last paragraph will be removed entirely.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>First paragraph remains in the document here.</p><p>Second paragraph remains in the document here.</p>"), &prev)

	assertSameID(t, prev, next, "First paragraph")
	assertSameID(t, prev, next, "Second paragraph")
}

// TestBlock_DeleteMiddleParagraph verifies that deleting a middle paragraph leaves flanking IDs unchanged.
func TestBlock_DeleteMiddleParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>First paragraph stays in the document here.</p><p>Middle paragraph will be removed from the document.</p><p>Last paragraph stays in the document here.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>First paragraph stays in the document here.</p><p>Last paragraph stays in the document here.</p>"), &prev)

	assertSameID(t, prev, next, "First paragraph stays")
	assertSameID(t, prev, next, "Last paragraph stays")
}

// TestBlock_MinorEditPreservesID verifies that a single-word change (>= 80% similarity) keeps the ID.
func TestBlock_MinorEditPreservesID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>The quick brown fox jumps over the lazy dog always.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>The quick brown fox jumped over the lazy dog always.</p>"), &prev)

	assertSameID(t, prev, next, "quick brown fox")
}

// TestBlock_LargeRewriteGetsNewID verifies that a completely different paragraph gets a new ID.
func TestBlock_LargeRewriteGetsNewID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>The quick brown fox jumps over the lazy dog always.</p>")
	prevID := prev.Blocks()[0].ID

	next := mustWithIDs(t, mustContent(t, "<p>Entirely different content with nothing in common at all.</p>"), &prev)

	assert.NotEqual(t, prevID, next.Blocks()[0].ID)
}

func TestBlock_CopiedOldIDOnLargeRewriteGetsNewID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>The quick brown fox jumps over the lazy dog while writing a long enough paragraph.</p>")
	prevID := prev.Blocks()[0].ID

	nextHTML := `<p id="` + prevID + `">Database migrations require careful planning, backup validation, and staged rollout procedures.</p>`
	next := mustWithIDs(t, mustContent(t, nextHTML), &prev)

	assert.NotEqual(t, prevID, next.Blocks()[0].ID)
}

// TestBlock_ExactSameContentKeepsIDs verifies that identical content round-trips with the same IDs.
func TestBlock_ExactSameContentKeepsIDs(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>Same paragraph one here.</p><p>Same paragraph two here.</p>")
	next := mustWithIDs(t, mustContent(t, "<p>Same paragraph one here.</p><p>Same paragraph two here.</p>"), &prev)

	assertSameID(t, prev, next, "Same paragraph one")
	assertSameID(t, prev, next, "Same paragraph two")
}

// TestBlock_ReorderTwoParagraphs verifies that exact moved paragraphs retain IDs.
func TestBlock_ReorderTwoParagraphs(t *testing.T) {
	p1 := "<p>Dogs are loyal companions that enjoy outdoor activities and fetch.</p>"
	p2 := "<p>Mathematics underpins all engineering disciplines across every field.</p>"
	prev := mustContentWithNewBlocks(t, p1+p2)
	next := mustWithIDs(t, mustContent(t, p2+p1), &prev)

	// Both blocks must have valid IDs and no duplicates.
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
	// At least one block must have kept its ID from prev.
	keptIDs := 0
	prevIDs := map[string]bool{}
	for _, b := range prev.Blocks() {
		prevIDs[b.ID] = true
	}
	for _, b := range next.Blocks() {
		if prevIDs[b.ID] {
			keptIDs++
		}
	}
	assert.GreaterOrEqual(t, keptIDs, 1, "at least one block must keep its ID when two paragraphs are swapped")
}

// TestBlock_ReorderThreeParagraphs verifies that exact matches survive larger reorders.
func TestBlock_ReorderThreeParagraphs(t *testing.T) {
	p1 := "<p>Dogs are loyal companions that enjoy outdoor activities and fetch.</p>"
	p2 := "<p>Mathematics underpins all engineering disciplines across every field.</p>"
	p3 := "<p>Ancient Romans built remarkable infrastructure spanning their vast empire.</p>"
	prev := mustContentWithNewBlocks(t, p1+p2+p3)
	next := mustWithIDs(t, mustContent(t, p3+p1+p2), &prev)

	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
	// Exact moved paragraphs keep their IDs.
	assertSameID(t, prev, next, "Dogs are loyal")
	assertSameID(t, prev, next, "Mathematics underpins")
}

func TestBlock_ExactMovedParagraphsKeepIDs(t *testing.T) {
	p1 := "<p>Dogs are loyal companions that enjoy outdoor activities and fetch.</p>"
	p2 := "<p>Mathematics underpins all engineering disciplines across every field.</p>"
	p3 := "<p>Ancient Romans built remarkable infrastructure spanning their vast empire.</p>"
	prev := mustContentWithNewBlocks(t, p1+p2+p3)
	next := mustWithIDs(t, mustContent(t, p3+p1+p2), &prev)

	assertSameID(t, prev, next, "Dogs are loyal")
	assertSameID(t, prev, next, "Mathematics underpins")
	assertSameID(t, prev, next, "Ancient Romans")
}

func TestBlock_ExactMovedContainerKeepsID(t *testing.T) {
	s1 := "<section><p>First section content that is substantial enough for matching.</p></section>"
	s2 := "<section><p>Second section content that is also substantial enough for matching.</p></section>"
	prev := mustContentWithNewBlocks(t, s1+s2)
	next := mustWithIDs(t, mustContent(t, s2+s1), &prev)

	assertSameID(t, prev, next, "First section content")
	assertSameID(t, prev, next, "Second section content")
}

// TestBlock_SplitParagraph verifies that splitting a paragraph produces two valid unique IDs.
// Neither half is > 80% similar to the combined original (Levenshtein ~50%), so both get
// fresh IDs; the assertion is that splitting is safe and produces no duplicates.
func TestBlock_SplitParagraph(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>First half of the paragraph. Second half of the paragraph.</p>")

	next := mustWithIDs(t, mustContent(t, "<p>First half of the paragraph.</p><p>Second half of the paragraph.</p>"), &prev)

	blocks := next.Blocks()
	require.Len(t, blocks, 2)
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

// TestBlock_TypeChangePToH2GetsNewID verifies that changing tag type (p→h2) yields a new ID.
func TestBlock_TypeChangePToH2GetsNewID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>This text will change its element type entirely.</p>")
	next := mustWithIDs(t, mustContent(t, "<h2>This text will change its element type entirely.</h2>"), &prev)

	// Both versions have the same text but different types — must get a new ID.
	assertDiffID(t, prev, next, "change its element type")
}

// TestBlock_ShortParagraphExactMatchKeepsID verifies that short text (< 20 chars) with exact match keeps ID.
func TestBlock_ShortParagraphExactMatchKeepsID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>Short text here!</p>")
	next := mustWithIDs(t, mustContent(t, "<p>Short text here!</p>"), &prev)

	assertSameID(t, prev, next, "Short text here")
}

// TestBlock_ShortParagraphOneCharChangeGetsNewID verifies that short text requires exact match.
func TestBlock_ShortParagraphOneCharChangeGetsNewID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>Short text here!</p>")
	next := mustWithIDs(t, mustContent(t, "<p>Short text here?</p>"), &prev)

	assertDiffID(t, prev, next, "Short text here")
}

// TestBlock_BlockquoteSlightEditKeepsID verifies that a < 10% change in a blockquote keeps the ID.
func TestBlock_BlockquoteSlightEditKeepsID(t *testing.T) {
	// Original and edited share > 90% similarity.
	prev := mustContentWithNewBlocks(t, "<blockquote>To be or not to be that is the question whether tis nobler in the mind to suffer.</blockquote>")
	next := mustWithIDs(t, mustContent(t, "<blockquote>To be or not to be that is the question whether tis nobler in the mind to endure.</blockquote>"), &prev)

	assertSameID(t, prev, next, "To be or not")
}

// TestBlock_BlockquoteModerateEditGetsNewID verifies the 90% threshold for blockquote.
func TestBlock_BlockquoteModerateEditGetsNewID(t *testing.T) {
	// Change roughly 20% of the text — passes 80% threshold but fails 90%.
	orig := "<blockquote>To be or not to be that is the question whether tis nobler in the mind to suffer the slings.</blockquote>"
	edited := "<blockquote>To be or not to be that is the question whether tis nobler a completely different ending here now yes.</blockquote>"

	prev := mustContentWithNewBlocks(t, orig)
	next := mustWithIDs(t, mustContent(t, edited), &prev)

	assertDiffID(t, prev, next, "To be or not to be")
}

// TestBlock_CodeBlockExactMatchKeepsID verifies that an unchanged code block keeps its ID.
func TestBlock_CodeBlockExactMatchKeepsID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<pre><code>func main() {\n    fmt.Println(\"hello\")\n}</code></pre>")
	next := mustWithIDs(t, mustContent(t, "<pre><code>func main() {\n    fmt.Println(\"hello\")\n}</code></pre>"), &prev)

	assertSameID(t, prev, next, "func main")
}

// TestBlock_ListAddItemAtStart verifies that adding a list item at the start keeps existing items' IDs.
func TestBlock_ListAddItemAtStart(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<ul><li>Item alpha with some text here.</li><li>Item beta with some text here.</li><li>Item gamma with some text here.</li></ul>")
	next := mustWithIDs(t, mustContent(t, "<ul><li>New item prepended to the list here.</li><li>Item alpha with some text here.</li><li>Item beta with some text here.</li><li>Item gamma with some text here.</li></ul>"), &prev)

	assertSameIDOfType(t, prev, next, "li", "Item alpha")
	assertSameIDOfType(t, prev, next, "li", "Item beta")
	assertSameIDOfType(t, prev, next, "li", "Item gamma")
}

// TestBlock_ListAddItemAtEnd verifies that adding a list item at the end keeps existing items' IDs.
func TestBlock_ListAddItemAtEnd(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<ul><li>Item alpha with some text here.</li><li>Item beta with some text here.</li><li>Item gamma with some text here.</li></ul>")
	next := mustWithIDs(t, mustContent(t, "<ul><li>Item alpha with some text here.</li><li>Item beta with some text here.</li><li>Item gamma with some text here.</li><li>New item appended to the list here.</li></ul>"), &prev)

	assertSameIDOfType(t, prev, next, "li", "Item alpha")
	assertSameIDOfType(t, prev, next, "li", "Item beta")
	assertSameIDOfType(t, prev, next, "li", "Item gamma")
}

// TestBlock_ListReorderItems verifies reorder behaviour for list items.
func TestBlock_ListReorderItems(t *testing.T) {
	li1 := "<li>Dogs are loyal companions that enjoy outdoor activities everywhere.</li>"
	li2 := "<li>Mathematics underpins engineering disciplines across every field worldwide.</li>"
	li3 := "<li>Ancient Romans built remarkable infrastructure spanning their vast empire.</li>"
	prev := mustContentWithNewBlocks(t, "<ul>"+li1+li2+li3+"</ul>")
	next := mustWithIDs(t, mustContent(t, "<ul>"+li3+li1+li2+"</ul>"), &prev)

	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
	// Exact moved list items keep their IDs.
	assertSameIDOfType(t, prev, next, "li", "Dogs are loyal")
	assertSameIDOfType(t, prev, next, "li", "Mathematics underpins")
}

func TestBlock_ExactMovedListItemsKeepIDs(t *testing.T) {
	li1 := "<li>Dogs are loyal companions that enjoy outdoor activities everywhere.</li>"
	li2 := "<li>Mathematics underpins engineering disciplines across every field worldwide.</li>"
	li3 := "<li>Ancient Romans built remarkable infrastructure spanning their vast empire.</li>"
	prev := mustContentWithNewBlocks(t, "<ul>"+li1+li2+li3+"</ul>")
	next := mustWithIDs(t, mustContent(t, "<ul>"+li3+li1+li2+"</ul>"), &prev)

	assertSameIDOfType(t, prev, next, "li", "Dogs are loyal")
	assertSameIDOfType(t, prev, next, "li", "Mathematics underpins")
	assertSameIDOfType(t, prev, next, "li", "Ancient Romans")
}

// TestBlock_ListDeleteOneItem verifies that deleting one list item keeps the rest's IDs.
func TestBlock_ListDeleteOneItem(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<ul><li>Item alpha with some text here.</li><li>Item beta gets deleted now.</li><li>Item gamma with some text here.</li></ul>")
	next := mustWithIDs(t, mustContent(t, "<ul><li>Item alpha with some text here.</li><li>Item gamma with some text here.</li></ul>"), &prev)

	assertSameIDOfType(t, prev, next, "li", "Item alpha")
	assertSameIDOfType(t, prev, next, "li", "Item gamma")
}

// TestBlock_MixedH2ParagraphDocument verifies that editing a paragraph leaves the heading ID unchanged.
func TestBlock_MixedH2ParagraphDocument(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<h2>Section heading stays the same</h2><p>Paragraph content that will be edited here.</p>")
	next := mustWithIDs(t, mustContent(t, "<h2>Section heading stays the same</h2><p>Paragraph content that has been updated now.</p>"), &prev)

	assertSameID(t, prev, next, "Section heading")
	// The paragraph changed enough to potentially get a new ID, but let's at least
	// confirm the heading kept its ID and both blocks have valid IDs.
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

// TestBlock_FigureUnchangedKeepsID verifies that an unmodified figure keeps its ID.
func TestBlock_FigureUnchangedKeepsID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, `<figure><img src="https://example.com/photo.jpg" alt="A photo"/><figcaption>A caption for the image.</figcaption></figure>`)
	next := mustWithIDs(t, mustContent(t, `<figure><img src="https://example.com/photo.jpg" alt="A photo"/><figcaption>A caption for the image.</figcaption></figure>`), &prev)

	assertSameID(t, prev, next, "caption for the image")
}

func TestBlock_ImageAltEditKeepsID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, `<p>Image follows.</p><img src="https://example.com/photo.jpg" alt="Before">`)
	next := mustWithIDs(t, mustContent(t, `<p>Image follows.</p><img src="https://example.com/photo.jpg" alt="After">`), &prev)

	assert.Equal(t, firstIDOfType(t, prev, "img"), firstIDOfType(t, next, "img"))
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

func TestBlock_ImageSrcEditGetsNewID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, `<p>Image follows.</p><img src="https://example.com/before.jpg" alt="Image">`)
	next := mustWithIDs(t, mustContent(t, `<p>Image follows.</p><img src="https://example.com/after.jpg" alt="Image">`), &prev)

	assert.NotEqual(t, firstIDOfType(t, prev, "img"), firstIDOfType(t, next, "img"))
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

// TestBlock_CJKMinorEditPreservesID verifies that rune-based similarity works for CJK text.
func TestBlock_CJKMinorEditPreservesID(t *testing.T) {
	// Japanese sentence: "The quick brown fox jumps over the lazy dog" roughly.
	// Original and edited differ by one character (one CJK rune).
	orig := "<p>日本語のテキストで短い変更をテストしています。これはサンプル文章です。</p>"
	// Change one character near the end: テスト -> テスタ
	edited := "<p>日本語のテキストで短い変更をテストしています。これはサンプル文章です!</p>"

	prev := mustContentWithNewBlocks(t, orig)
	next := mustWithIDs(t, mustContent(t, edited), &prev)

	assertSameID(t, prev, next, "日本語")
}

func TestBlock_NewRichTextWithChangedBlocksRequiresPreviousContent(t *testing.T) {
	c := mustContentWithNewBlocks(t, "<p>First paragraph of a brand new post.</p><p>Second paragraph of a brand new post.</p>")

	_, err := c.withPreviousState(nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, errPreviousStateRequired)
}

// TestBlock_DuplicateBlockIDRepaired verifies that duplicate id values in submitted HTML
// are repaired so every block ends up with a unique valid ID.
func TestBlock_DuplicateBlockIDRepaired(t *testing.T) {
	knownID := "sdb_cv3j6fld0aaaab44kk50"
	raw := `<p id="` + knownID + `">First block with the same ID.</p><p id="` + knownID + `">Second block with the same ID.</p>`
	result := mustContentWithNewBlocks(t, raw)

	blocks := result.Blocks()
	require.Len(t, blocks, 2)
	assertAllHaveIDs(t, result)
	assertUniqueIDs(t, result)
	assert.NotEqual(t, blocks[0].ID, blocks[1].ID)
}

func TestBlock_DuplicatedCopiedBlockKeepsOnlyOneOldID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>Reusable paragraph that gets copied and pasted in an editor.</p>")
	prevID := prev.Blocks()[0].ID

	nextHTML := `<p id="` + prevID + `">Reusable paragraph that gets copied and pasted in an editor.</p>` +
		`<p id="` + prevID + `">Reusable paragraph that gets copied and pasted in an editor.</p>`
	next := mustWithIDs(t, mustContent(t, nextHTML), &prev)

	kept := 0
	for _, b := range next.Blocks() {
		if b.ID == prevID {
			kept++
		}
	}
	assert.Equal(t, 1, kept)
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

func TestBlock_ChangedBlocksDoesNotMutateInputs(t *testing.T) {
	previous := mustContentWithNewBlocks(t, `<p>Existing paragraph keeps its ID.</p><p>Another existing paragraph keeps its ID.</p>`)
	next := mustContent(t, `<p>Existing paragraph keeps its ID.</p><p>Another existing paragraph keeps its ID with an edit.</p>`)
	previousHTML := previous.HTML()
	nextHTML := next.HTML()

	stable, err := NewRichTextWithChangedBlocks(previous, next)
	require.NoError(t, err)

	assert.Equal(t, previousHTML, previous.HTML())
	assert.Equal(t, nextHTML, next.HTML())
	assert.NotContains(t, next.HTML(), ` id=`)
	assertAllHaveIDs(t, stable.Content)
	assertUniqueIDs(t, stable.Content)
}

func TestBlock_ChangedBlocksPreservesNestedTreeShape(t *testing.T) {
	previous := mustContentWithNewBlocks(t, `<section><blockquote><p>Nested quote text stays the same.</p></blockquote><figure><img src="https://example.com/photo.jpg" alt="Before"></figure></section>`)
	next := mustWithIDs(t, mustContent(t, `<section><blockquote><p>Nested quote text stays the same.</p></blockquote><figure><img src="https://example.com/photo.jpg" alt="After"></figure></section>`), &previous)

	assert.Equal(t, []string{"section", "blockquote", "p", "figure", "img"}, blockTypes(next.Blocks()))
	assertSameIDOfType(t, previous, next, "p", "Nested quote")
	assert.Equal(t, firstIDOfType(t, previous, "img"), firstIDOfType(t, next, "img"))
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

func TestBlock_MinorEditWithTrustedIncomingIDKeepsID(t *testing.T) {
	prev := mustContentWithNewBlocks(t, "<p>This paragraph has enough words to survive a small trusted editor round trip.</p>")
	prevID := prev.Blocks()[0].ID

	next := mustWithIDs(t, mustContent(t, `<p id="`+prevID+`">This paragraph has enough words to survive one small trusted editor round trip.</p>`), &prev)

	assertSameID(t, prev, next, "This paragraph has enough")
}

func TestBlock_DuplicatePreviousIDsAreNotReused(t *testing.T) {
	knownID := "sdb_cv3j6fld0aaaab44kk50"
	prev := Content{
		html: mustParseBody(t, `<p id="`+knownID+`">First previous block with duplicate ID.</p><p id="`+knownID+`">Second previous block with duplicate ID.</p>`),
	}
	next := mustWithIDs(t, mustContent(t, `<p>First previous block with duplicate ID.</p><p>Second previous block with duplicate ID.</p>`), &prev)

	for _, b := range next.Blocks() {
		assert.NotEqual(t, knownID, b.ID)
	}
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

func TestBlock_AmbiguousIdenticalBlocksOnlyReuseDisambiguatedContext(t *testing.T) {
	repeated := "<p>Repeated paragraph with identical wording that cannot identify a specific original.</p>"
	anchor := "<p>Anchor paragraph that changes the surrounding order for context.</p>"
	prev := mustContentWithNewBlocks(t, anchor+repeated+repeated)
	next := mustWithIDs(t, mustContent(t, repeated+anchor+repeated), &prev)

	prevRepeatedIDs := idsOfType(prev, "p", "Repeated paragraph")
	nextRepeatedIDs := idsOfType(next, "p", "Repeated paragraph")
	require.Len(t, prevRepeatedIDs, 2)
	require.Len(t, nextRepeatedIDs, 2)

	assert.NotEqual(t, prevRepeatedIDs[0], nextRepeatedIDs[0], "ambiguous moved duplicate should not arbitrarily reuse the first old ID")
	assert.Equal(t, prevRepeatedIDs[1], nextRepeatedIDs[1], "same-index duplicate is disambiguated by context and may keep its ID")
	assertAllHaveIDs(t, next)
	assertUniqueIDs(t, next)
}

func mustParseBody(t *testing.T, raw string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(raw))
	require.NoError(t, err)
	body := findBody(doc)
	require.NotNil(t, body)
	return body
}

func blockTypes(blocks []Block) []string {
	out := make([]string, len(blocks))
	for i, b := range blocks {
		out[i] = b.Type
	}
	return out
}

func firstIDOfType(t *testing.T, c Content, typ string) string {
	t.Helper()
	for _, b := range c.Blocks() {
		if b.Type == typ {
			return b.ID
		}
	}
	t.Fatalf("no %s block found", typ)
	return ""
}
