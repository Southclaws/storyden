package warning_repo

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/warning"
	"github.com/Southclaws/storyden/internal/ent"
	ent_warning "github.com/Southclaws/storyden/internal/ent/warning"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context,
	givenBy account.AccountID,
	receiveBy account.AccountID,
	reason string,
) (*warning.Warning, error) {
	create := r.db.Warning.Create().
		SetAccountID(xid.ID(receiveBy)).
		SetAuthorID(xid.ID(givenBy)).
		SetReason(reason)

	created, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r.Get(ctx, receiveBy, created.ID)
}

func (r *Repository) ListByAccountID(ctx context.Context, accountID account.AccountID) (warning.Warnings, error) {
	items, err := r.db.Warning.Query().
		Where(ent_warning.AccountIDEQ(xid.ID(accountID))).
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) {
				arq.WithRole()
			})
		}).
		Order(ent_warning.ByCreatedAt(sql.OrderDesc()), ent_warning.ByID(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	out, err := dt.MapErr(items, warning.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return out, nil
}

func (r *Repository) UpdateReason(ctx context.Context, accountID account.AccountID, warningID warning.ID, reason string) (*warning.Warning, error) {
	record, err := r.Get(ctx, accountID, warningID)
	if err != nil {
		return nil, err
	}

	if _, err := r.db.Warning.UpdateOneID(record.ID).SetReason(reason).Save(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err = r.Get(ctx, accountID, warningID)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *Repository) Delete(ctx context.Context, accountID account.AccountID, warningID warning.ID) error {
	result, err := r.db.Warning.Query().
		Where(
			ent_warning.IDEQ(warningID),
			ent_warning.AccountIDEQ(xid.ID(accountID)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.db.Warning.DeleteOneID(result.ID).Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, accountID account.AccountID, warningID warning.ID) (*warning.Warning, error) {
	result, err := r.db.Warning.Query().
		Where(
			ent_warning.IDEQ(warningID),
			ent_warning.AccountIDEQ(xid.ID(accountID)),
		).
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) {
				arq.WithRole()
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	warning, err := warning.Map(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return warning, nil
}
