package node_auth

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

func AuthoriseNodeMutation(ctx context.Context, acc *account.AccountWithEdges, n *library.Node) error {
	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		ownsNode := n.Owner.ID == acc.ID

		if !ownsNode {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the page and do not have the Manage Library permission."),
			)
		}

		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func AuthoriseNodeParentChildMutation(ctx context.Context, acc *account.AccountWithEdges, cnode, pnode *library.Node) error {
	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		ownsChild := cnode.Owner.ID == acc.ID
		ownsParent := pnode.Owner.ID == acc.ID
		ownsNeither := !ownsChild && !ownsParent

		if ownsNeither {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of both of the pages being affected and do not have the Manage Library permission."),
			)
		}

		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
