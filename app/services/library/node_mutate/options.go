package node_mutate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_properties"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type preMutationResult struct {
	opts []node_writer.Option

	// Ideally, this API should only return node writer options, but because of
	// a weird public API design choice I made, the PATCH /nodes endpoint also
	// returns tag suggestions which can be opted out of being applied directly.
	// This may change in future but it would require breaking public API change
	// and it works pretty well at the moment as an API design, so not critical.
	tags       opt.Optional[tag_ref.Names]
	title      opt.Optional[string]
	content    opt.Optional[datagraph.Content]
	properties opt.Optional[library.PropertyMutationList]
}

// preMutation constructs node_writer options for a create or partial update.
func (s *Manager) preMutation(ctx context.Context, p Partial, current opt.Optional[library.Node]) (*preMutationResult, error) {
	opts := []node_writer.Option{}

	// Apply all primitive options. These are just basic partial updates.
	p.Name.Call(func(value string) { opts = append(opts, node_writer.WithName(value)) })
	p.Slug.Call(func(value mark.Slug) { opts = append(opts, node_writer.WithSlug(value.String())) })
	p.PrimaryImage.Call(func(value xid.ID) {
		opts = append(opts, node_writer.WithPrimaryImage(value))
	}, func() {
		opts = append(opts, node_writer.WithPrimaryImageRemoved())
	})
	p.Content.Call(func(value datagraph.Content) { opts = append(opts, node_writer.WithContent(value)) })
	p.Metadata.Call(func(value map[string]any) { opts = append(opts, node_writer.WithMetadata(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssetsRemoved(value)) })
	p.Visibility.Call(func(value visibility.Visibility) { opts = append(opts, node_writer.WithVisibility(value)) })
	p.HideChildren.Call(func(value bool) { opts = append(opts, node_writer.WithHideChildren(value)) })

	// If the mutation includes a parent node, we need to query it because the
	// WithParent API only accepts a node ID, not a node mark (slug or ID).
	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.nodeQuerier.Get(ctx, parentSlug)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node_writer.WithParent(library.NodeID(parent.Mark.ID())))
	}

	// If the mutation includes asset sources (so, URLs to assets to be added)
	// download them and append them to the node's asset list.
	if v, ok := p.AssetSources.Get(); ok {
		o, err := s.buildAssetSourcesOpts(ctx, v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		opts = append(opts, o...)
	}

	// -
	// Generative Fill Rules
	// -

	//
	// The content to use during pre-mutation tasks such as tag suggestion, auto
	// title generation and content summarisation. If it's a new node, this will
	// be the content submitted for the new node, if it's an update, either pick
	// the new content if specified in the partial, or the current node content.
	//
	// The first fill rule to run is summarisation, which generates a summary of
	// the content either from a URL or from the current or new content. This
	// mutates the current context content to be the summary thus resulting in
	// title or tag generation to be run on the summary result not the original.
	//

	content := p.Content.Or(current.OrZero().Content.OrZero())

	var contentSuggestion opt.Optional[datagraph.Content]
	var titleSuggestion opt.Optional[string]
	var tagSuggestions opt.Optional[tag_ref.Names]
	var tagsToWrite opt.Optional[tag_ref.Names]

	// If there's a URL being applied, fetch it and if it returns new content,
	// set that as the content to be used for fill rules.
	if u, ok := p.URL.Get(); ok {
		ln, wc, err := s.fetcher.ScrapeAndStore(ctx, u)
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))

			// If the request intends to run fill-rules on URL content, set it.
			if fs, ok := p.FillSource.Get(); ok && fs == asset.FillSourceURL {
				content = wc.Content
				titleSuggestion = opt.New(wc.Title)
			}
		}
	}

	if cfr, ok := p.ContentFill.Get(); ok {
		suggested, err := s.buildSummaryOpts(ctx, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		switch cfr.FillRule {
		case asset.ContentFillRuleQuery:
			content = *suggested
			contentSuggestion = opt.New(*suggested)

		case asset.ContentFillRuleReplace:
			opts = append(opts, node_writer.WithContent(*suggested))
		}
	}

	if tf, ok := p.TitleFill.Get(); ok {
		title, err := s.buildTitleSuggestionOpts(ctx, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if suggested, ok := title.Get(); ok {
			if tf.FillRule == datagraph.TitleFillRuleReplace {
				opts = append(opts, node_writer.WithName(suggested))
			} else {
				titleSuggestion = title
			}
		}
	}

	if tfr, ok := p.TagFill.Get(); ok {
		suggested, err := s.buildTagSuggestionOpts(ctx, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		switch tfr.FillRule {
		case tag.TagFillRuleQuery:
			tagSuggestions = opt.New(suggested)

		case tag.TagFillRuleReplace:
			tagsToWrite = opt.New(suggested)
		}
	} else {
		// If not running a generative fill command use the partial update tags.
		tagsToWrite = p.Tags
	}

	// -
	// Saving tags
	//
	// Happens after any generative options due to how tag-fill rules may yield
	// new tags that need to be saved before being applied to the node. But even
	// if there's no tag-fill rule set this is applied to any tags in the patch.
	// -

	if t, ok := tagsToWrite.Get(); ok {
		n, ok := current.Get()
		if ok {
			tagOpts, err := s.createDeleteTagsForExistingNode(ctx, &n, t)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			opts = append(opts, tagOpts...)
		} else {
			tagOpts, err := s.createDeleteTagsForNewNode(ctx, t)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			opts = append(opts, tagOpts...)
		}
	}

	return &preMutationResult{
		opts:       opts,
		tags:       tagSuggestions,
		title:      titleSuggestion,
		content:    contentSuggestion,
		properties: p.Properties,
	}, nil
}

func (s *Manager) buildAssetSourcesOpts(ctx context.Context, sources []string) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	for _, source := range sources {
		a, err := s.fetcher.CopyAsset(ctx, source)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
	}

	return opts, nil
}

func (s *Manager) createDeleteTagsForNewNode(ctx context.Context, tags tag_ref.Names) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	newTags, err := s.tagWriter.Add(ctx, tags...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	addIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })

	opts = append(opts, node_writer.WithTagsAdd(addIDs...))

	return opts, nil
}

