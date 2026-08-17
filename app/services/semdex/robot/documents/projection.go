package documents

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	adksession "google.golang.org/adk/v2/session"
)

const (
	tinyDocumentBytes  = 1000
	projectionRunes    = 6000
	leafRunes          = 8000
	previewRunes       = 180
	searchPreviewRunes = 260
	structurePageNodes = 25
)

type Projection struct {
	DocumentID string
	SourceType SourceType
	SourceID   string
	Title      string
	NodeID     string
	Page       int
	TotalPages int
	Previous   *int
	Next       *int
	ItemStart  int
	ItemEnd    int
	TotalItems int
	Text       string
	Truncated  bool
	NextAction string
}

type SearchMatch struct {
	NodeID  string
	Kind    string
	Path    string
	Preview string
}

type rankedSearchMatch struct {
	match SearchMatch
	score int
	order int
}

type documentNode struct {
	ID       string
	Kind     string
	Text     string
	Level    int
	Parent   *documentNode
	Children []*documentNode
}

func project(snapshot Snapshot, nodeID string, page int) (Projection, error) {
	content, root, err := parseSnapshot(snapshot)
	if err != nil {
		return Projection{}, err
	}

	node := root
	if nodeID != RootNodeID {
		node = findNode(root, nodeID)
		if node == nil {
			return Projection{}, fmt.Errorf("%w: %s", ErrNodeNotFound, nodeID)
		}
	}
	if page <= 0 {
		page = 1
	}

	result := Projection{
		DocumentID: snapshot.DocumentID,
		SourceType: snapshot.SourceType,
		SourceID:   snapshot.SourceID,
		Title:      snapshot.Title,
		NodeID:     node.ID,
		Page:       page,
		TotalPages: 1,
	}

	if node == root && len([]byte(content.Plaintext())) <= tinyDocumentBytes {
		if page != 1 {
			return Projection{}, invalidPage(page, 1)
		}
		result.Text = renderCompleteDocument(snapshot.Title, root)
		result.NextAction = "The complete document is shown; use document_search only if a focused lookup would help."
		return result, nil
	}

	if node == root {
		result.Text, result.Truncated, result.TotalPages, result.ItemStart, result.ItemEnd, result.TotalItems, err = renderRoot(snapshot.Title, root, page)
		if err != nil {
			return Projection{}, err
		}
		setPageLinks(&result)
		result.NextAction = projectionNextAction(result, fmt.Sprintf("Call document_get with document_id %q and a listed node_id to inspect that location, or use document_search for a focused lookup.", snapshot.DocumentID))
		return result, nil
	}

	if node.Kind == "heading" {
		result.Text, result.Truncated, result.TotalPages, result.ItemStart, result.ItemEnd, result.TotalItems, err = renderHeading(node, page)
		if err != nil {
			return Projection{}, err
		}
		setPageLinks(&result)
		result.NextAction = projectionNextAction(result, fmt.Sprintf("Call document_get with document_id %q and a listed child node_id for complete leaf content, or search within node_id %q.", snapshot.DocumentID, node.ID))
		return result, nil
	}
	if page != 1 {
		return Projection{}, invalidPage(page, 1)
	}

	result.Text, result.Truncated = truncateRunes(strings.TrimSpace(node.Text), leafRunes)
	if result.Truncated {
		result.NextAction = fmt.Sprintf("Use document_search scoped to node_id %q with more specific terms to locate the needed passage.", node.ID)
	} else {
		result.NextAction = "This structural location is shown completely; return to the document root or search elsewhere if more context is needed."
	}
	return result, nil
}

func setPageLinks(result *Projection) {
	if result.Page > 1 {
		previous := result.Page - 1
		result.Previous = &previous
	}
	if result.Page < result.TotalPages {
		next := result.Page + 1
		result.Next = &next
	}
}

