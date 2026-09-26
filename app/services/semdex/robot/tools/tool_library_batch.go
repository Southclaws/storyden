package tools

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/lib/mcp"
)

func (lt *libraryTools) newLibraryPagesCreateTool() *Tool {
	def := mcp.GetLibraryPagesCreateTool()

	return &Tool{
		Definition: def,
		Handler:    makeHandler(lt.ExecuteLibraryPagesCreate),
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name:        def.Name,
				Description: def.Description,
				InputSchema: def.InputSchema,
			}, func(ctx agent.Context, input mcp.ToolLibraryPagesCreateInput) (*mcp.ToolLibraryPagesCreateOutput, error) {
				return lt.ExecuteLibraryPagesCreate(ctx, input)
			})
		},
	}
}

func (lt *libraryTools) newLibraryPagesUpdateTool() *Tool {
	def := mcp.GetLibraryPagesUpdateTool()

	return &Tool{
		Definition: def,
		Handler:    makeHandler(lt.ExecuteLibraryPagesUpdate),
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name:        def.Name,
				Description: def.Description,
				InputSchema: def.InputSchema,
			}, func(ctx agent.Context, input mcp.ToolLibraryPagesUpdateInput) (*mcp.ToolLibraryPagesUpdateOutput, error) {
				return lt.ExecuteLibraryPagesUpdate(ctx, input)
			})
		},
	}
}

func (lt *libraryTools) ExecuteLibraryPagesCreate(ctx context.Context, input mcp.ToolLibraryPagesCreateInput) (*mcp.ToolLibraryPagesCreateOutput, error) {
	refs := make([]string, len(input.Items))
	for i, item := range input.Items {
		refs[i] = item.Ref
	}

	if err := validateLibraryBatchRefs(refs); err != nil {
		return nil, err
	}

	ids := make(map[string]library.NodeID, len(refs))
	for _, ref := range refs {
		ids[ref] = library.NodeID(xid.New())
	}

	items := make([]node_mutate.BatchCreate, len(input.Items))
	for i, item := range input.Items {
		inputFields := libraryBatchFields{
			Name:    &item.Name,
			Slug:    item.Slug,
			Content: item.Content,
			Parent:  item.Parent,
			URL:     item.Url,
			Tags:    item.Tags,
		}
		if item.Visibility != nil {
			value := string(*item.Visibility)
			inputFields.Visibility = &value
		}

		fields, err := inputFields.parse()
		if err != nil {
			return nil, fmt.Errorf("item %q: %w", item.Ref, err)
		}

		if item.ParentRef != nil {
			if item.Parent != nil {
				return nil, fmt.Errorf("item %q cannot set both parent and parent_ref", item.Ref)
			}

			id, ok := ids[*item.ParentRef]
			if !ok {
				return nil, fmt.Errorf("item %q references missing parent %q", item.Ref, *item.ParentRef)
			}

			fields.Parent = opt.New(library.NewID(xid.ID(id)))
		}

		fields.Visibility = opt.New(fields.Visibility.Or(visibility.VisibilityPublished))
		items[i] = node_mutate.BatchCreate{ID: ids[item.Ref], Fields: fields}
	}

	results, err := lt.nodeMutator.CreateMany(ctx, items)
	if err != nil {
		return nil, err
	}

	return &mcp.ToolLibraryPagesCreateOutput{
		Results:    lt.batchResults(refs, results, mcp.LibraryPageBatchResultYamlStatusCreated),
		NextAction: "All pages were created. Report them using their browser URLs.",
	}, nil
}

