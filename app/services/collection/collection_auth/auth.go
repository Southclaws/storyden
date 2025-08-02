package collection_auth

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func CheckCollectionMutationPermissions(ctx context.Context, col collection.Collection) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := session.Authorise(ctx, func() error {
		if accountID != col.Owner.ID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the collection and do not have the Manage Collections permission."),
			)
		}
		return nil
	}, rbac.PermissionManageCollections); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

var ErrNotResourcesOwner = fault.New("not the owner of either collection or item", ftag.With(ftag.PermissionDenied))

func CheckCollectionItemMutationPermissions(ctx context.Context, acc account.Account, cis collection.CollectionItemStatus) (collection.MembershipType, error) {
	if item, ok := cis.Item.Get(); ok && item.Author.ID != acc.ID {
		return collection.MembershipTypeSubmissionAccepted, nil
	}

	if cis.Collection.Owner.ID != acc.ID {
		return collection.MembershipTypeSubmissionReview, nil
	}

	if acc.ID == cis.Collection.Owner.ID {
		return collection.MembershipTypeNormal, nil
	}

	return collection.MembershipType{}, fault.Wrap(ErrNotResourcesOwner, fctx.With(ctx))
}
