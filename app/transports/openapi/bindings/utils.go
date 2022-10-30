package bindings

import (
	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

func serialiseAccount(acc *account.Account) openapi.Account {
	return openapi.Account{
		Id:        openapi.Identifier(acc.ID.String()),
		Handle:    (*openapi.AccountHandle)(&acc.Handle),
		Name:      utils.Ref(acc.Name),
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(acc.DeletedAt),
	}
}

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
		Meta:   t.Meta,

		Category: serialiseCategory(&t.Category),
		Pinned:   t.Pinned,
		Posts:    utils.Ref(len(t.Posts)),
		Reacts:   reacts(t.Reacts),
		Tags:     t.Tags,
	}
}

func serialiseThread(t *thread.Thread) openapi.Thread {
	return openapi.Thread{
		Author:    serialiseProfileReference(t.Author),
		Category:  serialiseCategory(&t.Category),
		CreatedAt: t.CreatedAt,
		// DeletedAt: t.DeletedAt,
		Id:        openapi.Identifier(t.ID.String()),
		Meta:      t.Meta,
		Pinned:    t.Pinned,
		Reacts:    dt.Map(t.Reacts, serialiseReact),
		Short:     &t.Short,
		Slug:      &t.Slug,
		Tags:      t.Tags,
		Title:     t.Title,
		UpdatedAt: t.UpdatedAt,
	}
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

func serialiseTag(t tag.Tag) openapi.Tag {
	return openapi.Tag{
		Id:   openapi.Identifier(t.ID),
		Name: t.Name,
	}
}

func tagID(t openapi.Tag) xid.ID {
	return t.Id.XID()
}

// tagsIDs just applies tagID to a slice so we get a slice of IDs back.
func tagsIDs(i []openapi.Tag) []xid.ID {
	return dt.Map(i, tagID)
}
