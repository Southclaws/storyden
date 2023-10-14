package item

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/services/authentication"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Manager interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		opts ...item.Option) (*datagraph.Item, error)
	Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
	Update(ctx context.Context, slug datagraph.ItemSlug, p Partial) (*datagraph.Item, error)
	Archive(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	ImageURL    opt.Optional[string]
	Description opt.Optional[string]
	Properties  opt.Optional[any]
}

type service struct {
	cr item.Repository
}

func New(cr item.Repository) Manager {
	return &service{cr: cr}
}

func (s *service) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	slug string,
	desc string,
	opts ...item.Option,
) (*datagraph.Item, error) {
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
	accountID, err := authentication.GetAccountID(ctx)
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

	opts := []item.Option{}

	p.Name.Call(func(value string) { opts = append(opts, item.WithName(value)) })
	p.Slug.Call(func(value string) { opts = append(opts, item.WithSlug(value)) })
	p.ImageURL.Call(func(value string) { opts = append(opts, item.WithImageURL(value)) })
	p.Description.Call(func(value string) { opts = append(opts, item.WithDescription(value)) })
	p.Properties.Call(func(value any) { opts = append(opts, item.WithProperties(value)) })

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
