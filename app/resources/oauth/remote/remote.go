package remote

import (
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	ent_connection "github.com/Southclaws/storyden/internal/ent/oauthremoteconnection"
	"github.com/rs/xid"
)

type ConnectionID xid.ID
type FlowID xid.ID

func (id ConnectionID) String() string { return xid.ID(id).String() }
func (id FlowID) String() string       { return xid.ID(id).String() }

func NewConnectionID(s string) (ConnectionID, error) {
	id, err := xid.FromString(s)
	return ConnectionID(id), err
}

type Mode string
type Status string

const (
	ModeCIMD   Mode = "cimd"
	ModeDCR    Mode = "dcr"
	ModeManual Mode = "manual"

	StatusPending   Status = "pending"
	StatusConnected Status = "connected"
	StatusError     Status = "error"
)

type Connection struct {
	ID        ConnectionID
	CreatedAt time.Time
	UpdatedAt time.Time

	ResourceURL                 string
	Resource                    string
	ResourceName                string
	ProtectedResourceMetadata   map[string]any
	AuthorizationServer         string
	AuthorizationServerMetadata map[string]any
	Mode                        Mode
	Status                      Status

	ClientID                string
	ClientSecret            string
	HasClientSecret         bool
	AuthorizationEndpoint   string
	TokenEndpoint           string
	RegistrationEndpoint    string
	TokenEndpointAuthMethod string
	RedirectURIs            []string
	RedirectURI             string
	Scope                   string

	AccessToken           string
	RefreshToken          string
	HasAccessToken        bool
	HasRefreshToken       bool
	TokenType             string
	TokenExpiry           *time.Time
	TokenRefreshStartedAt *time.Time
	LastError             *string
	AddedBy               account.AccountID
}

type Flow struct {
	ID           FlowID
	ConnectionID ConnectionID
	StateHash    string
	PKCEVerifier string
	RedirectURI  string
	ExpiresAt    time.Time
	ConsumedAt   *time.Time
	Connection   *Connection
}

type ConnectionCreate struct {
	ResourceURL                 string
	Resource                    string
	ResourceName                string
	ProtectedResourceMetadata   map[string]any
	AuthorizationServer         string
	AuthorizationServerMetadata map[string]any
	Mode                        Mode
	Status                      Status
	ClientID                    string
	ClientSecret                string
	AuthorizationEndpoint       string
	TokenEndpoint               string
	RegistrationEndpoint        string
	TokenEndpointAuthMethod     string
	RedirectURIs                []string
	RedirectURI                 string
	Scope                       string
	AddedBy                     account.AccountID
}

type TokenUpdate struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	TokenExpiry  *time.Time
	Scope        string
}

func MapConnection(in *ent.OAuthRemoteConnection) Connection {
	out := Connection{
		ID:                          ConnectionID(in.ID),
		CreatedAt:                   in.CreatedAt,
		UpdatedAt:                   in.UpdatedAt,
		ResourceURL:                 in.ResourceURL,
		Resource:                    in.Resource,
		ResourceName:                in.ResourceName,
		ProtectedResourceMetadata:   cloneMap(in.ProtectedResourceMetadata),
		AuthorizationServer:         in.AuthorizationServer,
		AuthorizationServerMetadata: cloneMap(in.AuthorizationServerMetadata),
		Mode:                        Mode(in.Mode.String()),
		Status:                      Status(in.Status.String()),
		ClientID:                    in.ClientID,
		ClientSecret:                in.ClientSecret,
		HasClientSecret:             in.ClientSecret != "",
		AuthorizationEndpoint:       in.AuthorizationEndpoint,
		TokenEndpoint:               in.TokenEndpoint,
		RegistrationEndpoint:        in.RegistrationEndpoint,
		TokenEndpointAuthMethod:     in.TokenEndpointAuthMethod,
		RedirectURIs:                append([]string(nil), in.RedirectUris...),
		RedirectURI:                 in.RedirectURI,
		Scope:                       in.Scope,
		AccessToken:                 in.AccessToken,
		RefreshToken:                in.RefreshToken,
		HasAccessToken:              in.AccessToken != "",
		HasRefreshToken:             in.RefreshToken != "",
		TokenType:                   in.TokenType,
		TokenExpiry:                 in.TokenExpiry,
		TokenRefreshStartedAt:       in.TokenRefreshStartedAt,
		LastError:                   in.LastError,
		AddedBy:                     account.AccountID(in.AddedBy),
	}
	return out
}

func MapFlow(in *ent.OAuthRemoteAuthorisationFlow) Flow {
	out := Flow{
		ID:           FlowID(in.ID),
		ConnectionID: ConnectionID(in.ConnectionID),
		StateHash:    in.StateHash,
		PKCEVerifier: in.PkceVerifier,
		RedirectURI:  in.RedirectURI,
		ExpiresAt:    in.ExpiresAt,
		ConsumedAt:   in.ConsumedAt,
	}
	if in.Edges.Connection != nil {
		conn := MapConnection(in.Edges.Connection)
		out.Connection = &conn
	}
	return out
}

func toEntMode(in Mode) ent_connection.Mode {
	switch in {
	case ModeCIMD:
		return ent_connection.ModeCimd
	case ModeDCR:
		return ent_connection.ModeDcr
	default:
		return ent_connection.ModeManual
	}
}

func toEntStatus(in Status) ent_connection.Status {
	switch in {
	case StatusConnected:
		return ent_connection.StatusConnected
	case StatusError:
		return ent_connection.StatusError
	default:
		return ent_connection.StatusPending
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
