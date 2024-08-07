package weaviate

import (
	"context"
	"net/url"
	"reflect"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/fx"
)

type WeaviateClassName string

func (w WeaviateClassName) String() string {
	return string(w)
}

type Configuration struct {
	Enabled   bool   `envconfig:"WEAVIATE_ENABLED"`
	URL       string `envconfig:"WEAVIATE_URL"`
	Token     string `envconfig:"WEAVIATE_API_TOKEN"`
	ClassName string `envconfig:"WEAVIATE_CLASS_NAME"`
	OpenAIKey string `envconfig:"OPENAI_API_KEY"`
}

func Build() fx.Option {
	return fx.Provide(newWeaviateClient)
}

func newWeaviateClient(lc fx.Lifecycle) (*weaviate.Client, WeaviateClassName, error) {
	cfg := Configuration{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, "", fault.Wrap(err)
	}

	if !cfg.Enabled {
		return nil, "", nil
	}

	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, "", fault.Wrap(err)
	}

	wc := weaviate.Config{
		Host:       u.Host,
		Scheme:     u.Scheme,
		AuthConfig: auth.ApiKey{Value: cfg.Token},
		Headers:    map[string]string{"X-OpenAI-Api-Key": cfg.OpenAIKey},
	}

	client, err := weaviate.NewClient(wc)
	if err != nil {
		return nil, "", fault.Wrap(err, fmsg.With("failed to create weaviate client"))
	}

	classMap := map[string]models.Class{
		"text2vec-transformers": {
			Class:      "ContentText2vecTransformers",
			Vectorizer: "text2vec-transformers",
			ModuleConfig: map[string]any{
				"text2vec-transformers": map[string]any{},
			},
		},
		"text2vec-openai": {
			Class:      "ContentOpenAI",
			Vectorizer: "text2vec-openai",
			Properties: []*models.Property{
				{
					Name:     "datagraph_id",
					DataType: []string{"text"},
					ModuleConfig: map[string]any{
						"text2vec-openai": map[string]any{
							"skip": false,
						},
					},
				},
				{
					Name:     "datagraph_type",
					DataType: []string{"text"},
					ModuleConfig: map[string]any{
						"text2vec-openai": map[string]any{
							"skip": false,
						},
					},
				},
				{
					Name:     "name",
					DataType: []string{"text"},
				},
				{
					Name:     "description",
					DataType: []string{"text"},
				},
				{
					Name:     "content",
					DataType: []string{"text"},
				},
			},
			ModuleConfig: map[string]any{
				"text2vec-openai": map[string]any{
					"model":      "text-embedding-3-large",
					"dimensions": "3072",
					"type":       "text",
				},
			},
		},
	}

	if cfg.ClassName == "text2vec-openai" && cfg.OpenAIKey == "" {
		return nil, "", fault.New("OpenAI API key is required for text2vec-openai class")
	}

	class, ok := classMap[cfg.ClassName]
	if !ok {
		return nil, "", fault.New("invalid class name")
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		r, err := client.Schema().
			ClassExistenceChecker().
			WithClassName(class.Class).
			Do(ctx)
		if err != nil {
			return fault.Wrap(err)
		}

		if !r {
			err := client.Schema().
				ClassCreator().
				WithClass(&class).
				Do(ctx)
			if err != nil {
				return fault.Wrap(err)
			}
		} else {
			//
			// MASSIVE HACK WARNING
			//
			// Weaviate does not support updating class properties but currently
			// it's still experimental, so the class structure MAY change in the
			// future so what happens here is the class is deleted which deletes
			// ALL vectorised content. Once this happens after a successful boot
			// the reindex job will re-index all content. On large instances, it
			// will be EXPENSIVE! So, until we settle on the properties, beware.
			//

			current, err := client.Schema().
				ClassGetter().
				WithClassName(class.Class).
				Do(ctx)
			if err != nil {
				return err
			}

			same := compareClassConfig(cfg.ClassName, *current, class)
			if !same {
				err = client.Schema().
					ClassDeleter().
					WithClassName(class.Class).
					Do(ctx)
				if err != nil {
					return fault.Wrap(err)
				}

				err = client.Schema().
					ClassCreator().
					WithClass(&class).
					Do(ctx)
				if err != nil {
					return fault.Wrap(err)
				}
			}
		}

		return nil
	}))

	return client, WeaviateClassName(class.Class), nil
}

func compareClassConfig(cn string, a, b models.Class) bool {
	if a.Class != b.Class {
		return false
	}

	if a.Vectorizer != b.Vectorizer {
		return false
	}

	if len(a.Properties) != len(b.Properties) {
		return false
	}

	for i := range a.Properties {
		if !comparePropertyConfig(cn, a.Properties[i], b.Properties[i]) {
			return false
		}
	}

	return true
}

func comparePropertyConfig(cn string, a, b *models.Property) bool {
	if a.Name != b.Name {
		return false
	}

	if !reflect.DeepEqual(a.DataType, b.DataType) {
		return false
	}

	if a.ModuleConfig != nil && b.ModuleConfig != nil {
		mca := a.ModuleConfig.(map[string]any)[cn].(map[string]any)
		mcb := b.ModuleConfig.(map[string]any)[cn].(map[string]any)

		if mca["skip"] != mcb["skip"] {
			return false
		}
	}

	return true
}