func (s *Manager) createDeleteTagsForExistingNode(ctx context.Context, n *library.Node, tags tag_ref.Names) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	currentTagNames := n.Tags.Names()

	toCreate, toRemove := lo.Difference(tags, currentTagNames)

	newTags, err := s.tagWriter.Add(ctx, toCreate...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	addIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })
	removeIDs := dt.Reduce(n.Tags, func(acc []tag_ref.ID, prev *tag_ref.Tag) []tag_ref.ID {
		if lo.Contains(toRemove, prev.Name) {
			acc = append(acc, prev.ID)
		}
		return acc
	}, []tag_ref.ID{})

	opts = append(opts, node_writer.WithTagsAdd(addIDs...))
	opts = append(opts, node_writer.WithTagsRemove(removeIDs...))

	return opts, nil
}

func (s *Manager) buildTitleSuggestionOpts(ctx context.Context, content datagraph.Content) (opt.Optional[string], error) {
	// Only bother if there's any actual content to work with!
	if content.IsEmpty() {
		return opt.NewEmpty[string](), nil
	}

	titles, err := s.titler.SuggestTitle(ctx, content)
	if err != nil {
		return opt.NewEmpty[string](), fault.Wrap(err, fctx.With(ctx))
	}

	if len(titles) == 0 {
		return opt.NewEmpty[string](), nil
	}

	return opt.New(titles[0]), nil
}

