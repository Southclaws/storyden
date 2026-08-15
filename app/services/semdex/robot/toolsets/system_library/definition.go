package system_library

const (
	ID          = "system.library"
	Name        = "Library management"
	Description = "Search, read, create, and maintain structured Storyden Library pages and their property schemas."
	Instruction = `Use the Library as the source of truth for durable, structured knowledge. Search or ask the user to select a page when the target is ambiguous. Read the current page or schema before changing it, preserve unrelated content and properties, and treat page, property, and schema mutations as side effects. Schema updates can remove omitted fields, so inspect first and change only what the task requires. Report the page affected and the substantive result.`
)
