package pluginbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	adksession "google.golang.org/adk/session"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
)

const pluginBuildTargetDifferentPluginMessage = "this chat is already working on a different plugin; start a new chat to work on another plugin"
const pluginBuilderPostInstallEditMessage = "this plugin is already installed in this chat; do not edit workspace files after a successful install unless the user requested another change. If a new change is required, set allow_after_install=true, then run plugin_validate and plugin_install again before saying done"

const (
	pluginBuildTargetStateKey     = "plugin_builder_target"
	pluginBuildTargetModeNew      = "new"
	pluginBuildTargetModeExisting = "existing"
)

type pluginBuildTarget struct {
	Mode           string `json:"mode"`
	InstallationID string `json:"installation_id,omitempty"`
	ManifestID     string `json:"manifest_id,omitempty"`
}

type pluginBuilderStateContext interface {
	State() adksession.State
}

type pluginBuilderReadonlyStateContext interface {
	ReadonlyState() adksession.ReadonlyState
}

type pluginBuilderSessionContext interface {
	SessionID() string
}

func pluginBuildTargetFromContext(ctx context.Context) (pluginBuildTarget, bool, error) {
	stateProvider, ok := ctx.(pluginBuilderStateContext)
	if !ok || stateProvider.State() == nil {
		return pluginBuildTarget{}, false, nil
	}

	if target, ok, err := pluginBuildTargetFromSessionState(stateProvider.State()); err != nil || ok {
		return target, ok, err
	}

	state := map[string]any{}
	for key, value := range stateProvider.State().All() {
		state[key] = value
	}
	mount, ok := workspacestate.MountFromState(state).Get()
	if !ok || mount.Metadata == nil {
		return pluginBuildTarget{}, false, nil
	}

	return pluginBuildTargetFromValue(mount.Metadata[pluginBuildTargetStateKey])
}

func pluginBuildTargetFromReadonlyContext(ctx context.Context) (pluginBuildTarget, bool, error) {
	stateProvider, ok := ctx.(pluginBuilderReadonlyStateContext)
	if !ok || stateProvider.ReadonlyState() == nil {
		return pluginBuildTarget{}, false, nil
	}

	if target, ok, err := pluginBuildTargetFromReadonlyState(stateProvider.ReadonlyState()); err != nil || ok {
		return target, ok, err
	}

	state := map[string]any{}
	for key, value := range stateProvider.ReadonlyState().All() {
		state[key] = value
	}
	mount, ok := workspacestate.MountFromState(state).Get()
	if !ok || mount.Metadata == nil {
		return pluginBuildTarget{}, false, nil
	}

	return pluginBuildTargetFromValue(mount.Metadata[pluginBuildTargetStateKey])
}

func pluginBuilderAllowUntrustedCommandsFromContext(ctx context.Context) bool {
	stateProvider, ok := ctx.(pluginBuilderStateContext)
	if !ok || stateProvider.State() == nil {
		return false
	}

	state := map[string]any{}
	for key, value := range stateProvider.State().All() {
		state[key] = value
	}
	mount, ok := workspacestate.MountFromState(state).Get()
	return ok && mount.AllowUntrustedCommands
}

func pluginBuilderAllowUntrustedCommandsFromReadonlyContext(ctx context.Context) bool {
	stateProvider, ok := ctx.(pluginBuilderReadonlyStateContext)
	if !ok || stateProvider.ReadonlyState() == nil {
		return false
	}

	state := map[string]any{}
	for key, value := range stateProvider.ReadonlyState().All() {
		state[key] = value
	}
	mount, ok := workspacestate.MountFromState(state).Get()
	return ok && mount.AllowUntrustedCommands
}

func pluginBuildTargetFromSessionState(state adksession.State) (pluginBuildTarget, bool, error) {
	return pluginBuildTargetFromReadonlyState(state)
}

func pluginBuildTargetFromReadonlyState(state adksession.ReadonlyState) (pluginBuildTarget, bool, error) {
	raw, err := state.Get(pluginBuildTargetStateKey)
	if err != nil {
		if errors.Is(err, adksession.ErrStateKeyNotExist) {
			return pluginBuildTarget{}, false, nil
		}
		return pluginBuildTarget{}, false, err
	}

	return pluginBuildTargetFromValue(raw)
}

func pluginBuildTargetFromValue(raw any) (pluginBuildTarget, bool, error) {
	if raw == nil {
		return pluginBuildTarget{}, false, nil
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return pluginBuildTarget{}, false, err
	}

	var target pluginBuildTarget
	if err := json.Unmarshal(data, &target); err != nil {
		return pluginBuildTarget{}, false, err
	}
	if target.Mode == "" && target.InstallationID == "" && target.ManifestID == "" {
		return pluginBuildTarget{}, false, nil
	}

	return target, true, nil
}

func (a *Agent) setPluginBuildTarget(ctx context.Context, target pluginBuildTarget) error {
	stateProvider, ok := ctx.(pluginBuilderStateContext)
	if !ok || stateProvider.State() == nil {
		return nil
	}

	if target.Mode == "" {
		target.Mode = pluginBuildTargetModeNew
	}
	if target.InstallationID == "" && target.Mode != pluginBuildTargetModeNew {
		return fmt.Errorf("plugin build target installation_id is required")
	}

	if err := stateProvider.State().Set(pluginBuildTargetStateKey, target); err != nil {
		return err
	}

	if a == nil || a.sessions == nil {
		return nil
	}
	sessionProvider, ok := ctx.(pluginBuilderSessionContext)
	if !ok || strings.TrimSpace(sessionProvider.SessionID()) == "" {
		return nil
	}
	sessionID, err := robotresource.NewSessionID(sessionProvider.SessionID())
	if err != nil {
		return err
	}

	state := map[string]any{}
	for key, value := range stateProvider.State().All() {
		state[key] = value
	}
	state[pluginBuildTargetStateKey] = target

	return a.sessions.UpdateState(ctx, sessionID, state)
}

func ensurePluginBuildTarget(ctx context.Context, manifestID string, installationID string) error {
	target, ok, err := pluginBuildTargetFromContext(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if target.ManifestID != "" && target.ManifestID != manifestID {
		return errors.New(pluginBuildTargetDifferentPluginMessage)
	}
	if target.InstallationID != "" && target.InstallationID != installationID {
		return errors.New(pluginBuildTargetDifferentPluginMessage)
	}
	return nil
}

func requirePostInstallEditIntent(ctx context.Context, allowed bool) (bool, error) {
	target, ok, err := pluginBuildTargetFromContext(ctx)
	if err != nil {
		return false, err
	}
	if !ok || strings.TrimSpace(target.InstallationID) == "" {
		return false, nil
	}
	if !allowed {
		return true, errors.New(pluginBuilderPostInstallEditMessage)
	}
	return true, nil
}

func workspaceEditNextAction(afterInstall bool) string {
	if afterInstall {
		return "Workspace changed after install; the installed plugin package is now stale. Run plugin_validate and plugin_install again before saying done."
	}
	return "Inspect or validate the affected behavior before installing."
}
