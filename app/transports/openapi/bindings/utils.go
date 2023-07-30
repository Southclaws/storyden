package bindings

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

func serialiseAccount(acc *account.Account) openapi.Account {
	return openapi.Account{
		Id:        openapi.Identifier(acc.ID.String()),
		Handle:    acc.Handle,
		Name:      utils.Ref(acc.Name),
		Bio:       utils.Ref(acc.Bio.OrZero()),
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
		Slug:   t.Slug,
		Short:  &t.Short,
		Meta:   t.Meta,

		Category:  serialiseCategoryReference(&t.Category),
		Pinned:    t.Pinned,
		PostCount: utils.Ref(len(t.Posts)),
		Reacts:    reacts(t.Reacts),
		Tags:      t.Tags,
		Assets:    dt.Map(t.Assets, serialiseAssetReference),
	}
}

func serialiseThread(t *thread.Thread) (*openapi.Thread, error) {
	posts, err := dt.MapErr(t.Posts, serialisePost)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &openapi.Thread{
		Author:    serialiseProfileReference(t.Author),
		Category:  serialiseCategoryReference(&t.Category),
		CreatedAt: t.CreatedAt,
		// DeletedAt: t.DeletedAt,
		Id:        openapi.Identifier(t.ID.String()),
		Meta:      t.Meta,
		Pinned:    t.Pinned,
		Reacts:    dt.Map(t.Reacts, serialiseReact),
		Short:     &t.Short,
		Slug:      t.Slug,
		Tags:      t.Tags,
		Posts:     posts,
		Title:     t.Title,
		UpdatedAt: t.UpdatedAt,
		Assets:    dt.Map(t.Assets, serialiseAssetReference),
	}, nil
}

func serialisePost(p *reply.Reply) (openapi.PostProps, error) {
	return openapi.PostProps{
		Id:        openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(p.DeletedAt),
		RootId:    p.RootPostID.String(),
		RootSlug:  p.RootThreadMark,
		Body:      p.Body,
		Author:    serialiseProfileReference(p.Author),
		Reacts:    dt.Map(p.Reacts, serialiseReact),
		Meta:      (*openapi.Metadata)(&p.Meta),
		Assets:    dt.Map(p.Assets, serialiseAssetReference),
	}, nil
}

func serialiseProfileReference(a account.Account) openapi.ProfileReference {
	return openapi.ProfileReference{
		Id:     *openapi.IdentifierFrom(xid.ID(a.ID)),
		Handle: (openapi.AccountHandle)(a.Handle),
		Name:   a.Name,
	}
}

func serialiseCategory(c *category.Category) openapi.Category {
	return openapi.Category{
		Id:          openapi.IdentifierFrom(xid.ID(c.ID)),
		Admin:       &c.Admin,
		Colour:      &c.Colour,
		Description: &c.Description,
		Name:        c.Name,
		PostCount:   c.PostCount,
		Sort:        c.Sort,
	}
}

func serialiseCategoryReference(c *category.Category) openapi.CategoryReference {
	return openapi.CategoryReference{
		Id:          openapi.IdentifierFrom(xid.ID(c.ID)),
		Admin:       &c.Admin,
		Colour:      &c.Colour,
		Description: &c.Description,
		Name:        c.Name,
		Sort:        c.Sort,
	}
}

func serialiseReact(r *react.React) openapi.React {
	return openapi.React{
		Id:    openapi.IdentifierFrom(xid.ID(r.ID)),
		Emoji: &r.Emoji,
	}
}

func serialiseAssetReference(a *asset.Asset) openapi.Asset {
	return openapi.Asset{
		Id:       string(a.ID),
		Url:      a.URL,
		MimeType: a.MIMEType,
		Width:    float32(a.Width),
		Height:   float32(a.Height),
	}
}

func serialiseTag(t tag.Tag) openapi.Tag {
	return openapi.Tag{
		Id:   openapi.Identifier(t.ID),
		Name: t.Name,
	}
}

func deserialiseID(t openapi.Identifier) xid.ID {
	return openapi.ParseID(t)
}

func tagsIDs(i openapi.TagListIDs) []xid.ID {
	return dt.Map(i, deserialiseID)
}
