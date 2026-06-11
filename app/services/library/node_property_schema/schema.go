package node_property_schema

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/ent"
)

type Updater struct {
	accountQuery   *account_querier.Querier
	nodeQuerier    *node_querier.Querier
	versionQuerier *node_version_querier.Querier
	nodeWriter     *node_writer.Writer
	nsr            *node_properties.SchemaWriter
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	versionQuerier *node_version_querier.Querier,
	nodeWriter *node_writer.Writer,
	nsr *node_properties.SchemaWriter,
) *Updater {
	return &Updater{
		accountQuery:   accountQuery,
		nodeQuerier:    nodeQuerier,
		versionQuerier: versionQuerier,
		nodeWriter:     nodeWriter,
		nsr:            nsr,
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

	if err := u.guardVersionedMutation(ctx, qk); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schema, err := u.nsr.UpdateSiblings(ctx, qk, schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = u.nodeWriter.Update(ctx, qk, node_writer.WithCurrentVersionCleared())
	if err != nil && !ent.IsNotFound(err) {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return schema, nil
}

func (u *Updater) guardVersionedMutation(ctx context.Context, qk library.QueryKey) error {
	hasDraft, err := u.versionQuerier.HasNodeDraft(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if hasDraft {
		return fault.New("node has a draft version",
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc("draft exists", "This node has a working draft. Apply or delete the draft before editing versioned page fields directly."),
		)
	}

	return nil
}
