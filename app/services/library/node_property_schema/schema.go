package node_property_schema

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Updater struct {
	nodeQuerier *node_querier.Querier
	nsr         *node_properties.SchemaWriter
}

func New(
	nodeQuerier *node_querier.Querier,
	nsr *node_properties.SchemaWriter,
) *Updater {
	return &Updater{
		nodeQuerier: nodeQuerier,
		nsr:         nsr,
	}
}

func (u *Updater) UpdateChildren(ctx context.Context, qk library.QueryKey, schemas node_properties.FieldSchemaMutations) (*library.PropertySchema, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := u.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := session.Authorise(ctx, func() error {
		if n.Owner.ID != accountID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the page and do not have the Manage Library permission."))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schema, err := u.nsr.UpdateChildren(ctx, qk, schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return schema, nil
}

func (u *Updater) UpdateSiblings(ctx context.Context, qk library.QueryKey, schemas node_properties.FieldSchemaMutations) (*library.PropertySchema, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := u.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := session.Authorise(ctx, func() error {
		if n.Owner.ID != accountID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the page and do not have the Manage Library permission."))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schema, err := u.nsr.UpdateSiblings(ctx, qk, schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return schema, nil
}
