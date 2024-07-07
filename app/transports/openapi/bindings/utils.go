package bindings

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

func serialiseAccount(acc *account.Account) openapi.Account {
	return openapi.Account{
		Id:        openapi.Identifier(acc.ID.String()),
		Handle:    acc.Handle,
		Name:      acc.Name,
		Bio:       acc.Bio.HTML(),
		Links:     serialiseExternalLinks(acc.ExternalLinks),
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(acc.DeletedAt),
		Admin:     acc.Admin,
	}
}

func serialiseExternalLinks(in []account.ExternalLink) openapi.ProfileExternalLinkList {
	return dt.Map(in, func(l account.ExternalLink) openapi.ProfileExternalLink {
		return openapi.ProfileExternalLink{
			Text: l.Text,
			Url:  l.URL.String(),
		}
	})
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
		Meta:   (*openapi.Metadata)(&t.Meta),

		Category:    serialiseCategoryReference(&t.Category),
		Pinned:      t.Pinned,
		PostCount:   utils.Ref(len(t.Posts)),
		Reacts:      reacts(t.Reacts),
		Tags:        t.Tags,
		Assets:      dt.Map(t.Assets, serialiseAssetReference),
		Collections: dt.Map(t.Collections, serialiseCollection),
		Link:        opt.Map(t.Links.Latest(), serialiseLink).Ptr(),
	}
}

func serialiseContentHTML(c content.Rich) string {
	return c.HTML()
}

func serialiseThread(t *thread.Thread) openapi.Thread {
	return openapi.Thread{
		Author:    serialiseProfileReference(t.Author),
		Category:  serialiseCategoryReference(&t.Category),
		CreatedAt: t.CreatedAt,
		// DeletedAt: t.DeletedAt,
		Id:             openapi.Identifier(t.ID.String()),
		Meta:           (*openapi.Metadata)(&t.Meta),
		Pinned:         t.Pinned,
		Reacts:         dt.Map(t.Reacts, serialiseReact),
		Short:          &t.Short,
		Slug:           t.Slug,
		Tags:           t.Tags,
		Posts:          dt.Map(t.Posts, serialisePost),
		Title:          t.Title,
		UpdatedAt:      t.UpdatedAt,
		Assets:         dt.Map(t.Assets, serialiseAssetReference),
		Collections:    dt.Map(t.Collections, serialiseCollection),
		Link:           opt.Map(t.Links.Latest(), serialiseLink).Ptr(),
		Recomentations: dt.Map(t.Related, serialiseDatagraphNodeReference),
	}
}

func serialisePost(p *reply.Reply) openapi.PostProps {
	return openapi.PostProps{
		Id:        openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(p.DeletedAt),
		RootId:    p.RootPostID.String(),
		RootSlug:  p.RootThreadMark,
		Body:      p.Content.HTML(),
		Author:    serialiseProfileReference(p.Author),
		Reacts:    dt.Map(p.Reacts, serialiseReact),
		Meta:      (*openapi.Metadata)(&p.Meta),
		Assets:    dt.Map(p.Assets, serialiseAssetReference),
		Links:     dt.Map(p.Links, serialiseLink),
	}
}

func serialiseProfileReference(a datagraph.Profile) openapi.ProfileReference {
	return openapi.ProfileReference{
		Id:     *openapi.IdentifierFrom(xid.ID(a.ID)),
		Handle: (openapi.AccountHandle)(a.Handle),
		Name:   a.Name,
		Admin:  a.Admin,
	}
}

func serialiseCategory(c *category.Category) openapi.Category {
	return openapi.Category{
		Id:          *openapi.IdentifierFrom(xid.ID(c.ID)),
		Name:        c.Name,
		Slug:        c.Slug,
		Colour:      c.Colour,
		Description: c.Description,
		PostCount:   c.PostCount,
		Admin:       c.Admin,
		Sort:        c.Sort,
		Meta:        (*openapi.Metadata)(&c.Metadata),
	}
}

func serialiseCategoryReference(c *category.Category) openapi.CategoryReference {
	return openapi.CategoryReference{
		Id:          *openapi.IdentifierFrom(xid.ID(c.ID)),
		Name:        c.Name,
		Slug:        c.Slug,
		Admin:       c.Admin,
		Colour:      c.Colour,
		Description: c.Description,
		Sort:        c.Sort,
		Meta:        (*openapi.Metadata)(&c.Metadata),
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
		Id:       a.ID.String(),
		Url:      a.URL,
		MimeType: a.Metadata.GetMIMEType(),
		Width:    float32(a.Metadata.GetWidth()),
		Height:   float32(a.Metadata.GetHeight()),
	}
}

func deserialiseAssetID(in string) asset.AssetID {
	return asset.AssetID(openapi.ParseID(in))
}

func deserialiseAssetIDs(ids []string) []asset.AssetID {
	return dt.Map(ids, deserialiseAssetID)
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

func serialiseLink(in *datagraph.Link) openapi.Link {
	return openapi.Link{
		Url:         in.URL,
		Title:       in.Title.Ptr(),
		Description: in.Description.Ptr(),
		Slug:        in.Slug,
		Domain:      in.Domain,
		Assets:      dt.Map(in.Assets, serialiseAssetReference),
	}
}

func deserialiseVisibility(in openapi.Visibility) (visibility.Visibility, error) {
	v, err := visibility.NewVisibility(string(in))
	if err != nil {
		return visibility.Visibility{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return v, nil
}

func serialiseVisibility(in visibility.Visibility) openapi.Visibility {
	return openapi.Visibility(in.String())
}

func deserialiseVisibilityList(in []openapi.Visibility) ([]visibility.Visibility, error) {
	v, err := dt.MapErr(in, deserialiseVisibility)
	if err != nil {
		return nil, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return v, nil
}

func serialiseVisibilityList(in []visibility.Visibility) []openapi.Visibility {
	return dt.Map(in, serialiseVisibility)
}
