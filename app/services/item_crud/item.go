package item_crud

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/hydrator"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		p Partial,
	) (*datagraph.Item, error)
	Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
	Update(ctx context.Context, slug datagraph.ItemSlug, p Partial) (*datagraph.Item, error)
	Archive(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	ImageURL    opt.Optional[string]
	URL         opt.Optional[string]
	Description opt.Optional[string]
	Content     opt.Optional[string]
	Properties  opt.Optional[any]
}

func (p Partial) Opts() (opts []item.Option) {
	p.Name.Call(func(value string) { opts = append(opts, item.WithName(value)) })
	p.Slug.Call(func(value string) { opts = append(opts, item.WithSlug(value)) })
	p.ImageURL.Call(func(value string) { opts = append(opts, item.WithImageURL(value)) })
	p.Description.Call(func(value string) { opts = append(opts, item.WithDescription(value)) })
	p.Content.Call(func(value string) { opts = append(opts, item.WithContent(value)) })
	p.Properties.Call(func(value any) { opts = append(opts, item.WithProperties(value)) })
	return
}

type service struct {
	l        *zap.Logger
	cr       item.Repository
	hydrator hydrator.Service
}

func New(
	l *zap.Logger,
	cr item.Repository,
	hydrator hydrator.Service,
) Manager {
	return &service{
		l:        l.With(zap.String("service", "cluster")),
		cr:       cr,
		hydrator: hydrator,
	}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	slug string,
	desc string,
	p Partial,
) (*datagraph.Item, error) {
	opts := p.Opts()
	opts = append(opts, s.hydrateLink(ctx, p)...)

	itm, err := s.cr.Create(ctx, owner, name, slug, desc, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}

func (s *service) Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error) {
	itm, err := s.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}

func (s *service) Update(ctx context.Context, slug datagraph.ItemSlug, p Partial) (*datagraph.Item, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	itm, err := s.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !itm.Owner.Admin {
		if itm.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	opts := p.Opts()
	opts = append(opts, s.hydrateLink(ctx, p)...)

	itm, err = s.cr.Update(ctx, itm.ID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}

func (s *service) Archive(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error) {
	itm, err := s.cr.Archive(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}

func (s *service) hydrateLink(ctx context.Context, partial Partial) (opts []item.Option) {
	text, textOK := partial.Content.Get()

	if !textOK && !partial.URL.Ok() {
		return
	}

	return s.hydrator.HydrateItem(ctx, text, partial.URL)
}
