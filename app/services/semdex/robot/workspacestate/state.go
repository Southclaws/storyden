package workspacestate

import (
	"encoding/json"

	"github.com/Southclaws/opt"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
)

const WorkspaceStateKey = "robot_workspace"

type ReadonlyState interface {
	Get(string) (any, error)
}

func MountToState(mount *robotresource.WorkspaceMount) map[string]any {
	return map[string]any{
		"workspace_id":             mount.WorkspaceID.String(),
		"workspace_instance_id":    mount.WorkspaceInstanceID.String(),
		"provider":                 string(mount.Provider),
		"provider_state":           mount.ProviderState,
		"allow_untrusted_commands": mount.AllowUntrustedCommands,
		"metadata":                 mount.Metadata,
	}
}

func MountFromState(state map[string]any) opt.Optional[robotresource.WorkspaceMount] {
	if state == nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}

	raw, ok := state[WorkspaceStateKey]
	if !ok {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}

	var payload struct {
		WorkspaceID            string         `json:"workspace_id"`
		WorkspaceInstanceID    string         `json:"workspace_instance_id"`
		Provider               string         `json:"provider"`
		ProviderState          map[string]any `json:"provider_state"`
		AllowUntrustedCommands bool           `json:"allow_untrusted_commands"`
		Metadata               map[string]any `json:"metadata"`
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}

	workspaceID, err := robotresource.NewWorkspaceID(payload.WorkspaceID)
	if err != nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}
	instanceID, err := robotresource.NewWorkspaceInstanceID(payload.WorkspaceInstanceID)
	if err != nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}

	return opt.New(robotresource.WorkspaceMount{
		WorkspaceID:            workspaceID,
		WorkspaceInstanceID:    instanceID,
		Provider:               robotresource.WorkspaceProvider(payload.Provider),
		ProviderState:          payload.ProviderState,
		AllowUntrustedCommands: payload.AllowUntrustedCommands,
		Metadata:               payload.Metadata,
	})
}

func MountFromReadonlyState(state ReadonlyState) opt.Optional[robotresource.WorkspaceMount] {
	if state == nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}

	value, err := state.Get(WorkspaceStateKey)
	if err != nil {
		return opt.NewEmpty[robotresource.WorkspaceMount]()
	}
	return MountFromState(map[string]any{WorkspaceStateKey: value})
}

func Available(state ReadonlyState) bool {
	_, ok := MountFromReadonlyState(state).Get()
	return ok
}
