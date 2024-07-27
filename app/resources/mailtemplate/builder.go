package mailtemplate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/matcornic/hermes/v2"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
)

type Rendered struct {
	HTML  string
	Plain string
}

type Builder struct {
	instanceURL string
	settings    settings.Repository
}

func New(
	ctx context.Context,
	cfg config.Config,
	set settings.Repository,
) (*Builder, error) {
	return &Builder{
		instanceURL: cfg.PublicWebAddress,
		settings:    set,
	}, nil
}

func (b *Builder) Build(ctx context.Context, name string, intros []string, actions []hermes.Action) (*Rendered, error) {
	s, err := b.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	instanceTitle := s.Title.Get()
	instanceURL := b.instanceURL

	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      instanceTitle,
			Link:      instanceURL,
			Copyright: "-",
		},
	}

	template := hermes.Email{
		Body: hermes.Body{
			Name:      name,
			Intros:    intros,
			Actions:   actions,
			Outros:    []string{},
			Signature: "Thanks",
		},
	}

	html, err := h.GenerateHTML(template)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	plain, err := h.GeneratePlainText(template)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Rendered{html, plain}, nil
}
