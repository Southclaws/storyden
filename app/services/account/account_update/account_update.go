package account_update

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

// TODO: Should be named profile updater tbh, is not account-specific.
type Updater struct {
	writer       *account_writer.Writer
	profileCache *profile_cache.Cache
	bus          *pubsub.Bus
}

func New(
	writer *account_writer.Writer,
	profileCache *profile_cache.Cache,
	bus *pubsub.Bus,
) *Updater {
	return &Updater{
		writer:       writer,
		profileCache: profileCache,
		bus:          bus,
	}
}

type Partial struct {
	Handle    opt.Optional[string]
	Name      opt.Optional[string]
	Bio       opt.Optional[string]
	Interests opt.Optional[[]xid.ID]
	Links     opt.Optional[[]account.ExternalLink]
	Meta      opt.Optional[map[string]any]
}

func (u *Updater) Update(ctx context.Context, id account.AccountID, params Partial) (*account.AccountWithEdges, error) {
	opts := []account_writer.Mutation{}

	if v, ok := params.Handle.Get(); ok {
		if err := account.ValidateHandle(ctx, v); err != nil {
			return nil, err
		}

		opts = append(opts, account_writer.SetHandle(v))
	}
	if v, ok := params.Name.Get(); ok {
		opts = append(opts, account_writer.SetName(v))
	}
	if v, ok := params.Bio.Get(); ok {
		opts = append(opts, account_writer.SetBio(v))
	}
	if v, ok := params.Interests.Get(); ok {
		opts = append(opts, account_writer.SetInterests(v))
	}
	if v, ok := params.Links.Get(); ok {
		opts = append(opts, account_writer.SetLinks(v))
	}
	if v, ok := params.Meta.Get(); ok {
		opts = append(opts, account_writer.SetMetadata(v))
	}

	err := u.profileCache.Invalidate(ctx, xid.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := u.writer.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: id,
	})

	return acc, nil
}
