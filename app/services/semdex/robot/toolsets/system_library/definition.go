package system_library

const (
	ID          = "system.library"
	Name        = "Library management"
	Description = "Search, read, create, and maintain structured Storyden Library pages and their property schemas."
	Instruction = `Use the Library as the source of truth for durable, structured knowledge. Search or ask the user to select a page when the target is ambiguous. Read the current page or schema before changing it, preserve unrelated content and properties, and treat page, property, and schema mutations as side effects. Schema updates can remove omitted fields, so inspect first and change only what the task requires. Use library_pages_create for multiple pages and use ref with parent_ref to describe their hierarchy in one request. Use library_pages_update for multiple existing pages. Each batch commits all page and tag changes or none; on error, correct the reported item and resubmit the whole batch. Link enrichment occurs after page creation and can fail independently. Batches do not deduplicate repeated requests, so inspect the Library before resubmitting after a lost response. Tag names create missing tags automatically; updates replace the tag set, omission preserves it, and an empty array clears it. Report the pages affected and the substantive result.`
)
