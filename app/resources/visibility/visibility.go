package visibility

import (
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/post"
)

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type visibilityEnum string

const (
	visibilityDraft     visibilityEnum = "draft"
	visibilityUnlisted  visibilityEnum = "unlisted"
	visibilityReview    visibilityEnum = "review"
	visibilityPublished visibilityEnum = "published"
)

func NewVisibilityFromEnt[T post.Visibility | node.Visibility | collection.Visibility](in T) Visibility {
	return Visibility{visibilityEnum(in)}
}
