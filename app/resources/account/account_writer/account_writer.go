package account_writer

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

type Writer struct {
	db             *ent.Client
	accountQuerier *account_querier.Querier
}

func New(db *ent.Client, accountQuerier *account_querier.Querier) *Writer {
	return &Writer{db: db, accountQuerier: accountQuerier}
}

type (
	Option   func(*account.Account)
	Mutation func(u *ent.AccountUpdateOne)
)

func WithID(id account.AccountID) Option {
	return func(a *account.Account) {
		a.ID = account.AccountID(id)
	}
}

func WithAdmin(admin bool) Option {
	return func(a *account.Account) {
		a.Admin = admin
	}
}

func WithName(name string) Option {
	return func(a *account.Account) {
		a.Name = name
	}
}

func WithBio(v datagraph.Content) Option {
	return func(a *account.Account) {
		a.Bio = v
	}
}

func WithInvitedBy(id xid.ID) Option {
	return func(a *account.Account) {
		a.InvitedByID = &id
	}
}

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

func SetAdmin(status bool) Mutation {
	return func(u *ent.AccountUpdateOne) {
		u.SetAdmin(status)
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

func (d *Writer) Create(ctx context.Context, handle string, opts ...Option) (*account.Account, error) {
	withrequired := account.Account{
		Handle: handle,
		Name:   handle, // default display name is just the handle
	}

	for _, v := range opts {
		v(&withrequired)
	}

	create := d.db.Account.Create()

	if !xid.ID(withrequired.ID).IsNil() {
		create.SetID(xid.ID(withrequired.ID))
	}

	a, err := create.
		SetHandle(withrequired.Handle).
		SetName(withrequired.Name).
		SetBio(withrequired.Bio.HTML()).
		SetAdmin(withrequired.Admin).
		SetNillableInvitedByID(withrequired.InvitedByID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return d.accountQuerier.GetByID(ctx, account.AccountID(a.ID))
}

func (d *Writer) Update(ctx context.Context, id account.AccountID, opts ...Mutation) (*account.Account, error) {
	update := d.db.Account.UpdateOneID(xid.ID(id))

	for _, fn := range opts {
		fn(update)
	}

	acc, err := update.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			// Currently the only valid constraint worth checking is the handle.
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("unique constraint violation", "The specified handle has already been used."))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.accountQuerier.GetByID(ctx, account.AccountID(acc.ID))
}
