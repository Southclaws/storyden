package account_repo

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/schema"
	"github.com/rs/xid"
)

type (
	Option func(*ent.AccountMutation)
)

func WithID(id account.AccountID) Option {
	return func(a *ent.AccountMutation) {
		a.SetID(xid.ID(id))
	}
}

func WithAdmin(admin bool) Option {
	return func(a *ent.AccountMutation) {
		a.SetAdmin(admin)
	}
}

func WithName(name string) Option {
	return func(a *ent.AccountMutation) {
		a.SetName(name)
	}
}

func WithBio(v datagraph.Content) Option {
	return func(a *ent.AccountMutation) {
		a.SetBio(v.HTML())
	}
}

func WithSignature(v datagraph.Content) Option {
	return func(a *ent.AccountMutation) {
		a.SetSignature(v.HTML())
	}
}

func WithInvitedBy(id xid.ID) Option {
	return func(a *ent.AccountMutation) {
		a.SetInvitedByID(id)
	}
}

type Mutation func(u *ent.AccountUpdateOne)

func SetHandle(handle string) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetHandle(handle)
	}
}

func SetName(name string) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetName(name)
	}
}

func SetBio(bio string) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetBio(bio)
	}
}

func SetSignature(signature string) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetSignature(signature)
	}
}

func SetAdmin(status bool) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetAdmin(status)
	}
}

func SetVerifiedStatus(status account.VerifiedStatus) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetVerifiedStatus(account_ent.VerifiedStatus(status.String()))
	}
}

func SetInterests(interests []xid.ID) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.ClearTags().AddTagIDs(interests...)
	}
}

func SetLinks(links []account.ExternalLink) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetLinks(dt.Map(links, func(i account.ExternalLink) schema.ExternalLink {
			return schema.ExternalLink{
				Text: i.Text,
				URL:  i.URL.String(),
			}
		}))
	}
}

func SetMetadata(m map[string]any) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetMetadata(m)
	}
}

func SetDeleted(t opt.Optional[time.Time]) Mutation {
	return func(u *ent.AccountUpdateOne) {
		if v, ok := t.Get(); ok {
			u.SetDeletedAt(v)
		} else {
			u.ClearDeletedAt()
		}
	}
}

func (r *Repository) Create(ctx context.Context, handle string, opts ...Option) (*account.AccountWithEdges, error) {
	create := r.db.Account.Create()
	mutate := create.Mutation()

	mutate.SetHandle(handle)
	mutate.SetName(handle)

	for _, v := range opts {
		v(mutate)
	}

	saved, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return r.GetByID(ctx, account.AccountID(saved.ID))
}

func (r *Repository) Update(ctx context.Context, id account.AccountID, opts ...Mutation) (*account.AccountWithEdges, error) {
	update := r.db.Account.UpdateOneID(xid.ID(id))

	for _, fn := range opts {
		fn(update)
	}

	saved, err := update.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("unique constraint violation", "The specified handle has already been used."))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r.GetByID(ctx, account.AccountID(saved.ID))
}
