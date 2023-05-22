package thread

import "github.com/Southclaws/storyden/internal/ent/post"

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type statusEnum string

const (
	statusDraft     statusEnum = "draft"
	statusPublished statusEnum = "published"
)

func NewStatusFromEnt(in post.Status) Status {
	return Status{statusEnum(in)}
}

func (s Status) ToEnt() post.Status {
	return post.Status(s.v)
}
