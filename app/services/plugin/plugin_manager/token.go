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
)

func (m *Manager) CycleExternalToken(ctx context.Context, id plugin.InstallationID) (string, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	if rec.Mode.Supervised() {
		return "", fault.Wrap(
			fault.New("cannot cycle supervised plugin token"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid plugin mode", "Only external plugin tokens can be cycled manually."),
		)
	}

	token, err := m.pluginWriter.CycleExternalToken(ctx, id)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	// Invalidate any currently connected external session so the plugin must
	// reconnect with the newly issued token.
	if sess, err := m.runner.GetSession(ctx, id); err == nil && sess != nil {
		if err := sess.SetActiveState(ctx, plugin.ActiveStateInactive); err == nil {
			_ = sess.SetActiveState(ctx, plugin.ActiveStateActive)
		}
	}

	return token, nil
}
