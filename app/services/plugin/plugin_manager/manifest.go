package plugin_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (m *Manager) UpdateManifest(
	ctx context.Context,
	id plugin.InstallationID,
	manifest rpc.Manifest,
) (*plugin.Record, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if rec.Mode.Supervised() {
		return nil, fault.Wrap(
			fault.New("cannot update supervised plugin manifest"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid plugin mode", "Only external plugin manifests can be updated. To update a supervised plugin, upload a new version of the plugin (ZIP or SDX file) archive."),
		)
	}

	updated, err := m.pluginWriter.UpdateManifest(ctx, id, manifest)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if _, err := m.runner.GetSession(ctx, id); err == nil {
		if err := m.runner.Unload(ctx, id); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if updated.State == plugin.ActiveStateActive {
		if err := m.runner.Load(ctx, *updated); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sess, err := m.runner.GetSession(ctx, id)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := sess.SetActiveState(ctx, plugin.ActiveStateActive); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	out, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if sess, err := m.runner.GetSession(ctx, id); err == nil && sess != nil {
		hydrateSession(out, sess)
	}

	return out, nil
}
