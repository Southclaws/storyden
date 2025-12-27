package bindings

import (
	"context"
	"net/url"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func getAuthStateCacheControl(ctx context.Context, rest string) string {
	authenticated := session.GetOptAccountID(ctx).Ok()
	if authenticated {
		return "private, " + rest
	} else {
		return "public, " + rest
	}
}

func serialiseAccount(acc *account.AccountWithEdges) openapi.Account {
	invitedBy := opt.Map(acc.InvitedBy, func(ib account.Account) openapi.ProfileReference {
		return serialiseProfileReferenceFromAccount(ib)
	})

	return openapi.Account{
		Id:             openapi.Identifier(acc.ID.String()),
		Joined:         acc.CreatedAt,
		Suspended:      acc.DeletedAt.Ptr(),
		Handle:         acc.Handle,
		Name:           acc.Name,
		Bio:            acc.Bio.HTML(),
		Meta:           acc.Metadata,
		Links:          serialiseExternalLinks(acc.ExternalLinks),
		CreatedAt:      acc.CreatedAt,
		UpdatedAt:      acc.UpdatedAt,
		DeletedAt:      acc.DeletedAt.Ptr(),
		Admin:          acc.Admin,
		VerifiedStatus: openapi.AccountVerifiedStatus(acc.VerifiedStatus.String()),
		EmailAddresses: dt.Map(acc.EmailAddresses, serialiseEmailAddressPtr),
		Roles:          serialiseHeldRoleList(acc.Roles),
		InvitedBy:      invitedBy.Ptr(),
	}
}

func serialiseEmailAddress(in account.EmailAddress) openapi.AccountEmailAddress {
	return openapi.AccountEmailAddress{
		Id:           in.ID.String(),
		EmailAddress: in.Email.Address,
		Verified:     in.Verified,
	}
}

func serialiseEmailAddressPtr(in *account.EmailAddress) openapi.AccountEmailAddress {
	return serialiseEmailAddress(*in)
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
		DeletedAt: t.DeletedAt.Ptr(),

		Title:       t.Title,
		Author:      serialiseProfileReference(t.Author),
		Slug:        t.Slug,
		Description: &t.Short,
		Body:        t.Content.HTML(),
		Meta:        (*openapi.Metadata)(&t.Meta),
		LastReplyAt: t.LastReplyAt.Ptr(),

		Category:    opt.Map(t.Category, serialiseCategoryReference).Ptr(),
		Visibility:  serialiseVisibility(t.Visibility),
		Pinned:      t.Pinned,
		ReadStatus:  opt.PtrMap(t.ReadStatus, serialiseReadStatus),
		ReplyStatus: serialiseReplyStatus(t.ReplyStatus),
		Likes:       serialiseLikeStatus(&t.Likes),
		Reacts:      serialiseReactList(t.Reacts),
		Tags:        serialiseTagReferenceList(t.Tags),
		Assets:      dt.Map(t.Assets, serialiseAssetPtr),
		Collections: serialiseCollectionStatus(t.Collections),
		Link:        opt.Map(t.WebLink, serialiseLinkRef).Ptr(),
	}
}

func serialiseContentHTML(c datagraph.Content) string {
	return c.HTML()
}

func serialiseThread(t *thread.Thread) openapi.Thread {
	return openapi.Thread{
		Assets:         dt.Map(t.Assets, serialiseAssetPtr),
		Author:         serialiseProfileReference(t.Author),
		Body:           serialiseContentHTML(t.Content),
		Category:       opt.Map(t.Category, serialiseCategoryReference).Ptr(),
		Likes:          serialiseLikeStatus(&t.Likes),
		Collections:    serialiseCollectionStatus(t.Collections),
		CreatedAt:      t.CreatedAt,
		DeletedAt:      t.DeletedAt.Ptr(),
		Description:    &t.Short,
		Visibility:     serialiseVisibility(t.Visibility),
		Id:             openapi.Identifier(t.ID.String()),
		Link:           opt.Map(t.WebLink, serialiseLinkRef).Ptr(),
		Meta:           (*openapi.Metadata)(&t.Meta),
		Pinned:         t.Pinned,
		ReadStatus:     opt.PtrMap(t.ReadStatus, serialiseReadStatus),
		ReplyStatus:    serialiseReplyStatus(t.ReplyStatus),
		Reacts:         dt.Map(t.Reacts, serialiseReact),
		Recomentations: dt.Map(t.Related, serialiseDatagraphItem),
		Replies:        serialiseThreadRepliesPaginatedList(t.Replies),
		Slug:           t.Slug,
		Tags:           serialiseTagReferenceList(t.Tags),
		Title:          t.Title,
		UpdatedAt:      t.UpdatedAt,
		LastReplyAt:    t.LastReplyAt.Ptr(),
	}
}

