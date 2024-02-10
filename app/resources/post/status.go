package post

import "github.com/Southclaws/storyden/internal/ent/post"

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type visibilityEnum string

const (
	visibilityDraft     visibilityEnum = "draft"
	visibilityReview    visibilityEnum = "review"
	visibilityPublished visibilityEnum = "published"
)

func NewVisibilityFromEnt(in post.Visibility) Visibility {
	return Visibility{visibilityEnum(in)}
}

func (s Visibility) ToEnt() post.Visibility {
	return post.Visibility(s.v)
}