func projectionNextAction(result Projection, fallback string) string {
	if result.Next != nil {
		return fmt.Sprintf("Call document_get with document_id %q, node_id %q, and page %d to continue this location.", result.DocumentID, result.NodeID, *result.Next)
	}
	if result.Previous != nil {
		return fmt.Sprintf("This is the final structural page. Inspect a listed node, search this document, or call document_get with page %d to move backward.", *result.Previous)
	}
	return fallback
}

func invalidPage(page, totalPages int) error {
	return fmt.Errorf("document page %d is out of range; this location has %d page(s)", page, totalPages)
}

func Search(state adksession.ReadonlyState, documentID, query string, scopeIDs []string, maxResults int) (string, string, []SearchMatch, bool, error) {
	snapshot, err := resolveSnapshot(state, documentID)
	if err != nil {
		return "", "", nil, false, err
	}
	_, root, err := parseSnapshot(snapshot)
	if err != nil {
		return "", "", nil, false, err
	}

	terms := strings.Fields(strings.ToLower(strings.TrimSpace(query)))
	if len(terms) == 0 {
		return "", "", nil, false, fmt.Errorf("search query is required")
	}
	if maxResults <= 0 {
		maxResults = 10
	}

	allowed := map[string]struct{}{}
	if len(scopeIDs) > 0 {
		for _, scopeID := range scopeIDs {
			scope := root
			if scopeID != RootNodeID {
				scope = findNode(root, scopeID)
			}
			if scope == nil {
				return "", "", nil, false, fmt.Errorf("%w: %s", ErrNodeNotFound, scopeID)
			}
			walkNodes(scope, func(node *documentNode) {
				allowed[node.ID] = struct{}{}
			})
		}
	}

	ranked := make([]rankedSearchMatch, 0, maxResults)
	order := 0
	walkNodes(root, func(node *documentNode) {
		if node == root || strings.TrimSpace(node.Text) == "" {
			return
		}
		if len(allowed) > 0 {
			if _, ok := allowed[node.ID]; !ok {
				return
			}
		}

		lower := strings.ToLower(node.Text)
		for _, term := range terms {
			if !strings.Contains(lower, term) {
				return
			}
		}

		score := 0
		for _, term := range terms {
			if containsWholeTerm(lower, term) {
				score++
			}
		}
		ranked = append(ranked, rankedSearchMatch{
			match: SearchMatch{
				NodeID:  node.ID,
				Kind:    node.Kind,
				Path:    nodePath(snapshot.Title, node),
				Preview: matchingPreview(node.Text, terms[0]),
			},
			score: score,
			order: order,
		})
		order++
	})

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].order < ranked[j].order
	})
	total := len(ranked)
	if len(ranked) > maxResults {
		ranked = ranked[:maxResults]
	}
	matches := make([]SearchMatch, 0, len(ranked))
	for _, candidate := range ranked {
		matches = append(matches, candidate.match)
	}

	return snapshot.DocumentID, strings.Join(terms, " "), matches, total > len(matches), nil
}

func containsWholeTerm(text, term string) bool {
	if term == "" {
		return false
	}
	first, _ := utf8.DecodeRuneInString(term)
	last, _ := utf8.DecodeLastRuneInString(term)
	if !isWordRune(first) || !isWordRune(last) {
		return false
	}

	for offset := 0; offset <= len(text)-len(term); {
		index := strings.Index(text[offset:], term)
		if index < 0 {
			return false
		}
		start := offset + index
		end := start + len(term)
		leftBoundary := start == 0
		if !leftBoundary {
			previous, _ := utf8.DecodeLastRuneInString(text[:start])
			leftBoundary = !isWordRune(previous)
		}
		rightBoundary := end == len(text)
		if !rightBoundary {
			next, _ := utf8.DecodeRuneInString(text[end:])
			rightBoundary = !isWordRune(next)
		}
		if leftBoundary && rightBoundary {
			return true
		}
		offset = start + len(term)
	}
	return false
}

func isWordRune(value rune) bool {
	return value == '_' || unicode.IsLetter(value) || unicode.IsNumber(value) || unicode.IsMark(value)
}

