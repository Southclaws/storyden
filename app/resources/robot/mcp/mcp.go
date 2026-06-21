package mcp

import (
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

type ServerID xid.ID

func (id ServerID) String() string {
	return xid.ID(id).String()
}

func NewServerID(s string) (ServerID, error) {
	id, err := xid.FromString(s)
	if err != nil {
		return ServerID{}, err
	}
	return ServerID(id), nil
}

type Server struct {
	ID        ServerID
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Slug        string
	Description string
	EndpointURL string
	Enabled     bool

	OAuthRemoteConnectionID *oauth_remote.ConnectionID
	OAuthAccessToken        string
	HasOAuthAccessToken     bool
	BearerToken             string
	HasBearerToken          bool

	LastRefreshedAt *time.Time
	LastError       *string

	AddedBy account.AccountID
	Tools   []Tool
}

type Tool struct {
	ID           string
	RemoteName   string
	CallableName string
	Title        string
	Description  string
	InputSchema  map[string]any
	OutputSchema map[string]any
	Annotations  map[string]any
	Enabled      bool
	LastSeenAt   time.Time

	ServerID   ServerID
	ServerSlug string
}

type ServerCreate struct {
	Name                    string
	Slug                    string
	Description             string
	EndpointURL             string
	Enabled                 bool
	BearerToken             string
	OAuthRemoteConnectionID *oauth_remote.ConnectionID
	AddedBy                 account.AccountID
}

type ServerUpdate struct {
	Name             *string
	Description      *string
	EndpointURL      *string
	Enabled          *bool
	BearerToken      *string
	ClearBearerToken bool
}

func MapServer(in *ent.RobotMCPServer) Server {
	out := Server{
		ID:              ServerID(in.ID),
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
		Name:            in.Name,
		Slug:            in.Slug,
		Description:     in.Description,
		EndpointURL:     in.EndpointURL,
		Enabled:         in.Enabled,
		BearerToken:     in.BearerToken,
		HasBearerToken:  in.BearerToken != "",
		LastRefreshedAt: in.LastRefreshedAt,
		LastError:       in.LastError,
		AddedBy:         account.AccountID(in.AddedBy),
	}
	if in.OauthRemoteConnectionID != nil {
		id := oauth_remote.ConnectionID(*in.OauthRemoteConnectionID)
		out.OAuthRemoteConnectionID = &id
	}
	if in.Edges.OauthRemoteConnection != nil {
		out.OAuthAccessToken = in.Edges.OauthRemoteConnection.AccessToken
		out.HasOAuthAccessToken = in.Edges.OauthRemoteConnection.AccessToken != ""
	}

	if in.Edges.Tools != nil {
		out.Tools = make([]Tool, 0, len(in.Edges.Tools))
		for _, tool := range in.Edges.Tools {
			out.Tools = append(out.Tools, MapTool(tool, in.Slug))
		}
	}

	return out
}

func MapTool(in *ent.RobotMCPTool, serverSlug string) Tool {
	return Tool{
		ID:           in.ToolID,
		RemoteName:   in.RemoteName,
		CallableName: in.CallableName,
		Title:        in.Title,
		Description:  in.Description,
		InputSchema:  cloneMap(in.InputSchema),
		OutputSchema: cloneMap(in.OutputSchema),
		Annotations:  cloneMap(in.Annotations),
		Enabled:      in.Enabled,
		LastSeenAt:   in.LastSeenAt,
		ServerID:     ServerID(in.ServerID),
		ServerSlug:   serverSlug,
	}
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
