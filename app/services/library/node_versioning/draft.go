package node_versioning

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/deletable"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type DraftPartial struct {
	Name               opt.Optional[string]
	Slug               opt.Optional[mark.Slug]
	Description        deletable.Value[string]
	Content            deletable.Value[datagraph.Content]
	PropertiesSnapshot opt.Optional[library.PropertyMutationList]
	Metadata           opt.Optional[map[string]any]
}

type draftSnapshot struct {
	Name        string
	Slug        string
	Description opt.Optional[string]
	Content     opt.Optional[datagraph.Content]
	Properties  []node_version.PropertySnapshot
	Metadata    map[string]any
}

func (s *Service) CreateDraft(ctx context.Context, qk library.QueryKey, p DraftPartial) (*node_version.NodeVersion, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionSubmitLibraryNodeChanges, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("cannot create draft", "You need permission to create draft versions."),
		)
	}

	acc, err := session.GetAccount(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodeFilter := versionNodeFilter(ctx)

	existing, err := s.versionQuerier.GetNodeDraft(ctx, qk, nodeFilter)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if existing.Ok() {
		var errorMessage string
		if acc.Roles.Permissions().HasAll(rbac.PermissionManageLibrary) {
			errorMessage = "This page already has a working draft. Apply or delete the draft before creating another one."
		} else {
			errorMessage = "This page already has a working draft. It must be applied by a member with permissions before you can create a draft."
		}

		return nil, fault.New("node already has a draft version",
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc("draft exists", errorMessage),
		)
	}

	n, err := s.nodeReader.GetBySlug(ctx, qk, opt.NewEmpty[node_querier.ChildSortRule]())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	snapshot := snapshotFromNode(n)
	if err := applyDraftPartial(&snapshot, p); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	version, err := s.versionWriter.Create(ctx, qk, acc.ID, nodeFilter, snapshotOptions(snapshot)...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeVersionDraftCreated{
		NodeID:    version.NodeID,
		NodeSlug:  n.GetSlug(),
		VersionID: version.ID.String(),
		AuthorID:  acc.ID,
	})

	return version, nil
}

func (s *Service) UpdateDraft(ctx context.Context, qk library.QueryKey, id node_version.VersionID, p DraftPartial) (*node_version.NodeVersion, error) {
	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	v, err := s.versionQuerier.GetForNode(ctx, qk, id, versionNodeFilter(ctx))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil, versionNotFound(ctx)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.updateDraft(ctx, v, callerID, p)
}

func (s *Service) updateDraft(ctx context.Context, v *node_version.NodeVersion, callerID account.AccountID, p DraftPartial) (*node_version.NodeVersion, error) {
	if err := authoriseVisibleDraftMutation(ctx, callerID, v); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	snapshot := snapshotFromVersion(v)
	if err := applyDraftPartial(&snapshot, p); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := s.versionWriter.Update(ctx, v.ID, snapshotOptions(snapshot)...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeVersionDraftUpdated{
		NodeID:    updated.NodeID,
		NodeSlug:  updated.Slug,
		VersionID: updated.ID.String(),
		AuthorID:  updated.Author.ID,
	})

	return updated, nil
}

func (s *Service) GetVisibleDraft(ctx context.Context, qk library.QueryKey) (*node_version.NodeVersion, error) {
	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, draftNotFound(ctx)
	}

	return s.getVisibleDraft(ctx, qk, callerID)
}

func (s *Service) getVisibleDraft(ctx context.Context, qk library.QueryKey, callerID account.AccountID) (*node_version.NodeVersion, error) {
	nodeFilter := versionNodeFilter(ctx)

	draft, err := s.versionQuerier.GetNodeDraft(ctx, qk, nodeFilter)
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil, draftNotFound(ctx)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	v, ok := draft.Get()
	if !ok {
		return nil, draftNotFound(ctx)
	}

	if err := authoriseDraftVisible(ctx, callerID, &v); err != nil {
		return nil, draftNotFound(ctx)
	}

	return &v, nil
}

func (s *Service) UpdateVisibleDraft(ctx context.Context, qk library.QueryKey, p DraftPartial) (*node_version.NodeVersion, error) {
	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	v, err := s.getVisibleDraft(ctx, qk, callerID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.updateDraft(ctx, v, callerID, p)
}

func (s *Service) ListVisible(ctx context.Context, qk library.QueryKey, page pagination.Parameters) (pagination.Result[*node_version.NodeVersion], error) {
	versions, err := s.versionQuerier.ListVisible(ctx, qk, versionNodeFilter(ctx), page)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return versions, nil
}

func (s *Service) ListAllDrafts(ctx context.Context, page pagination.Parameters) (pagination.Result[*node_version.NodeVersionWithNode], error) {
	if _, err := session.GetAccountID(ctx); err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("sign in required", "Sign in to list draft versions."),
		)
	}

	drafts, err := s.versionQuerier.ListAllDrafts(ctx, versionNodeFilter(ctx), page)
	if err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return drafts, nil
}

