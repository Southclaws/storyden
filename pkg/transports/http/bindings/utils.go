package bindings

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/react"
	"github.com/Southclaws/storyden/pkg/resources/thread"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

func serialiseThread(t *thread.Thread) openapi.Thread {
	return openapi.Thread{
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
