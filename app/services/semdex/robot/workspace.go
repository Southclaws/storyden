package robot

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_workspace"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
)

const WorkspaceStateKey = workspacestate.WorkspaceStateKey

type WorkspaceMountSpec = agent_registry.WorkspaceMountSpec

type WorkspaceManager struct {
	logger      *slog.Logger
	repo        *robot_workspace.Repository
	sessionRepo *robot_session.Repository
	providers   *workspaceprovider.Registry
}

func NewWorkspaceManager(
	logger *slog.Logger,
	repo *robot_workspace.Repository,
	sessionRepo *robot_session.Repository,
	providers *workspaceprovider.Registry,
) *WorkspaceManager {
	return &WorkspaceManager{
		logger:      logger,
		repo:        repo,
		sessionRepo: sessionRepo,
		providers:   providers,
	}
}

func (m *WorkspaceManager) Mount(
	ctx context.Context,
	sessionID robotresource.SessionID,
	accountID account.AccountID,
	spec WorkspaceMountSpec,
) (*robotresource.WorkspaceMount, error) {
	_, hasWorkspaceID := spec.WorkspaceID.Get()
	_, hasInstanceID := spec.WorkspaceInstanceID.Get()
	if hasWorkspaceID == hasInstanceID {
		return nil, fault.Wrap(fault.New("provide exactly one workspace_id or workspace_instance_id"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var instance *robotresource.WorkspaceInstance
	var workspace *robotresource.Workspace
	var err error
	if workspaceID, ok := spec.WorkspaceID.Get(); ok {
		workspace, err = m.repo.Get(ctx, workspaceID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		instance, err = m.repo.CreateInstance(ctx, workspaceID, accountID, map[string]any{}, spec.Metadata)
	} else {
		instanceID, _ := spec.WorkspaceInstanceID.Get()
		instance, err = m.repo.GetInstance(ctx, instanceID)
		if err == nil {
			workspace, err = m.repo.Get(ctx, instance.WorkspaceID)
		}
	}
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	provider, err := m.providers.Get(ctx, instance.Provider)
	if err != nil {
		return nil, err
	}
	providerState, err := provider.Mount(ctx, instance)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !mapsEqual(instance.ProviderState, providerState) {
		instance, err = m.repo.UpdateInstanceProviderState(ctx, instance.ID, providerState)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	mount := &robotresource.WorkspaceMount{
		WorkspaceID:            instance.WorkspaceID,
		WorkspaceInstanceID:    instance.ID,
		Provider:               instance.Provider,
		ProviderState:          providerState,
		AllowUntrustedCommands: workspace.AllowUntrustedCommands,
		Metadata:               instance.Metadata,
	}

	if err := m.storeMount(ctx, sessionID, mount); err != nil {
		return nil, err
	}

	return mount, nil
}

func (m *WorkspaceManager) DeleteInstance(ctx context.Context, id robotresource.WorkspaceInstanceID) error {
	instance, err := m.repo.GetInstance(ctx, id)
	if err != nil {
		return err
	}

	if err := m.repo.DeleteInstance(ctx, id); err != nil {
		return err
	}

	provider, err := m.providers.Get(ctx, instance.Provider)
	if err != nil {
		m.logger.Warn("workspace provider not registered for instance cleanup",
			slog.String("workspace_instance_id", id.String()),
			slog.String("provider", string(instance.Provider)),
			slog.String("error", err.Error()))

		return nil
	}
	if err := provider.Cleanup(ctx, instance); err != nil {
		m.logger.Warn("failed to clean up workspace instance",
			slog.String("workspace_instance_id", id.String()),
			slog.String("provider", string(instance.Provider)),
			slog.String("error", err.Error()))
	}

	return nil
}

func (m *WorkspaceManager) storeMount(ctx context.Context, sessionID robotresource.SessionID, mount *robotresource.WorkspaceMount) error {
	sess, _, err := m.sessionRepo.Get(ctx, sessionID, robotresource.NewMessageCursorParams(opt.NewEmpty[robotresource.MessageID](), 1))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	state := sess.State
	if state == nil {
		state = make(map[string]any)
	}
	state[WorkspaceStateKey] = workspacestate.MountToState(mount)

	return m.sessionRepo.UpdateState(ctx, sessionID, state)
}

func WorkspaceMountFromState(state map[string]any) opt.Optional[robotresource.WorkspaceMount] {
	return workspacestate.MountFromState(state)
}

func mapsEqual(a, b map[string]any) bool {
	aj, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aj) == string(bj)
}
