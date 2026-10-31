package system_library

const (
	ID          = "system.library"
	Name        = "Library management"
	Description = "Search, read, create, and maintain structured Storyden Library pages and their property schemas."
	Instruction = `Use the Library for durable, structured community knowledge.

Use a supplied page identifier for an explicit target. For references such as "this page", use datagraph_item in <invocation-context-json> when its kind is "node": id is the Library page ID and slug is its slug. Resolve the target before searching by keywords. Use library_page_get for metadata or library_page_open for content. If search finds no match, use a suitable direct lookup or visibility-filtered library_page_list before concluding that the page cannot be found. Request a page selection when the target remains ambiguous.

Read the current content or schema when needed to preserve unrelated material. Schema updates can remove omitted fields, so inspect the schema before changing it.

Use library_pages_create for multiple new pages and library_pages_update for multiple existing pages. Use parent_ref for parent-child relationships within a creation batch. Follow the tool descriptions for transaction, tag replacement, and retry behaviour.

Providing a URL requests best-effort link enrichment. Page creation does not establish that enrichment succeeded; report it as confirmed only when the available result establishes it.

Report the affected pages and the substantive changes.`
)
