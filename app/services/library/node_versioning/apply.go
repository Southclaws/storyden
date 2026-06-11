package node_versioning

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Service) ApplyVersion(ctx context.Context, qk library.QueryKey, id node_version.VersionID) (*node_version.NodeVersion, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	v, err := s.versionQuerier.GetForNode(ctx, qk, id, versionNodeFilter(ctx))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v.Status != node_version.VersionStatusDraft {
		return nil, fault.New("version is not a draft",
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("not a draft", "Only draft versions can be applied."),
		)
	}

	partial := node_mutate.Partial{}

	partial.Name = opt.New(v.Name)

	slug, err := mark.NewSlug(v.Slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	partial.Slug = opt.New(*slug)

	if desc, ok := v.Description.Get(); ok {
		partial.Description = opt.New(desc)
	}

	if content, ok := v.Content.Get(); ok {
		partial.Content = opt.New(content)
	}

	if props, ok := v.PropertiesSnapshot.Get(); ok {
		mutations := snapshotToPropertyMutations(props)
		partial.Properties = opt.New(mutations)
	}

	partial.Metadata = opt.New(v.Metadata)

	_, err = s.nodeMutator.UpdateFromVersion(ctx, qk, partial, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := s.versionWriter.Update(ctx, id,
		node_version_writer.WithStatus(node_version.VersionStatusApplied),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeVersionDraftApplied{
		NodeID:    updated.NodeID,
		NodeSlug:  updated.Slug,
		VersionID: updated.ID.String(),
		AuthorID:  updated.Author.ID,
	})

	return updated, nil
}