func (lt *libraryTools) ExecuteLibraryPagesUpdate(ctx context.Context, input mcp.ToolLibraryPagesUpdateInput) (*mcp.ToolLibraryPagesUpdateOutput, error) {
	refs := make([]string, len(input.Items))
	for i, item := range input.Items {
		refs[i] = item.Ref
	}

	if err := validateLibraryBatchRefs(refs); err != nil {
		return nil, err
	}

	items := make([]node_mutate.BatchUpdate, len(input.Items))
	for i, item := range input.Items {
		inputFields := libraryBatchFields{
			Name:    item.Name,
			Slug:    item.Slug,
			Content: item.Content,
			Parent:  item.Parent,
			URL:     item.Url,
			Tags:    item.Tags,
		}
		if item.Visibility != nil {
			value := string(*item.Visibility)
			inputFields.Visibility = &value
		}

		fields, err := inputFields.parse()
		if err != nil {
			return nil, fmt.Errorf("item %q: %w", item.Ref, err)
		}

		items[i] = node_mutate.BatchUpdate{Key: library.NewKey(item.Id), Fields: fields}
	}

	results, err := lt.nodeMutator.UpdateMany(ctx, items)
	if err != nil {
		return nil, err
	}

	return &mcp.ToolLibraryPagesUpdateOutput{
		Results:    lt.batchResults(refs, results, mcp.LibraryPageBatchResultYamlStatusUpdated),
		NextAction: "All pages were updated. Report them using their browser URLs.",
	}, nil
}

func validateLibraryBatchRefs(refs []string) error {
	if len(refs) == 0 || len(refs) > node_mutate.MaxBatchItems {
		return fmt.Errorf("batch must contain 1 to %d items", node_mutate.MaxBatchItems)
	}

	seen := make(map[string]bool, len(refs))
	for _, ref := range refs {
		if strings.TrimSpace(ref) == "" || len(ref) > 128 {
			return fmt.Errorf("item ref must contain 1 to 128 characters")
		}

		if seen[ref] {
			return fmt.Errorf("duplicate item ref %q", ref)
		}

		seen[ref] = true
	}

	return nil
}

type libraryBatchFields struct {
	Name       *string
	Slug       *string
	Content    *string
	Parent     *string
	URL        *string
	Tags       []string
	Visibility *string
}

func (f libraryBatchFields) parse() (node_mutate.BatchFields, error) {
	fields := node_mutate.BatchFields{Name: opt.NewPtr(f.Name)}
	if f.Slug != nil {
		slug, err := mark.NewSlug(*f.Slug)
		if err != nil {
			return fields, err
		}

		fields.Slug = opt.New(*slug)
	}

	if f.Content != nil {
		content, err := datagraph.NewRichText(*f.Content)
		if err != nil {
			return fields, err
		}

		fields.Content = opt.New(content)
	}

	if f.Parent != nil {
		fields.Parent = opt.New(library.NewKey(*f.Parent))
	}

	if f.URL != nil {
		value, err := url.Parse(*f.URL)
		if err != nil {
			return fields, err
		}

		fields.URL = opt.New(*value)
	}

	if f.Visibility != nil {
		value, err := visibility.NewVisibility(*f.Visibility)
		if err != nil {
			return fields, err
		}

		fields.Visibility = opt.New(value)
	}

	if f.Tags != nil {
		names := make(tag_ref.Names, len(f.Tags))
		for i, name := range f.Tags {
			names[i] = tag_ref.NewName(name)
		}

		fields.Tags = opt.New(names)
	}

	return fields, nil
}

func (lt *libraryTools) batchResults(refs []string, results []node_mutate.BatchResult, status mcp.LibraryPageBatchResultYamlStatus) []mcp.LibraryPageBatchResultYaml {
	mapped := make([]mcp.LibraryPageBatchResultYaml, len(results))
	for i, result := range results {
		id, slug := result.Mark.ID().String(), result.Mark.Slug()
		browserURL := datagraph.CanonicalResolveURL(lt.webAddress, datagraph.KindNode, id+"-"+slug).String()
		mapped[i] = mcp.LibraryPageBatchResultYaml{
			Ref:        refs[i],
			Status:     status,
			Id:         id,
			Slug:       slug,
			Name:       result.Name,
			BrowserUrl: browserURL,
		}
	}

	return mapped
}
