package resources

import (
	"context"
	"net/url"
	"path"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/transports/http/mcp/mcp_schema"
)

type Provider struct {
	logger      *zap.Logger
	nodeLister  node_traversal.Repository
	nodeQuerier *node_querier.Querier
}

func New(
	logger *zap.Logger,
	nodeLister node_traversal.Repository,
	nodeQuerier *node_querier.Querier,
) *Provider {
	return &Provider{
		logger:      logger,
		nodeLister:  nodeLister,
		nodeQuerier: nodeQuerier,
	}
}

func (p *Provider) ListResources(ctx context.Context) (mcp_schema.ListResourcesResult, error) {
	resources, err := p.getResources(ctx)
	if err != nil {
		return mcp_schema.ListResourcesResult{}, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp_schema.ListResourcesResult{
		Resources: resources,
	}, nil
}

func (p *Provider) ReadResource(ctx context.Context, req mcp_schema.ReadResourceRequest) (mcp_schema.ReadResourceResult, error) {
	nu, err := url.Parse(req.Params.Uri)
	if err != nil {
		return mcp_schema.ReadResourceResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	_, nid := path.Split(nu.Path)

	n, err := p.nodeQuerier.Get(ctx, library.NewKey(nid))
	if err != nil {
		return mcp_schema.ReadResourceResult{}, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp_schema.ReadResourceResult{
		Contents: []any{
			mcp_schema.TextResourceContents{
				Text:     n.Content.OrZero().HTML(),
				MimeType: opt.New("text/html").Ptr(),
				Uri:      req.Params.Uri,
			},
		},
	}, nil
}

func (p *Provider) getResources(ctx context.Context) ([]mcp_schema.Resource, error) {
	nodes, err := p.nodeLister.Subtree(ctx, nil, true)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resources := dt.Map(nodes, func(n *library.Node) mcp_schema.Resource {
		return mcp_schema.Resource{
			Annotations: &mcp_schema.ResourceAnnotations{
				Audience: []mcp_schema.Role{
					mcp_schema.RoleAssistant,
					mcp_schema.RoleUser,
				},
			},
			Description: n.Description.Ptr(),
			MimeType:    opt.New("text/html").Ptr(),
			Name:        n.Name,
			Uri:         "node://" + n.GetID().String(),
		}
	})

	return resources, nil
}
