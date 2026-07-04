package datagraph

import (
	"sort"
	"strings"

	"github.com/agext/levenshtein"
	"github.com/rs/xid"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	// short/long branching thresholds
	shortBlockRunes = 20  // Below this length, require exact text/HTML to avoid overmatching tiny blocks.
	longBlockRunes  = 160 // At or above this length, use trigram similarity instead of Levenshtein.

	// trigram and distance
	trigramRuneWidth    = 3 // Number of runes per gram for long-text Dice similarity.
	nearbyIndexDistance = 2 // Maximum document-order distance that earns a small context boost.

	// scoring thresholds
	ambiguityMargin = 0.005 // Scores within this margin are treated as ties unless ID/index context separates them.
	minLengthRatio  = 0.33  // Reject fuzzy matches when one text is less than this fraction of the other.

	// similarity thresholds
	defaultSimilarityMinimum = 0.80 // Minimum fuzzy text similarity for ordinary block types.
	strictSimilarityMinimum  = 0.90 // Minimum fuzzy text similarity for blocks where false positives are especially costly.

	// matching
	exactHTMLMatchScore = 1.0  // Base score for identical normalized HTML after stripping Storyden block IDs.
	exactTextMatchScore = 0.96 // Base score for identical normalized text when markup changed.
	imageSrcMatchScore  = 0.99 // Base score for identical image source when non-identity attributes changed.

	// scoring nodes, boosts
	incomingIDMatchBoost  = 0.20 // Confidence boost when a valid unique incoming ID matches the old block ID.
	sameIndexMatchBoost   = 0.03 // Context boost for blocks that stayed at the same projected index.
	nearbyIndexMatchBoost = 0.01 // Smaller context boost for blocks that moved only slightly.
)

type contentBlockNode struct {
	n        *html.Node
	id       string
	typ      string
	index    int
	normText string
	normHTML string
	identity string
}

func collectBlockNodes(n *html.Node, out *[]contentBlockNode) {
	if n == nil {
		return
	}
	if isBlockNode(n) {
		*out = append(*out, contentBlockNode{
			n:        n,
			id:       getAttr(n, blockIDAttributeName),
			typ:      n.Data,
			index:    len(*out),
			normText: normText(textFromNode(n, n.Data == "pre")),
			normHTML: renderNodeWithoutBlockID(n),
			identity: blockIdentity(n),
		})
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectBlockNodes(c, out)
	}
}

func blockIdentity(n *html.Node) string {
	if n.DataAtom != atom.Img {
		return ""
	}
	return strings.TrimSpace(getAttr(n, "src"))
}

func isValidBlockID(id string) bool {
	if !strings.HasPrefix(id, blockIDPrefix) {
		return false
	}
	_, err := xid.FromString(strings.TrimPrefix(id, blockIDPrefix))
	return err == nil
}

func normText(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return spaces.ReplaceAllString(s, " ")
}

func similarityThreshold(typ string) float64 {
	switch typ {
	case "blockquote", "pre", "table":
		return strictSimilarityMinimum
	default:
		return defaultSimilarityMinimum
	}
}

type blockCandidate struct {
	oldIdx          int
	newIdx          int
	score           float64
	incomingIDMatch bool
	sameIndex       bool
}

func assignBlockIDs(oldNodes, newNodes []contentBlockNode) []string {
	result := make([]string, len(newNodes))

	if len(newNodes) == 0 {
		return result
	}

	newIDCount := countValidIDs(newNodes)
	if len(oldNodes) == 0 {
		for i, nn := range newNodes {
			if isValidBlockID(nn.id) && newIDCount[nn.id] == 1 {
				result[i] = nn.id
			} else {
				result[i] = newBlockID()
			}
		}
		return result
	}

	oldIDCount := countValidIDs(oldNodes)
	oldValidIDs := make(map[string]bool, len(oldIDCount))
	for id := range oldIDCount {
		oldValidIDs[id] = true
	}

	reusableOld := make(map[int]bool, len(oldNodes))
	for i, ob := range oldNodes {
		reusableOld[i] = isValidBlockID(ob.id) && oldIDCount[ob.id] == 1
	}

	candidates := buildBlockCandidates(oldNodes, newNodes, reusableOld, newIDCount)
	candidates = rejectAmbiguousCandidates(candidates)
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			if candidates[i].newIdx == candidates[j].newIdx {
				return candidates[i].oldIdx < candidates[j].oldIdx
			}
			return candidates[i].newIdx < candidates[j].newIdx
		}
		return candidates[i].score > candidates[j].score
	})

	usedOld := make(map[int]bool, len(oldNodes))
	usedNew := make(map[int]bool, len(newNodes))
	for _, c := range candidates {
		if usedOld[c.oldIdx] || usedNew[c.newIdx] {
			continue
		}
		result[c.newIdx] = oldNodes[c.oldIdx].id
		usedOld[c.oldIdx] = true
		usedNew[c.newIdx] = true
	}

	usedIDs := make(map[string]bool, len(newNodes))
	for _, id := range result {
		if id != "" {
			usedIDs[id] = true
		}
	}
	for i, nn := range newNodes {
		if result[i] != "" {
			continue
		}
		if isValidBlockID(nn.id) && newIDCount[nn.id] == 1 && !oldValidIDs[nn.id] && !usedIDs[nn.id] {
			result[i] = nn.id
		} else {
			result[i] = newBlockID()
		}
		usedIDs[result[i]] = true
	}

	return result
}

