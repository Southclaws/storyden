package plugin_manager

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (m *Manager) UpdatePackage(
	ctx context.Context,
	id plugin.InstallationID,
	reader io.Reader,
) (*plugin.Record, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if rec.Mode.External() {
		return nil, fault.Wrap(
			fault.New("cannot update external plugin package"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"invalid plugin mode",
				"Only supervised plugins support package updates. External plugins can update manifests instead.",
			),
		)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	validated, err := plugin.Binary(data).Validate(ctx)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid plugin package", "The uploaded plugin file is invalid: "+getValidationDetails(err)),
		)
	}

	if rec.Manifest.ID != validated.Metadata.ID {
		return nil, fault.Wrap(
			fault.Newf("plugin manifest id mismatch: installed=%q uploaded=%q", rec.Manifest.ID, validated.Metadata.ID),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"manifest ID mismatch",
				"The uploaded package must have the same manifest ID as the installed plugin.",
			),
		)
	}

	wasActive := rec.State == plugin.ActiveStateActive
	_, sessErr := m.runner.GetSession(ctx, id)
	hadSession := sessErr == nil

	updated, err := m.pluginWriter.UpdatePackage(ctx, id, validated)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Refresh the in-memory session so future activations run the updated
	// package even if the plugin is currently inactive.
	if hadSession || wasActive {
		if err := m.runner.Unload(ctx, id); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := m.runner.Load(ctx, *updated); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if wasActive {
			if err := m.activatePlugin(ctx, updated); err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
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
