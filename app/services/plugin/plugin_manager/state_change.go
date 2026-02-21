package plugin_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
)

const stateChangeTimeout = 5 * time.Minute

// Handle state transitions for plugins
// Activate and deactivate via runner sessions.
func (m *Manager) SetActiveState(ctx context.Context, id plugin.InstallationID, desiredState plugin.ActiveState) error {
	// Get the plugin record
	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if rec.Mode.External() {
		return fault.Wrap(
			fault.New("cannot set active state for external plugins"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
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
		// Session doesn't exist, load it first
		if err := m.runner.Load(ctx, *rec); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		sess, err = m.runner.GetSession(ctx, rec.InstallationID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	// Set active state - for supervised plugins, this starts the process
	if err := sess.SetActiveState(ctx, plugin.ActiveStateActive); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if sess.Supervised() != nil {
		if err := m.waitForState(ctx, sess, plugin.ReportedStateActive, stateChangeTimeout); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (m *Manager) deactivatePlugin(ctx context.Context, rec *plugin.Record) error {
	sess, err := m.runner.GetSession(ctx, rec.InstallationID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Set inactive state - for supervised plugins, this stops the process
	if err := sess.SetActiveState(ctx, plugin.ActiveStateInactive); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.waitForState(ctx, sess, plugin.ReportedStateInactive, stateChangeTimeout); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) waitForState(ctx context.Context, sess plugin_runner.Session, desiredState plugin.ReportedState, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			currentState := sess.GetReportedState()
			errorMessage := sess.GetErrorMessage()

			if currentState == plugin.ReportedStateError && errorMessage != "" {
				return fault.Wrap(
					fault.New(errorMessage),
					ftag.With(ftag.Cancelled),
					fmsg.With("plugin reported an error while changing state"),
					fmsg.WithDesc(
						"plugin failed while changing state",
						"The plugin reported an error while changing state. Check plugin logs for details.",
					),
				)
			}

			err := fault.Newf(
				"plugin did not reach %q before timeout (current state: %q, message: %q)",
				desiredState.String(),
				currentState.String(),
				errorMessage,
			)

			return fault.Wrap(
				err,
				ftag.With(ftag.InvalidArgument),
				fmsg.With("timeout waiting for state change"),
				fmsg.WithDesc(
					"plugin did not become ready in time",
					fmt.Sprintf(
						"The plugin did not reach %q before timing out. Current state is %q. Check plugin logs for more details.",
						desiredState.String(),
						currentState.String(),
					),
				),
			)

		case <-ticker.C:
			currentState := sess.GetReportedState()
			if currentState == desiredState {
				return nil
			}

			if currentState == plugin.ReportedStateError {
				message := sess.GetErrorMessage()
				if message == "" {
					message = "plugin entered error state"
				}

				return fault.Wrap(
					fault.New(message),
					ftag.With(ftag.Cancelled),
					fmsg.With("plugin reported an error while changing state"),
					fmsg.WithDesc(
						"plugin failed while changing state",
						"The plugin reported an error while changing state. Check plugin logs for details.",
					),
				)
			}
		}
	}
}
