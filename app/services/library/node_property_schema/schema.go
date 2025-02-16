package node_property_schema

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
)

type Updater struct {
	accountQuery *account_querier.Querier
	nodeQuerier  *node_querier.Querier
	nsr          *node_properties.SchemaWriter
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nsr *node_properties.SchemaWriter,
) *Updater {
	return &Updater{
		accountQuery: accountQuery,
		nodeQuerier:  nodeQuerier,
		nsr:          nsr,
	}
}

func (u *Updater) UpdateChildren(ctx context.Context, qk library.QueryKey, schemas node_properties.FieldSchemaMutations) (*library.PropertySchema, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := u.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := u.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeMutation(ctx, acc, n); err != nil {
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

	acc, err := u.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := u.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeMutation(ctx, acc, n); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	schema, err := u.nsr.UpdateSiblings(ctx, qk, schemas)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return schema, nil
}
