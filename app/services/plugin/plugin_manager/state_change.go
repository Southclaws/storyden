package plugin_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/plugin"
)

// Handle state transitions for plugins
// Activate and deactivate via runner sessions.
func (m *Manager) SetActiveState(ctx context.Context, id plugin.ID, desiredState plugin.ActiveState) error {
	// Get the plugin record
	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	switch desiredState {
	case plugin.ActiveStateActive:
		if err := m.activatePlugin(ctx, rec); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

	case plugin.ActiveStateInactive:
		if err := m.deactivatePlugin(ctx, rec); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

	default:
		return fault.New("unknown desired state")
	}

	// Update active state?
	// We need to make a decision around what the DB actually holds:
	// - desired state, and only use it internally for state reconciliation
	// - current state, and only update it when we actually change state
	// m.pluginWriter.Update

	return nil
}

func (m *Manager) activatePlugin(ctx context.Context, rec *plugin.Record) error {
	_, err := m.runner.GetSession(ctx, rec.Manifest.ID)
	if err != nil {
		bin, err := m.pluginQuerier.LoadBinary(ctx, rec.ID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err = m.runner.Load(ctx, bin)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := m.runner.StartPlugin(ctx, rec.Manifest.ID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) deactivatePlugin(ctx context.Context, rec *plugin.Record) error {
	if err := m.runner.StopPlugin(ctx, rec.Manifest.ID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
