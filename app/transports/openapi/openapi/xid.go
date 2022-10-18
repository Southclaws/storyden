package openapi

import (
	"context"
	"errors"

	"github.com/Southclaws/fault/errtag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
)

// XID converts an openapi identifier to an xid. This is to work around an issue
// with the oapi-codegen generated code which generates the identifier as a new
// type instead of an alias which results in the marshal functions being hidden.
func (i Identifier) XID() xid.ID {
	v, err := xid.FromString(string(i))
	if err != nil {
		return xid.NilID()
	}

	return v
}

// id converts any arbitrary xid.ID derivative to an *openapi.Identifier type.
func IdentifierFrom(id xid.ID) *Identifier {
	oid := Identifier(id.String())
	return &oid
}

// ID will resolve the Unique value to an account's ID using a repository. Since
// this is used a lot as most of the system supports using IDs and handles
// interchangably, the repository this uses should make heavy use of caching.
func (u AccountHandle) ID(ctx context.Context, r account_repo.Repository) (account_repo.AccountID, error) {
	if id, err := xid.FromString(string(u)); err == nil {
		a, err := r.GetByID(ctx, account_repo.AccountID(id))
		if err != nil {
			if model.IsNotFound(err) {
				// NOTE: In the unlikely chance that a user sets their handle to
				// a valid XID string, this will produce a confusing result.
				return account_repo.AccountID(xid.NilID()), errtag.Wrap(errors.New("an account does not exist with the specified ID"), errtag.NotFound{})
			}

			return account_repo.AccountID(xid.NilID()), err
		}

		return account_repo.AccountID(a.ID), nil
	}

	a, err := r.GetByHandle(ctx, string(u))
	if err != nil {
		if model.IsNotFound(err) {
			return account_repo.AccountID(xid.NilID()), errtag.Wrap(errors.New("an account does not exist with the specified handle"), errtag.NotFound{})
		}

		return account_repo.AccountID(xid.NilID()), err
	}

	return account_repo.AccountID(a.ID), nil
}

func (u *AccountHandle) OptionalID(ctx context.Context, r account_repo.Repository) (opt.Optional[account.AccountID], error) {
	if u == nil {
		return opt.NewEmpty[account.AccountID](), nil
	}

	id, err := u.ID(ctx, r)
	if err != nil {
		return nil, err
	}

	return opt.New(id), nil
}
