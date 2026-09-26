package node_mutate

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

const MaxBatchItems = 100

type BatchFields struct {
	Name       opt.Optional[string]
	Slug       opt.Optional[mark.Slug]
	Content    opt.Optional[datagraph.Content]
	Parent     opt.Optional[library.QueryKey]
	URL        opt.Optional[url.URL]
	Tags       opt.Optional[tag_ref.Names]
	Visibility opt.Optional[visibility.Visibility]
}

type BatchCreate struct {
	ID     library.NodeID
	Fields BatchFields
}

type BatchUpdate struct {
	Key    library.QueryKey
	Fields BatchFields
}

type BatchResult struct {
	Mark library.Mark
	Name string
}

func (f BatchFields) versioned() bool {
	return f.Name.Ok() || f.Slug.Ok() || f.Content.Ok()
}

func validateBatchFields(fields []BatchFields) error {
	if len(fields) == 0 || len(fields) > MaxBatchItems {
		return fmt.Errorf("batch must contain 1 to %d items", MaxBatchItems)
	}

	for i, field := range fields {
		if err := field.validate(); err != nil {
			return fmt.Errorf("item %d: %w", i+1, err)
		}
	}

	return nil
}

func (f BatchFields) validate() error {
	if name, ok := f.Name.Get(); ok && strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if slug, ok := f.Slug.Get(); ok && slug.String() == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if parent, ok := f.Parent.Get(); ok && parent.String() == "" {
		return fmt.Errorf("parent must identify a page")
	}

	if u, ok := f.URL.Get(); ok && ((u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "") {
		return fmt.Errorf("URL must be an absolute HTTP or HTTPS URL")
	}

	if v, ok := f.Visibility.Get(); ok && v != visibility.VisibilityPublished && v != visibility.VisibilityDraft {
		return fmt.Errorf("visibility must be published or draft")
	}

	for _, name := range f.Tags.OrZero() {
		if strings.TrimSpace(name.String()) == "" {
			return fmt.Errorf("tag names cannot be empty")
		}
	}

	return nil
}

func indexMutationStates(states map[library.NodeID]node_querier.MutationState) map[string]library.NodeID {
	index := make(map[string]library.NodeID, len(states))
	for id, state := range states {
		index[state.Mark.Slug()] = id
	}

	return index
}

func resolveBatchKey(key library.QueryKey, states map[library.NodeID]node_querier.MutationState, slugs map[string]library.NodeID) (library.NodeID, error) {
	if id, ok := key.ID().Get(); ok {
		if _, exists := states[library.NodeID(id)]; exists {
			return library.NodeID(id), nil
		}
	} else if id, ok := slugs[key.String()]; ok {
		return id, nil
	}

	return library.NodeID{}, fmt.Errorf("page %q does not exist", key.String())
}

func validateBatchHierarchy(parents map[library.NodeID]library.NodeID) error {
	complete := make(map[library.NodeID]bool, len(parents))
	for start := range parents {
		path := make(map[library.NodeID]bool)
		for current := start; current != (library.NodeID{}) && !complete[current]; current = parents[current] {
			if path[current] {
				return fmt.Errorf("parent references contain a cycle at page %s", current)
			}

			path[current] = true
		}

		for id := range path {
			complete[id] = true
		}
	}

	return nil
}

func (s *Manager) batchOptions(ctx context.Context, fields []BatchFields, current []node_querier.MutationState, parents []opt.Optional[library.NodeID]) [][]node_writer.Option {
	options := make([][]node_writer.Option, len(fields))
	links := make(map[string]opt.Optional[xid.ID])
	for i, field := range fields {
		opts := []node_writer.Option{}
		field.Name.Call(func(v string) { opts = append(opts, node_writer.WithName(v)) })
		field.Slug.Call(func(v mark.Slug) { opts = append(opts, node_writer.WithSlug(v.String())) })
		field.Visibility.Call(func(v visibility.Visibility) { opts = append(opts, node_writer.WithVisibility(v)) })
		parents[i].Call(func(v library.NodeID) { opts = append(opts, node_writer.WithParent(v)) })
		field.Content.Call(func(v datagraph.Content) {
			content := contentWithStableBlocks(current[i].Content, v, s.logger)
			opts = append(opts, node_writer.WithContent(content))
		})

		if u, ok := field.URL.Get(); ok {
			link, fetched := links[u.String()]
			if !fetched {
				value, _, err := s.fetcher.ScrapeAndStore(ctx, u)
				if err == nil {
					link = opt.New(xid.ID(value.ID))
				}

				links[u.String()] = link
			}

			link.Call(func(id xid.ID) { opts = append(opts, node_writer.WithLink(id)) })
		}

		options[i] = opts
	}

	return options
}
