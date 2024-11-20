package mailtemplate

import (
	"context"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/matcornic/hermes/v2"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

type Action = hermes.Action

type Builder struct {
	instanceURL url.URL
	settings    *settings.SettingsRepository
}

func New(
	ctx context.Context,
	cfg config.Config,
	set *settings.SettingsRepository,
) (*Builder, error) {
	return &Builder{
		instanceURL: cfg.PublicWebAddress,
		settings:    set,
	}, nil
}

func (b *Builder) Build(ctx context.Context, name string, intros []string, actions []hermes.Action) (*mailer.Content, error) {
	s, err := b.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	instanceTitle := s.Title.Or(settings.DefaultTitle)
	instanceURL := b.instanceURL

	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      instanceTitle,
			Link:      instanceURL.String(),
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

	return mailer.NewContent(html, plain)
}