func countValidIDs(nodes []contentBlockNode) map[string]int {
	counts := make(map[string]int, len(nodes))
	for _, n := range nodes {
		if isValidBlockID(n.id) {
			counts[n.id]++
		}
	}
	return counts
}

func buildBlockCandidates(oldNodes, newNodes []contentBlockNode, reusableOld map[int]bool, newIDCount map[string]int) []blockCandidate {
	candidates := []blockCandidate{}
	for oi, old := range oldNodes {
		if !reusableOld[oi] {
			continue
		}
		for ni, next := range newNodes {
			uniqueIncomingID := newIDCount[next.id] == 1
			score, ok := blockMatchScore(old, next, uniqueIncomingID)
			if !ok {
				continue
			}
			candidates = append(candidates, blockCandidate{
				oldIdx:          oi,
				newIdx:          ni,
				score:           score,
				incomingIDMatch: uniqueIncomingID && next.id == old.id,
				sameIndex:       old.index == next.index,
			})
		}
	}
	return candidates
}

func rejectAmbiguousCandidates(candidates []blockCandidate) []blockCandidate {
	if len(candidates) < 2 {
		return candidates
	}

	ambiguous := make([]bool, len(candidates))
	for i, c := range candidates {
		for j, other := range candidates {
			if i == j {
				continue
			}
			if c.oldIdx != other.oldIdx && c.newIdx != other.newIdx {
				continue
			}
			if c.clearlyBeats(other) || other.clearlyBeats(c) {
				continue
			}
			if scoreDistance(c.score, other.score) <= ambiguityMargin {
				ambiguous[i] = true
				break
			}
		}
	}

	filtered := candidates[:0]
	for i, c := range candidates {
		if !ambiguous[i] {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func (c blockCandidate) clearlyBeats(other blockCandidate) bool {
	if c.incomingIDMatch != other.incomingIDMatch {
		return c.incomingIDMatch
	}
	if c.sameIndex != other.sameIndex {
		return c.sameIndex
	}
	return c.score > other.score+ambiguityMargin
}

func scoreDistance(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func blockMatchScore(old, next contentBlockNode, uniqueIncomingID bool) (float64, bool) {
	if old.typ != next.typ {
		return 0, false
	}

	exactHTML := old.normHTML == next.normHTML
	exactText := old.normText != "" && old.normText == next.normText

	score := 0.0
	switch {
	case exactHTML:
		score = exactHTMLMatchScore
	case exactText:
		score = exactTextMatchScore
	case old.typ == "img" && old.identity != "" && old.identity == next.identity:
		score = imageSrcMatchScore
	default:
		sim, ok := blockTextSimilarity(old, next)
		if !ok || sim < similarityThreshold(old.typ) {
			return 0, false
		}
		score = sim
	}

	if uniqueIncomingID && next.id == old.id {
		score += incomingIDMatchBoost
	}
	if old.index == next.index {
		score += sameIndexMatchBoost
	} else {
		distance := old.index - next.index
		if distance < 0 {
			distance = -distance
		}
		if distance <= nearbyIndexDistance {
			score += nearbyIndexMatchBoost
		}
	}

	return score, true
}

func blockTextSimilarity(old, next contentBlockNode) (float64, bool) {
	oldText, nextText := old.normText, next.normText
	oldLen, nextLen := len([]rune(oldText)), len([]rune(nextText))

	if oldLen < shortBlockRunes || nextLen < shortBlockRunes {
		return 0, false
	}

	lo, hi := oldLen, nextLen
	if lo > hi {
		lo, hi = hi, lo
	}
	if hi > 0 && float64(lo)/float64(hi) < minLengthRatio {
		return 0, false
	}

	if oldLen >= longBlockRunes || nextLen >= longBlockRunes {
		return trigramDiceSimilarity(oldText, nextText), true
	}

	return levenshtein.Similarity(oldText, nextText, nil), true
}

func trigramDiceSimilarity(a, b string) float64 {
	aTrigrams := trigrams(a)
	bTrigrams := trigrams(b)
	if len(aTrigrams) == 0 || len(bTrigrams) == 0 {
		return 0
	}

	intersection := 0
	for k, av := range aTrigrams {
		if bv, ok := bTrigrams[k]; ok {
			if av < bv {
				intersection += av
			} else {
				intersection += bv
			}
		}
	}

	total := 0
	for _, v := range aTrigrams {
		total += v
	}
	for _, v := range bTrigrams {
		total += v
	}
	if total == 0 {
		return 0
	}

	return (2 * float64(intersection)) / float64(total)
}

func trigrams(s string) map[string]int {
	runes := []rune(s)
	if len(runes) < trigramRuneWidth {
		return nil
	}

	out := make(map[string]int, len(runes)-trigramRuneWidth+1)
	for i := 0; i <= len(runes)-trigramRuneWidth; i++ {
		out[string(runes[i:i+trigramRuneWidth])]++
	}
	return out
}