func serialiseThreadRepliesPaginatedList(in pagination.Result[*reply.Reply]) openapi.PaginatedReplyList {
	return openapi.PaginatedReplyList{
		CurrentPage: in.CurrentPage,
		NextPage:    in.NextPage.Ptr(),
		PageSize:    in.Size,
		Results:     in.Results,
		Replies:     dt.Map(in.Items, serialiseReplyPtr),
		TotalPages:  in.TotalPages,
	}
}

func serialiseReplyStatus(s post.ReplyStatus) openapi.ReplyStatus {
	return openapi.ReplyStatus{
		Replies: s.Count,
		Replied: s.Replied,
	}
}

func serialiseReadStatus(s post.ReadStatus) openapi.ReadStatus {
	return openapi.ReadStatus{
		RepliesSince: s.Count,
		LastReadAt:   s.LastReadAt,
	}
}

func serialiseReplyPtr(p *reply.Reply) openapi.Reply {
	description := p.Content.Short()
	return openapi.Reply{
		Id:          openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt.Ptr(),
		RootId:      p.RootPostID.String(),
		RootSlug:    p.RootThreadMark,
		Slug:        p.Slug,
		Title:       p.RootThreadTitle,
		Body:        p.Content.HTML(),
		Visibility:  serialiseVisibility(p.Visibility),
		Description: &description,
		Author:      serialiseProfileReference(p.Author),
		Likes:       serialiseLikeStatus(&p.Likes),
		Reacts:      dt.Map(p.Reacts, serialiseReact),
		Meta:        (*openapi.Metadata)(&p.Meta),
		Assets:      dt.Map(p.Assets, serialiseAssetPtr),
		ReplyTo:     opt.Map(p.ReplyTo, serialiseReply).Ptr(),
	}
}

func serialiseReply(p reply.Reply) openapi.Reply {
	return serialiseReplyPtr(&p)
}

func serialisePost(p *post.Post) openapi.Post {
	description := p.Content.Short()
	return openapi.Post{
		Id:          openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt:   p.CreatedAt,
		DeletedAt:   p.DeletedAt.Ptr(),
		UpdatedAt:   p.UpdatedAt,
		Title:       p.Title,
		Description: &description,
		Slug:        p.Slug,
		Body:        p.Content.HTML(),
		Author:      serialiseProfileReference(p.Author),
		Assets:      dt.Map(p.Assets, serialiseAssetPtr),
		Collections: openapi.CollectionStatus{},
		Likes:       serialiseLikeStatus(&p.Likes),
		Reacts:      dt.Map(p.Reacts, serialiseReact),
		Meta:        (*openapi.Metadata)(&p.Meta),
		Visibility:  serialiseVisibility(p.Visibility),
	}
}

func serialisePostRef(p *post.Post) openapi.PostReference {
	return openapi.PostReference{
		Id:         openapi.Identifier(xid.ID(p.ID).String()),
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  p.DeletedAt.Ptr(),
		Author:     serialiseProfileReference(p.Author),
		Likes:      serialiseLikeStatus(&p.Likes),
		Reacts:     dt.Map(p.Reacts, serialiseReact),
		Meta:       (*openapi.Metadata)(&p.Meta),
		Assets:     dt.Map(p.Assets, serialiseAssetPtr),
		Visibility: serialiseVisibility(p.Visibility),
	}
}

func deserialisePostID(s string) post.ID {
	return post.ID(openapi.ParseID(s))
}

func deserialiseID(t openapi.Identifier) xid.ID {
	return openapi.ParseID(t)
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

func deserialiseMark(s string) mark.Queryable {
	return mark.NewQueryKey(s)
}

func deserialiseOptionalFloat(in *float32) opt.Optional[float64] {
	return opt.NewPtrMap(in, func(s float32) float64 { return float64(s) })
}

func serialiseOptionalFloat(in opt.Optional[float64]) *float32 {
	return opt.PtrMap(in, func(s float64) float32 { return float32(s) })
}

func seraliseOptionalURL(in opt.Optional[url.URL]) *string {
	return opt.PtrMap(in, func(s url.URL) string { return s.String() })
}

// NOTE: Page query parameters are string types due to a bug in oapi-codegen.
func deserialisePageParams(p *string, pageSize uint) pagination.Parameters {
	pageNumber := opt.NewPtrMap(p, func(s string) uint {
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return 1
		}

		return max(1, uint(v))
	}).Or(1)

	return pagination.NewPageParams(pageNumber, pageSize)
}
