package datagraph

import (
	"github.com/Southclaws/fault"
	"github.com/rs/xid"
	"golang.org/x/net/html"
)

func newBlockID() string {
	return blockIDPrefix + xid.New().String()
}

var errPreviousStateRequired = fault.New("previous content state is required")

// Blocks returns all addressable block-level elements in document order.
// Each Block carries the current id (empty if not yet assigned),
// the HTML tag type, plain text, and raw outer HTML.
func (c Content) Blocks() []Block {
	if c.html == nil {
		return nil
	}
	var blocks []Block
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}
		if isBlockNode(n) {
			blocks = append(blocks, Block{
				ID:   getAttr(n, blockIDAttributeName),
				Type: n.Data,
				Text: textFromNode(n, n.Data == "pre"),
				HTML: renderNode(n),
			})
		}
		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			walk(ch)
		}
	}
	walk(c.html)
	return blocks
}

// withBlockIDs returns a Content where every block-level element carries a
// valid id attribute. Existing valid IDs are preserved; blocks without
// a valid ID get a fresh XID. This is the no-comparison path used by
// NewRichTextWithNewBlocks on initial creation.
func (c Content) withBlockIDs() (Content, error) {
	if c.html == nil {
		return c, nil
	}

	var nodes []contentBlockNode
	collectBlockNodes(c.html, &nodes)
	if len(nodes) == 0 {
		return c, nil
	}

	idCount := countValidIDs(nodes)
	allValidUnique := true
	for _, n := range nodes {
		if !isValidBlockID(n.id) || idCount[n.id] > 1 {
			allValidUnique = false
			break
		}
	}
	if allValidUnique {
		return c, nil
	}

	freshBody := cloneNodeDeep(c.html)

	var freshNodes []contentBlockNode
	collectBlockNodes(freshBody, &freshNodes)
	ids := assignBlockIDs(nil, freshNodes)
	for i, nn := range freshNodes {
		setBlockID(nn.n, ids[i])
	}

	// Replace only the html field; short, plain, links, and sdrs are not
	// affected by block ID attributes.
	result := c
	result.html = freshBody
	return result, nil
}

// withPreviousState returns a new Content where every block-level element
// carries an id attribute, preserving IDs from the previous version where
// possible using scored one-to-one reconciliation.
func (c Content) withPreviousState(previous *Content) (Content, error) {
	if previous == nil || previous.html == nil {
		return c, errPreviousStateRequired
	}

	// Collect old block nodes from previous content in a single walk.
	var oldNodes []contentBlockNode
	collectBlockNodes(previous.html, &oldNodes)

	body := cloneNodeDeep(c.html)

	var newNodes []contentBlockNode
	collectBlockNodes(body, &newNodes)

	if len(newNodes) == 0 {
		return c, nil
	}

	ids := assignBlockIDs(oldNodes, newNodes)

	for i, nn := range newNodes {
		setBlockID(nn.n, ids[i])
	}

	result := c
	result.html = body
	return result, nil
}
