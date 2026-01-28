package plugin_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/plugin"
)

// Handle state transitions for plugins
// Activate and deactivate via runner sessions.
func (m *Manager) SetActiveState(ctx context.Context, id plugin.InstallationID, desiredState plugin.ActiveState) error {
	// Get the plugin record
	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	switch desiredState {
	case plugin.ActiveStateActive:

		if err := m.pluginWriter.SetActiveState(ctx, rec.InstallationID, plugin.ActiveStateActive); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		if err := m.activatePlugin(ctx, rec); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

	case plugin.ActiveStateInactive:
		if err := m.pluginWriter.SetActiveState(ctx, rec.InstallationID, plugin.ActiveStateInactive); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		if err := m.deactivatePlugin(ctx, rec); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

	default:
		return fault.New("unknown desired state")
	}

	return nil
}

func (m *Manager) activatePlugin(ctx context.Context, rec *plugin.Record) error {
	sess, err := m.runner.GetSession(ctx, rec.InstallationID)
	if err != nil {
		bin, err := m.pluginQuerier.LoadBinary(ctx, rec.InstallationID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		sess, err = m.runner.Load(ctx, rec.InstallationID, bin)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := sess.Start(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) deactivatePlugin(ctx context.Context, rec *plugin.Record) error {
	sess, err := m.runner.GetSession(ctx, rec.InstallationID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := sess.Stop(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