func parseSnapshot(snapshot Snapshot) (datagraph.Content, *documentNode, error) {
	parsed, err := datagraph.NewRichTextWithBlocks(snapshot.Content)
	if err != nil {
		return datagraph.Content{}, nil, fmt.Errorf("parse document snapshot: %w", err)
	}

	root := &documentNode{ID: RootNodeID, Kind: "document", Text: snapshot.Title}
	headings := []*documentNode{root}
	for _, block := range parsed.Blocks() {
		kind, level, include := classifyBlock(block.Type)
		if !include || block.ID == "" {
			continue
		}
		text := strings.TrimSpace(block.Text)
		if kind == "image" && text == "" {
			text = block.HTML
		}
		node := &documentNode{ID: block.ID, Kind: kind, Text: text, Level: level}
		if kind == "heading" {
			for len(headings) > 1 && headings[len(headings)-1].Level >= level {
				headings = headings[:len(headings)-1]
			}
			parent := headings[len(headings)-1]
			node.Parent = parent
			parent.Children = append(parent.Children, node)
			headings = append(headings, node)
			continue
		}

		parent := headings[len(headings)-1]
		node.Parent = parent
		parent.Children = append(parent.Children, node)
	}
	return parsed.Content, root, nil
}

func classifyBlock(typ string) (kind string, level int, include bool) {
	switch typ {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return "heading", int(typ[1] - '0'), true
	case "p":
		return "paragraph", 0, true
	case "li":
		return "list_item", 0, true
	case "pre":
		return "code", 0, true
	case "table":
		return "table", 0, true
	case "figure":
		return "figure", 0, true
	case "img":
		return "image", 0, true
	default:
		return "", 0, false
	}
}

func renderRoot(title string, root *documentNode, page int) (string, bool, int, int, int, int, error) {
	var lines []string
	headings := directAndNestedHeadings(root)
	entries := headings
	if len(entries) == 0 {
		entries = root.Children
	}
	paged, totalPages, itemStart, itemEnd, totalItems, err := structuralPage(entries, page)
	if err != nil {
		return "", false, totalPages, 0, 0, totalItems, err
	}
	lines = append(lines, "# "+title+"\n\n"+structuralPageSummary(page, totalPages, itemStart, itemEnd, totalItems))
	if len(headings) > 0 {
		for _, heading := range paged {
			depth := headingDepth(heading)
			lines = append(lines, strings.Repeat("  ", max(0, depth-1))+"- ["+heading.ID+"] "+heading.Text)
		}
	} else {
		for _, child := range paged {
			preview, _ := truncateRunes(child.Text, previewRunes)
			lines = append(lines, fmt.Sprintf("- [%s] (%s) %s", child.ID, child.Kind, preview))
		}
	}
	text, clipped := truncateRunes(strings.Join(lines, "\n"), projectionRunes)
	return text, clipped || len(root.Children) > 0 || totalPages > 1, totalPages, itemStart, itemEnd, totalItems, nil
}

func renderHeading(node *documentNode, page int) (string, bool, int, int, int, int, error) {
	children, totalPages, itemStart, itemEnd, totalItems, err := structuralPage(node.Children, page)
	if err != nil {
		return "", false, totalPages, 0, 0, totalItems, err
	}
	lines := []string{
		strings.Repeat("#", max(1, node.Level)) + " " + node.Text,
		structuralPageSummary(page, totalPages, itemStart, itemEnd, totalItems),
	}
	previewClipped := false
	for _, child := range children {
		if child.Kind == "heading" {
			lines = append(lines, fmt.Sprintf("- [%s] %s", child.ID, child.Text))
			continue
		}
		preview, clipped := truncateRunes(child.Text, previewRunes)
		previewClipped = previewClipped || clipped
		lines = append(lines, fmt.Sprintf("- [%s] (%s) %s", child.ID, child.Kind, preview))
	}
	text, clipped := truncateRunes(strings.Join(lines, "\n"), projectionRunes)
	return text, clipped || previewClipped || hasNestedContent(node) || totalPages > 1, totalPages, itemStart, itemEnd, totalItems, nil
}

