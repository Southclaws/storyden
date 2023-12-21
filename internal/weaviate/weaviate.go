package weaviate

import (
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/kelseyhightower/envconfig"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"go.uber.org/fx"
)

type Configuration struct {
	Enabled   bool   `envconfig:"WEAVIATE_ENABLED"`
	Host      string `envconfig:"WEAVIATE_HOST"`
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

	wc := weaviate.Config{
		Host:       cfg.Host,
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: cfg.Token},
		Headers:    map[string]string{"X-OpenAI-Api-Key": cfg.OpenAIKey},
	}

	client, err := weaviate.NewClient(wc)
	if err != nil {
		fmt.Println(err)
	}

	return client, nil
}
