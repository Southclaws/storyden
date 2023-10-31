package authentication

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	model_account "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/authentication"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context,
	id account.AccountID,
	service Service,
	identifier string,
	token string,
	metadata map[string]any,
	opts ...Option,
) (*Authentication, error) {
	create := d.db.Authentication.Create()
	mutate := create.Mutation()

	mutate.SetAccountID(xid.ID(id))
	mutate.SetService(string(service))
	mutate.SetIdentifier(identifier)
	mutate.SetToken(token)
	mutate.SetMetadata(metadata)

	for _, fn := range opts {
		fn(mutate)
	}

	r, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	r, err = d.db.Authentication.
		Query().
		Where(authentication.ID(r.ID)).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(r)
}

func (d *database) LookupByIdentifier(ctx context.Context, service Service, identifier string) (*Authentication, bool, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.IdentifierEQ(identifier),
			authentication.ServiceEQ(string(service)),
		).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	auth, err := FromModel(r)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return auth, true, nil
}

func (d *database) LookupByHandle(ctx context.Context, service Service, handle string) (*Authentication, bool, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.ServiceEQ(string(service)),
			authentication.HasAccountWith(model_account.Handle(handle)),
		).
		WithAccount().
		Order(authentication.ByCreatedAt(sql.OrderDesc())).
		// NOTE: We pick the latest here and don't handle multiple providers of
		// the same type very well (it's quite an edge case...)
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	auth, err := FromModel(r)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return auth, true, nil
}

func (d *database) GetAuthMethods(ctx context.Context, id account.AccountID) ([]*Authentication, error) {
	r, err := d.db.Authentication.
		Query().
		Where(authentication.HasAccountWith(model_account.IDEQ(xid.ID(id)))).
		WithAccount().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	auths, err := dt.MapErr(r, FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return auths, nil
}

func (d *database) Update(ctx context.Context, id ID, opts ...Option) (*Authentication, error) {
	update := d.db.Authentication.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	r, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	r, err = d.db.Authentication.Query().Where(authentication.ID(r.ID)).WithAccount().Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, err := FromModel(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}

func (d *database) DeleteByID(ctx context.Context, accountID account.AccountID, aid ID) (bool, error) {
	n, err := d.db.Authentication.
		Delete().
		Where(
			authentication.HasAccountWith(
				model_account.ID(xid.ID(accountID)),
			),
			authentication.ID(aid),
		).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}

		return false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return n > 0, nil
}

func (d *database) Delete(ctx context.Context, accountID account.AccountID, identifier string, service Service) (bool, error) {
	n, err := d.db.Authentication.
		Delete().
		Where(
			authentication.HasAccountWith(
				model_account.ID(xid.ID(accountID)),
			),
			authentication.IdentifierEQ(identifier),
			authentication.ServiceEQ(string(service)),
		).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}

		return false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return n > 0, nil
}
