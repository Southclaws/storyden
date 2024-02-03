package datagraph_searcher

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/rbac"
)

type Result struct {
	Id   xid.ID
	Name string
	Type string
}

type Searcher interface {
	Search(ctx context.Context, query string) ([]*Result, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	wc *weaviate.Client
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	wc *weaviate.Client,
) Searcher {
	return &service{
		l:    l.With(zap.String("service", "search")),
		rbac: rbac,
		wc:   wc,
	}
}

type WeaviateObject struct {
	DatagraphID   string `json:"datagraph_id"`
	DatagraphType string `json:"datagraph_type"`
	Name          string `json:"name"`
	Content       string `json:"content"`
}

type WeaviateContent struct {
	Content []WeaviateObject
}

type WeaviateResponse struct {
	Get WeaviateContent
}

func (s *service) Search(ctx context.Context, q string) ([]*Result, error) {
	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "name"},
		{Name: "content"},
	}

	arg := s.wc.GraphQL().
		HybridArgumentBuilder().
		WithAlpha(0.25).
		WithFusionType(graphql.Ranked).
		WithQuery(q)

	result, err := s.wc.GraphQL().Get().
		WithClassName("Content").
		WithFields(fields...).
		WithHybrid(arg).
		WithLimit(30).
		Do(context.Background())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	j, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	parsed := WeaviateResponse{}
	err = json.Unmarshal(j, &parsed)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results, err := dt.MapErr(parsed.Get.Content, func(v WeaviateObject) (*Result, error) {
		id, err := xid.FromString(v.DatagraphID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return &Result{
			Id:   id,
			Type: v.DatagraphType,
			Name: v.Name,
		}, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return results, nil
}
