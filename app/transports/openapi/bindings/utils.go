package bindings

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/openapi"
)

func serialiseAccount(acc *account.Account) openapi.Account {
	return openapi.Account{
		Id:             openapi.Identifier(acc.ID.String()),
		Handle:         acc.Handle,
		Name:           acc.Name,
		Bio:            acc.Bio.HTML(),
		Links:          serialiseExternalLinks(acc.ExternalLinks),
		Meta:           acc.Metadata,
		CreatedAt:      acc.CreatedAt,
		UpdatedAt:      acc.UpdatedAt,
		DeletedAt:      acc.DeletedAt.Ptr(),
		Admin:          acc.Admin,
		VerifiedStatus: openapi.AccountVerifiedStatus(acc.VerifiedStatus.String()),
		EmailAddresses: dt.Map(acc.EmailAddresses, serialiseEmailAddress),
	}
}

func serialiseEmailAddress(in *account.EmailAddress) openapi.AccountEmailAddress {
	return openapi.AccountEmailAddress{
		EmailAddress: in.Email.Address,
		IsAuth:       in.IsAuth,
		Verified:     in.Verified,
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
	postCount := len(t.Posts)

	return openapi.ThreadReference{
		Id:        openapi.Identifier(xid.ID(t.ID).String()),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt.Ptr(),

		Title:  t.Title,
		Author: serialiseProfileReference(t.Author),
		Slug:   t.Slug,
		Short:  &t.Short,
		Meta:   (*openapi.Metadata)(&t.Meta),

		Category:    serialiseCategoryReference(&t.Category),
		Pinned:      t.Pinned,
		PostCount:   &postCount,
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
		DeletedAt: p.DeletedAt.Ptr(),
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

func serialiseProfileReference(a profile.Public) openapi.ProfileReference {
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
		Filename: a.Name.String(),
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
