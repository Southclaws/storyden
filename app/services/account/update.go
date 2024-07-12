package account

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mq"
)

type Partial struct {
	Handle    opt.Optional[string]
	Name      opt.Optional[string]
	Bio       opt.Optional[string]
	Interests opt.Optional[[]xid.ID]
	Links     opt.Optional[[]account.ExternalLink]
	Meta      opt.Optional[map[string]any]
}

func (s *service) Update(ctx context.Context, id account.AccountID, params Partial) (*account.Account, error) {
	opts := []account.Mutation{}

	if v, ok := params.Handle.Get(); ok {
		opts = append(opts, account.SetHandle(v))
	}
	if v, ok := params.Name.Get(); ok {
		opts = append(opts, account.SetName(v))
	}
	if v, ok := params.Bio.Get(); ok {
		opts = append(opts, account.SetBio(v))
	}
	if v, ok := params.Interests.Get(); ok {
		opts = append(opts, account.SetInterests(v))
	}
	if v, ok := params.Links.Get(); ok {
		opts = append(opts, account.SetLinks(v))
	}
	if v, ok := params.Meta.Get(); ok {
		opts = append(opts, account.SetMetadata(v))
	}

	acc, err := s.account_repo.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexProfile{
		ID: id,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	return acc, nil
}