func (s *Service) GetVisible(ctx context.Context, qk library.QueryKey, id node_version.VersionID) (*node_version.NodeVersion, error) {
	v, err := s.versionQuerier.GetForNode(ctx, qk, id, versionNodeFilter(ctx))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil, versionNotFound(ctx)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v.Status == node_version.VersionStatusApplied {
		previous, err := s.versionQuerier.GetPreviousReference(ctx, v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		v.Previous = previous

		return v, nil
	}

	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, versionNotFound(ctx)
	}

	if err := session.Authorise(ctx, func() error {
		if v.Author.ID != callerID {
			return fault.New("not the version author",
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("not author", "You can only view your own draft versions."),
			)
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, versionNotFound(ctx)
	}

	return v, nil
}

func (s *Service) DiscardDraft(ctx context.Context, qk library.QueryKey, id node_version.VersionID) error {
	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("sign in required", "Sign in to discard a draft."),
		)
	}

	v, err := s.versionQuerier.GetForNode(ctx, qk, id, versionNodeFilter(ctx))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return versionNotFound(ctx)
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := authoriseVisibleDraftDelete(ctx, callerID, v); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.versionWriter.Delete(ctx, id); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeVersionDraftDeleted{
		NodeID:    v.NodeID,
		NodeSlug:  v.Slug,
		VersionID: v.ID.String(),
		AuthorID:  v.Author.ID,
	})

	return nil
}

func snapshotFromNode(n *library.Node) draftSnapshot {
	return draftSnapshot{
		Name:        n.Name,
		Slug:        n.GetSlug(),
		Description: n.Description,
		Content:     n.Content,
		Properties:  propertyTableToSnapshot(n.Properties),
		Metadata:    n.Metadata,
	}
}

func snapshotFromVersion(v *node_version.NodeVersion) draftSnapshot {
	return draftSnapshot{
		Name:        v.Name,
		Slug:        v.Slug,
		Description: v.Description,
		Content:     v.Content,
		Properties:  v.PropertiesSnapshot.Or([]node_version.PropertySnapshot{}),
		Metadata:    v.Metadata,
	}
}

func applyDraftPartial(snapshot *draftSnapshot, p DraftPartial) error {
	p.Name.Call(func(value string) { snapshot.Name = value })
	p.Slug.Call(func(value mark.Slug) { snapshot.Slug = value.String() })

	description, clearDescription := p.Description.Get()
	if clearDescription {
		snapshot.Description = opt.NewEmpty[string]()
	} else {
		description.Call(func(value string) { snapshot.Description = opt.New(value) })
	}

	content, clearContent := p.Content.Get()
	if clearContent {
		snapshot.Content = opt.NewEmpty[datagraph.Content]()
	} else {
		var err error
		content.Call(func(value datagraph.Content) {
			var stable datagraph.ContentWithBlocks
			if previous, ok := snapshot.Content.Get(); ok {
				stable, err = datagraph.NewRichTextWithChangedBlocks(previous, value)
			} else {
				stable, err = datagraph.NewRichTextWithNewBlocks(value)
			}
			if err != nil {
				return
			}
			snapshot.Content = opt.New(stable.Content)
		})
		if err != nil {
			return err
		}
	}

	if props, ok := p.PropertiesSnapshot.Get(); ok {
		snapshot.Properties = propertyMutationsToSnapshot(props)
	}

	p.Metadata.Call(func(value map[string]any) { snapshot.Metadata = value })

	return nil
}

func snapshotOptions(snapshot draftSnapshot) []node_version_writer.Option {
	return []node_version_writer.Option{
		node_version_writer.WithName(snapshot.Name),
		node_version_writer.WithSlug(snapshot.Slug),
		node_version_writer.WithDescription(snapshot.Description),
		node_version_writer.WithContent(snapshot.Content),
		node_version_writer.WithPropertiesSnapshot(snapshot.Properties),
		node_version_writer.WithMetadata(snapshot.Metadata),
	}
}

func draftNotFound(ctx context.Context) error {
	return fault.New("node draft version not found",
		fctx.With(ctx),
		ftag.With(ftag.NotFound),
		fmsg.WithDesc("not found", "The requested draft does not exist or is not visible to you."),
	)
}

func versionNotFound(ctx context.Context) error {
	return fault.New("node version not found",
		fctx.With(ctx),
		ftag.With(ftag.NotFound),
		fmsg.WithDesc("not found", "The requested version does not exist or is not visible to you."),
	)
}

func versionNodeFilter(ctx context.Context) node_version_querier.NodeFilter {
	return node_version_querier.NodeFilter{
		AccountID: session.GetOptAccountID(ctx),
		CanManage: session.Authorise(ctx, nil, rbac.PermissionManageLibrary) == nil,
	}
}
