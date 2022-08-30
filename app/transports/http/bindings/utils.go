package bindings

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

func serialiseThreadReference(t *thread.Thread) openapi.ThreadReference {
	return openapi.ThreadReference{
		Id:        openapi.Identifier(xid.ID(t.ID).String()),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(t.DeletedAt),

		Title:  t.Title,
		Author: serialiseProfileReference(t.Author),
		Slug:   &t.Slug,
		Short:  &t.Short,

		Category: serialiseCategory(&t.Category),
		Pinned:   t.Pinned,
		Posts:    utils.Ref(len(t.Posts)),
		Reacts:   reacts(t.Reacts),
		Tags:     t.Tags,
	}
}

func serialiseThread(t *thread.Thread) openapi.Thread {
	return openapi.Thread{}
}

func serialisePost(p *post.Post) openapi.Post {
	return openapi.Post{
		Id:        openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(p.DeletedAt),

		//
	}
}

func serialiseProfileReference(a thread.AuthorRef) openapi.ProfileReference {
	return openapi.ProfileReference{
		Id:   openapi.IdentifierFrom(xid.ID(a.ID)),
		Name: &a.Name,
	}
}

func serialiseCategory(c *category.Category) openapi.Category {
	return openapi.Category{
		Id: openapi.IdentifierFrom(xid.ID(c.ID)),
	}
}

func serialiseReact(r *react.React) openapi.React {
	return openapi.React{
		Id:    openapi.IdentifierFrom(xid.ID(r.ID)),
		Emoji: &r.Emoji,
	}
}
