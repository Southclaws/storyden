package weaviate

import (
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"go.uber.org/fx"
)

type Configuration struct {
	Enabled   bool   `envconfig:"WEAVIATE_ENABLED"`
	URL       string `envconfig:"WEAVIATE_URL"`
	Token     string `envconfig:"WEAVIATE_API_TOKEN"`
	OpenAIKey string `envconfig:"OPENAI_API_KEY"`
}

func Build() fx.Option {
	return fx.Provide(newWeaviateClient)
}

func newWeaviateClient() (*weaviate.Client, error) {
	cfg := Configuration{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fault.Wrap(err)
	}

	if !cfg.Enabled {
		return nil, nil
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	wc := weaviate.Config{
		Host:       u.Host,
		Scheme:     u.Scheme,
		AuthConfig: auth.ApiKey{Value: cfg.Token},
		Headers:    map[string]string{"X-OpenAI-Api-Key": cfg.OpenAIKey},
	}

	client, err := weaviate.NewClient(wc)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create weaviate client"))
	}

	return client, nil
}
