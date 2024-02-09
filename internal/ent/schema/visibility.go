package schema

var (
	VisibilityTypesDraft  = "draft"  // Items in draft are only accessible by the owner.
	VisibilityTypesReview = "review" // Items in review are only accessible by the owner and admins.
	VisibilityTypesPublic = "public" // Items in public are accessible by everyone.
)

var VisibilityTypes = []string{
	VisibilityTypesDraft,
	VisibilityTypesReview,
	VisibilityTypesPublic,
}
