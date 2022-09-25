package authentication

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/fault/errtag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	model_account "github.com/Southclaws/storyden/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/authentication"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context,
	id account.AccountID,
	service Service,
	identifier string,
	token string,
	metadata map[string]any,
) (*Authentication, error) {
	r, err := d.db.Authentication.Create().
		SetAccountID(xid.ID(id)).
		SetService(string(service)).
		SetIdentifier(identifier).
		SetToken(token).
		SetMetadata(metadata).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.AlreadyExists{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	r, err = d.db.Authentication.
		Query().
		Where(authentication.ID(r.ID)).
		WithAccount().
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(r), nil
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
		if model.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(r), true, nil
}

func (d *database) GetAuthMethods(ctx context.Context, id account.AccountID) ([]Authentication, error) {
	r, err := d.db.Authentication.
		Query().
		Where(authentication.HasAccountWith(model_account.IDEQ(xid.ID(id)))).
		All(ctx)
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModelMany(r), nil
}

func (d *database) IsEqual(ctx context.Context, id account.AccountID, identifier string, token string) (bool, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.HasAccountWith(model_account.IDEQ(xid.ID(id))),
			authentication.IdentifierEQ(identifier),
		).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return false, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return false, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return r.Token == token, nil
}
