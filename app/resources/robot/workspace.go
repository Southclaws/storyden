package robot

import (
	"time"

	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type WorkspaceID xid.ID

func (id WorkspaceID) String() string {
	return xid.ID(id).String()
}

func NewWorkspaceID(s string) (WorkspaceID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return WorkspaceID{}, err
	}
	return WorkspaceID(id), nil
}

type WorkspaceInstanceID xid.ID

func (id WorkspaceInstanceID) String() string {
	return xid.ID(id).String()
}

func NewWorkspaceInstanceID(s string) (WorkspaceInstanceID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return WorkspaceInstanceID{}, err
	}
	return WorkspaceInstanceID(id), nil
}

type WorkspaceProvider string

const (
	WorkspaceProviderLocal   WorkspaceProvider = "local"
	WorkspaceProviderSprites WorkspaceProvider = "sprites"
)

type Workspace struct {
	ID        WorkspaceID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Provider    WorkspaceProvider
	Config      map[string]any
	Metadata    map[string]any

	Creator account.Account
}

type WorkspaceInstance struct {
	ID        WorkspaceInstanceID
	CreatedAt time.Time
	UpdatedAt time.Time

	WorkspaceID   WorkspaceID
	Provider      WorkspaceProvider
	ProviderState map[string]any
	Metadata      map[string]any

	Creator account.Account
}

type WorkspaceMount struct {
	WorkspaceID         WorkspaceID
	WorkspaceInstanceID WorkspaceInstanceID
	Provider            WorkspaceProvider
	ProviderState       map[string]any
	Metadata            map[string]any
}

func MapWorkspace(in *ent.RobotWorkspace) (*Workspace, error) {
	creator, err := account.MapRef(in.Edges.Creator)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Workspace{
		ID:          WorkspaceID(in.ID),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Description: in.Description,
		Provider:    WorkspaceProvider(in.Provider),
		Config:      in.Config,
		Metadata:    in.Metadata,
		Creator:     *creator,
	}, nil
}

func MapWorkspaceInstance(in *ent.RobotWorkspaceInstance) (*WorkspaceInstance, error) {
	creator, err := account.MapRef(in.Edges.Creator)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &WorkspaceInstance{
		ID:            WorkspaceInstanceID(in.ID),
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
		WorkspaceID:   WorkspaceID(in.WorkspaceID),
		Provider:      WorkspaceProvider(in.Provider),
		ProviderState: in.ProviderState,
		Metadata:      in.Metadata,
		Creator:       *creator,
	}, nil
}
