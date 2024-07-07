package schema

var (
	VisibilityTypesDraft     = "draft"     // Items in draft are only accessible by the owner.
	VisibilityTypesUnlisted  = "unlisted"  // Items unlisted are not publicly listed but accessible via collections and other direct links.
	VisibilityTypesReview    = "review"    // Items in review are only accessible by the owner and admins.
	VisibilityTypesPublished = "published" // Items published are published globally and searchable.
)

var VisibilityTypes = []string{
	VisibilityTypesDraft,
	VisibilityTypesUnlisted,
	VisibilityTypesReview,
	VisibilityTypesPublished,
}
