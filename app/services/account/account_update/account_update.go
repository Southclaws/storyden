package account_update

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Updater struct {
	fx.In

	Log        *zap.Logger
	Writer     *account_writer.Writer
	IndexQueue pubsub.Topic[mq.IndexProfile]
}

type Partial struct {
	Handle    opt.Optional[string]
	Name      opt.Optional[string]
	Bio       opt.Optional[string]
	Interests opt.Optional[[]xid.ID]
	Links     opt.Optional[[]account.ExternalLink]
	Meta      opt.Optional[map[string]any]
}

func (u *Updater) Update(ctx context.Context, id account.AccountID, params Partial) (*account.Account, error) {
	opts := []account_writer.Mutation{}

	if v, ok := params.Handle.Get(); ok {
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

	acc, err := u.Writer.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := u.IndexQueue.Publish(ctx, mq.IndexProfile{
		ID: id,
	}); err != nil {
		u.Log.Error("failed to publish index post message", zap.Error(err))
	}

	return acc, nil
}
