package account_update

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/event"
)

// TODO: Should be named profile updater tbh, is not account-specific.
type Updater struct {
	writer *account_writer.Writer
	bus    *event.Bus
}

func New(
	writer *account_writer.Writer,
	bus *event.Bus,
) *Updater {
	return &Updater{
		writer: writer,
		bus:    bus,
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

	acc, err := u.writer.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u.bus.Publish(ctx, &mq.EventAccountUpdated{
		ID: id,
	})

	return acc, nil
}
