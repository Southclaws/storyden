package access_key

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/internal/ent"
	ent_auth "github.com/Southclaws/storyden/internal/ent/authentication"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, accountID account.AccountID, kind AccessKeyKind, name string, expiry opt.Optional[time.Time]) (*AccessKeyRecordWithSecret, error) {
	ak := newAccessKey(kind, expiry)

	authRecord, err := r.create(ctx, accountID, ak, name)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	akr, err := AccessKeyRecordFromAuthenticationRecord(*authRecord)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &AccessKeyRecordWithSecret{
		AccessKeyRecord: *akr,
		Name:            authRecord.Name.Or("Unnamed key"),
		secret:          ak.secret,
	}, nil
}

func (r *Repository) create(ctx context.Context, accountID account.AccountID, record AccessKeyRecordWithSecret, name string) (*authentication.Authentication, error) {
	auth, err := r.db.Authentication.Create().
		SetService(authentication.ServiceAccessKey.String()).
		SetTokenType(authentication.TokenTypePasswordHash.String()).
		SetIdentifier(record.GetAuthenticationRecordIdentifier()).
		SetToken(string(record.Hash)).
		SetName(name).
		SetNillableExpiresAt(record.Expires.Ptr()).
		SetAccountAuthentication(xid.ID(accountID)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r.getByID(ctx, auth.ID)
}

func (r *Repository) List(ctx context.Context, accountID account.AccountID) ([]*authentication.Authentication, error) {
	auths, err := r.db.Authentication.Query().
		Where(
			ent_auth.AccountAuthentication(xid.ID(accountID)),
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
		).
		WithAccount().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := dt.MapErr(auths, authentication.FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) ListAllAsAdmin(ctx context.Context) ([]*authentication.Authentication, error) {
	auths, err := r.db.Authentication.Query().
		Where(
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
		).
		WithAccount().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := dt.MapErr(auths, authentication.FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) LookupByToken(ctx context.Context, token *AccessKeyToken) (*authentication.Authentication, error) {
	identifier := token.GetAuthenticationRecordIdentifier()

	auth, err := r.db.Authentication.Query().
		Where(
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
			ent_auth.Identifier(identifier),
		).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, ftag.With(ftag.NotFound), fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := authentication.FromModel(auth)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) Revoke(ctx context.Context, accountID account.AccountID, authID authentication.ID) (*authentication.Authentication, error) {
	auth, err := r.db.Authentication.UpdateOneID(authID).
		SetDisabled(true).
		Where(
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
			ent_auth.AccountAuthentication(xid.ID(accountID)),
		).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, ftag.With(ftag.NotFound), fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := r.getByID(ctx, auth.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) RevokeAsAdmin(ctx context.Context, authID authentication.ID) (*authentication.Authentication, error) {
	auth, err := r.db.Authentication.UpdateOneID(authID).
		SetDisabled(true).
		Where(
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
		).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, ftag.With(ftag.NotFound), fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := r.getByID(ctx, auth.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}

func (r *Repository) getByID(ctx context.Context, id xid.ID) (*authentication.Authentication, error) {
	auth, err := r.db.Authentication.Query().
		Where(
			ent_auth.Service(authentication.ServiceAccessKey.String()),
			ent_auth.TokenType(authentication.TokenTypePasswordHash.String()),
			ent_auth.ID(id),
		).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, ftag.With(ftag.NotFound), fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := authentication.FromModel(auth)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result, nil
}
