package pinecone

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

type Client struct {
	*pinecone.Client
	size   int32
	cloud  pinecone.Cloud
	region string
}

type Index = pinecone.IndexConnection

type Vector = pinecone.Vector

type Metadata = pinecone.Metadata

type MetadataFilter = pinecone.MetadataFilter

type QueryByVectorValuesRequest = pinecone.QueryByVectorValuesRequest

type ScoredVector = pinecone.ScoredVector

func Build() fx.Option {
	return fx.Provide(newPinecone)
}

func newPinecone(cfg config.Config) (*Client, error) {
	c, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: cfg.PineconeAPIKey,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: c,
		size:   cfg.PineconeDimentions,
		cloud:  pinecone.Cloud(cfg.PineconeCloud),
		region: cfg.PineconeRegion,
	}, nil
}

func (c *Client) GetOrCreateIndex(ctx context.Context, name string) (*Index, error) {
	desc, err := func() (*pinecone.Index, error) {
		index, err := c.DescribeIndex(ctx, name)
		if err == nil {
			return index, nil
		}

		if !isNotFound(err) {
			return nil, err
		}

		index, err = c.CreateServerlessIndex(ctx, &pinecone.CreateServerlessIndexRequest{
			Name:      name,
			Dimension: c.size,
			Metric:    "cosine",
			Cloud:     c.cloud,
			Region:    c.region,
		})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return index, nil
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	idxConnection, err := c.Index(pinecone.NewIndexConnParams{Host: desc.Host, Namespace: "storyden"})
	if err != nil {
		return nil, err
	}

	return idxConnection, nil
}

func isNotFound(err error) bool {
	pe := &pinecone.PineconeError{}
	if errors.As(err, &pe) {
		if pe.Code == 404 {
			return true
		}
	}

	return false
}
