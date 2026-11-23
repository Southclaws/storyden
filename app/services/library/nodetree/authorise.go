package nodetree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func AuthoriseOwner(ctx context.Context, node *library.Node) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return session.Authorise(ctx, func() error {
		if node.Owner.ID != accountID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the page and do not have the Manage Library permission."))
		}

		return nil
	}, rbac.PermissionManageLibrary)
}

func AuthoriseParent(ctx context.Context, parent *library.Node, childVisibility visibility.Visibility) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return session.Authorise(ctx, func() error {
		if parent.Owner.ID == accountID {
			return nil
		}

		isPublished := childVisibility == visibility.VisibilityPublished
		if !isPublished && parent.Visibility == visibility.VisibilityPublished {
			return nil
		}

		return fault.Wrap(rbac.ErrPermissions,
			fctx.With(ctx),
			fmsg.WithDesc("not parent owner", "You do not own the parent page and do not have the Manage Library permission. Only unpublished pages may be added to another member's published page."),
		)
	}, rbac.PermissionManageLibrary)
}
