package openapi

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	model "github.com/Southclaws/storyden/internal/ent"
)

var (
	ErrNoAccountWithID     = errors.New("an account does not exist with the specified ID")
	ErrNoAccountWithHandle = errors.New("an account does not exist with the specified handle")
)

// XID converts an openapi identifier to an xid. This is to work around an fmsg
// with the oapi-codegen generated code which generates the identifier as a new
// type instead of an alias which results in the marshal functions being hidden.
func GetAccountID(i Identifier) xid.ID {
	v, err := xid.FromString(string(i))
	if err != nil {
		return xid.NilID()
	}

	return v
}

// XID converts a thread mark (id-slug) to just the XID, same as above.
func ParseID(i Identifier) xid.ID {
	if len(i) < 20 {
		return xid.NilID()
	}

	v, err := xid.FromString(string(i[:20]))
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
func ResolveHandle(ctx context.Context, r *profile_querier.Querier, u AccountHandle) (account.AccountID, error) {
	if id, err := xid.FromString(string(u)); err == nil {
		a, err := r.GetByID(ctx, account.AccountID(id))
		if err != nil {
			if model.IsNotFound(err) {
				// NOTE: In the unlikely chance that a user sets their handle to
				// a valid XID string, this will produce a confusing result.
				return account.AccountID(xid.NilID()), fault.Wrap(ErrNoAccountWithID, fctx.With(ctx), ftag.With(ftag.NotFound))
			}

			return account.AccountID(xid.NilID()), err
		}

		return account.AccountID(a.ID), nil
	}

	a, found, err := r.LookupByHandle(ctx, string(u))
	if err != nil {
		return account.AccountID(xid.NilID()), err
	}

	if !found {
		return account.AccountID(xid.NilID()), fault.Wrap(ErrNoAccountWithHandle, fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	return account.AccountID(a.ID), nil
}

func OptionalID(ctx context.Context, r *profile_querier.Querier, u *AccountHandle) (opt.Optional[account.AccountID], error) {
	if u == nil {
		return opt.NewEmpty[account.AccountID](), nil
	}

	id, err := ResolveHandle(ctx, r, *u)
	if err != nil {
		return nil, err
	}

	return opt.New(id), nil
}