func renderCompleteDocument(title string, root *documentNode) string {
	lines := []string{"# " + title}
	walkNodes(root, func(node *documentNode) {
		if node == root || strings.TrimSpace(node.Text) == "" {
			return
		}
		if node.Kind == "heading" {
			lines = append(lines, strings.Repeat("#", max(1, node.Level))+" "+node.Text)
			return
		}
		lines = append(lines, node.Text)
	})
	return strings.Join(lines, "\n\n")
}

func structuralPageSummary(page, totalPages, itemStart, itemEnd, totalItems int) string {
	if totalItems == 0 {
		return fmt.Sprintf("Structural page %d of %d; no child items", page, totalPages)
	}
	return fmt.Sprintf("Structural page %d of %d; items %d-%d of %d", page, totalPages, itemStart, itemEnd, totalItems)
}

func structuralPage(nodes []*documentNode, page int) ([]*documentNode, int, int, int, int, error) {
	totalPages := max(1, (len(nodes)+structurePageNodes-1)/structurePageNodes)
	if page < 1 || page > totalPages {
		return nil, totalPages, 0, 0, len(nodes), invalidPage(page, totalPages)
	}
	start := (page - 1) * structurePageNodes
	end := min(len(nodes), start+structurePageNodes)
	if len(nodes) == 0 {
		return nodes[start:end], totalPages, 0, 0, 0, nil
	}
	return nodes[start:end], totalPages, start + 1, end, len(nodes), nil
}

func directAndNestedHeadings(root *documentNode) []*documentNode {
	var headings []*documentNode
	walkNodes(root, func(node *documentNode) {
		if node.Kind == "heading" {
			headings = append(headings, node)
		}
	})
	return headings
}

func headingDepth(node *documentNode) int {
	depth := 0
	for current := node.Parent; current != nil && current.Kind != "document"; current = current.Parent {
		if current.Kind == "heading" {
			depth++
		}
	}
	return depth
}

func hasNestedContent(node *documentNode) bool {
	for _, child := range node.Children {
		if child.Kind == "heading" && len(child.Children) > 0 {
			return true
		}
	}
	return false
}

func findNode(root *documentNode, id string) *documentNode {
	var found *documentNode
	walkNodes(root, func(node *documentNode) {
		if found == nil && node.ID == id {
			found = node
		}
	})
	return found
}

func walkNodes(node *documentNode, visit func(*documentNode)) {
	if node == nil {
		return
	}
	visit(node)
	for _, child := range node.Children {
		walkNodes(child, visit)
	}
}

func nodePath(title string, node *documentNode) string {
	parts := []string{title}
	for current := node; current != nil && current.Kind != "document"; current = current.Parent {
		if current.Kind == "heading" {
			parts = append(parts, current.Text)
		}
	}
	slices.Reverse(parts[1:])
	return strings.Join(parts, " > ")
}

func matchingPreview(text, term string) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) <= searchPreviewRunes {
		return text
	}
	lower := strings.ToLower(text)
	byteIndex := strings.Index(lower, term)
	if byteIndex < 0 {
		preview, _ := truncateRunes(text, searchPreviewRunes)
		return preview
	}
	runeIndex := utf8.RuneCountInString(lower[:byteIndex])
	start := max(0, runeIndex-searchPreviewRunes/3)
	end := min(len(runes), start+searchPreviewRunes)
	preview := strings.TrimSpace(string(runes[start:end]))
	if start > 0 {
		preview = "…" + preview
	}
	if end < len(runes) {
		preview += "…"
	}
	return preview
}

func truncateRunes(text string, limit int) (string, bool) {
	runes := []rune(text)
	if len(runes) <= limit {
		return text, false
	}
	return strings.TrimSpace(string(runes[:limit])) + "…", true
}
