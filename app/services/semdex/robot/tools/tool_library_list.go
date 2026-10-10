package tools

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/lib/mcp"
)

func libraryPageListParams(args mcp.ToolLibraryPageTreeInput) (node_traversal.ListParams, pagination.Parameters, error) {
	params := node_traversal.ListParams{
		CreatedSince:  opt.NewPtr(args.CreatedSince),
		CreatedBefore: opt.NewPtr(args.CreatedBefore),
	}
	page := opt.NewPtr(args.Page).Or(1)
	size := opt.NewPtr(args.PageSize).Or(10)
	if page < 1 || size < 1 || size > 100 {
		return params, pagination.Parameters{}, fault.New("page must be positive and page_size must be between 1 and 100", ftag.With(ftag.InvalidArgument))
	}
	pp := pagination.NewPageParams(uint(page), uint(size))
	if args.Depth != nil && *args.Depth < -1 {
		return params, pp, fault.New("depth must be -1 or greater", ftag.With(ftag.InvalidArgument))
	}
	if args.CreatedSince != nil && args.CreatedBefore != nil && !args.CreatedSince.Before(*args.CreatedBefore) {
		return params, pp, fault.New("created_since must be before created_before", ftag.With(ftag.InvalidArgument))
	}
	if args.Sort != nil {
		switch *args.Sort {
		case mcp.ToolLibraryPageTreeInputSortNewest:
			params.NewestFirst = opt.New(true)
		case mcp.ToolLibraryPageTreeInputSortOldest:
			params.NewestFirst = opt.New(false)
		case mcp.ToolLibraryPageTreeInputSortHierarchy:
		default:
			return params, pp, fault.New("sort must be hierarchy, newest or oldest", ftag.With(ftag.InvalidArgument))
		}
	}
	return params, pp, nil
}
