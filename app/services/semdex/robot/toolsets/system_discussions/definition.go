package system_discussions

const (
	ID          = "system.discussions"
	Name        = "Discussion management"
	Description = "Browse, read, create, and maintain Storyden discussion threads, replies, categories, and tags."
	Instruction = `Use discussion data as the source of truth for forum activity. Read the current thread before updating or replying when existing context matters. Choose categories and tags from their listing tools rather than inventing identifiers. Treat thread creation, updates, and replies as user-visible side effects; preserve unrelated fields and report exactly what was created or changed.`
)