func (s *Manager) buildTagSuggestionOpts(ctx context.Context, content datagraph.Content) (tag_ref.Names, error) {
	// Only bother if there's any actual content to work with!
	if content.IsEmpty() {
		return nil, nil
	}

	gathered, err := s.tagger.Gather(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return gathered, nil
}

func (s *Manager) buildSummaryOpts(ctx context.Context, content datagraph.Content) (*datagraph.Content, error) {
	summary, err := s.summariser.Summarise(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	newContent, err := datagraph.NewRichText(summary)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &newContent, nil
}

type postMutationResult struct {
	properties opt.Optional[library.PropertyTable]
}

// NOTE: "post" as in afterwards, not "post" as in thread/reply post...
func (s *Manager) postMutation(ctx context.Context, n *library.Node, pre *preMutationResult) (*postMutationResult, error) {
	updatedProperties := opt.New(library.PropertyTable{})

	if properties, ok := pre.properties.Get(); ok {
		schema, hasSchema := n.Properties.Get()

		migration, err := schema.Schema.Split(properties)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if !hasSchema {
			mutations, err := dt.MapErr(migration.NewProps, mapNewPropertyMutation)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			newSchema, err := s.schemaWriter.CreateForNode(ctx, library.NodeID(n.Mark.ID()), mutations)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema
		} else {
			schemaUpdates := []*node_properties.SchemaFieldMutation{}

			schemaUpdates = lo.FilterMap(migration.ExistingProps, func(pm *library.ExistingPropertyMutation, _ int) (*node_properties.SchemaFieldMutation, bool) {
				if !pm.IsSchemaChanged {
					return nil, false
				}

				return &node_properties.SchemaFieldMutation{
					ID:   opt.New(pm.ID),
					Name: pm.Name,
					Type: pm.Type,
					Sort: pm.Sort,
				}, true
			})

			newSchema, err := s.schemaWriter.UpdateSiblings(ctx, library.QueryKey{n.Mark.Queryable()}, schemaUpdates)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema
		}

		for _, newProp := range migration.NewProps {
			newSchemaProp, found := lo.Find(schema.Schema.Fields, func(f *library.PropertySchemaField) bool {
				return f.Name == newProp.Name
			})
			if !found {
				continue
			}
			for i, mutProp := range properties {
				if newProp.Name == mutProp.Name {
					properties[i].ID = opt.New(newSchemaProp.ID)
				}
			}
		}

		// TODO: Remove all this code below and move other migrations into the
		// above UpdateSiblings call. Currently that call only does existing
		// field updates but it should perform all migrations. This means we
		// would no longer need to call .Split() twice and mutate properties.

		// re-validate the schema properties mutation plan.
		migration, err = schema.Schema.Split(properties)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if len(migration.NewProps) > 0 {
			newSchemaFields, err := dt.MapErr(migration.NewProps, mapNewPropertyMutation)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			newSchema, err := s.schemaWriter.AddFields(ctx, schema.Schema.ID, newSchemaFields)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema

			for _, newProp := range migration.NewProps {
				newSchemaProp, found := lo.Find(schema.Schema.Fields, func(f *library.PropertySchemaField) bool {
					return f.Name == newProp.Name
				})
				if !found {
					continue
				}

				migration.ExistingProps = append(migration.ExistingProps, &library.ExistingPropertyMutation{
					PropertySchemaField: *newSchemaProp,
					Value:               newProp.Value,
				})
			}
		}

		if len(migration.RemovedProps) > 0 {
			removedSchemaFields, err := dt.MapErr(migration.RemovedProps, mapExistingPropertyMutation)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			newSchema, err := s.schemaWriter.RemoveFields(ctx, schema.Schema.ID, removedSchemaFields)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			schema.Schema = *newSchema
		}

		// Assumption: all schema changes are done by this point. Update no
		// longer needs to actually check the schema, just write the data.
		updated, err := s.propWriter.Update(ctx, library.NodeID(n.GetID()), schema.Schema, migration.ExistingProps)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		updatedProperties = opt.NewPtr(updated)
	}

	return &postMutationResult{
		properties: updatedProperties,
	}, nil
}

func mapNewPropertyMutation(pm *library.PropertyMutation) (*node_properties.SchemaFieldMutation, error) {
	ft, ok := pm.Type.Get()
	if !ok {
		return nil, fault.Wrap(fault.New("no type on new field"), ftag.With(ftag.InvalidArgument), fmsg.WithDesc("missing type", "You must provide a field type when adding a new property."))
	}
	return &node_properties.SchemaFieldMutation{
		Name: pm.Name,
		Type: ft,
		Sort: pm.Sort.OrZero(),
	}, nil
}

func mapExistingPropertyMutation(pm *library.ExistingPropertyMutation) (*node_properties.SchemaFieldMutation, error) {
	return &node_properties.SchemaFieldMutation{
		ID:   opt.New(pm.ID),
		Name: pm.Name,
		Type: pm.Type,
		Sort: pm.Sort,
	}, nil
}
