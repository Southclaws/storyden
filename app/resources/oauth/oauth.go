package oauth

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type (
	ClientID               xid.ID
	AuthorisationCodeID    xid.ID
	AuthorisationRequestID xid.ID
	DeviceAuthorisationID  xid.ID
	RefreshTokenID         xid.ID
)

const (
	OAuthAccessKeyKind    = "sdoak"
	OAuthAccessSecretKind = "sdoas"

	OAuthAccessKeyPrefix    = OAuthAccessKeyKind + "_"
	OAuthAccessSecretPrefix = OAuthAccessSecretKind + "_"
)

//go:generate go run github.com/Southclaws/enumerator

type clientTypeEnum string

type scopePolicyEnum string

const (
	clientTypePublic       clientTypeEnum = "public"
	clientTypeConfidential clientTypeEnum = "confidential"

	scopePolicyExplicit               scopePolicyEnum = "explicit"
	scopePolicyInheritUserPermissions scopePolicyEnum = "inherit"
)

type Client struct {
	ID                      ClientID
	CreatedAt               time.Time
	UpdatedAt               time.Time
	AccountID               opt.Optional[account.AccountID]
	ClientID                string
	ClientSecretHash        opt.Optional[string]
	Name                    string
	Type                    ClientType
	ScopePolicy             ScopePolicy
	TokenEndpointAuthMethod string
	RedirectURIs            []string
	AllowedScopes           []string
	AllowedGrants           []string
}

type AuthorisationCode struct {
	ID                  AuthorisationCodeID
	CreatedAt           time.Time
	ClientID            ClientID
	AccountID           account.AccountID
	CodeHash            string
	RedirectURI         string
	Scope               string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
	ConsumedAt          opt.Optional[time.Time]
}

type AuthorisationRequest struct {
	ID                  AuthorisationRequestID
	CreatedAt           time.Time
	ClientID            ClientID
	AccountID           account.AccountID
	RequestIDHash       string
	RedirectURI         string
	Scope               string
	State               opt.Optional[string]
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
	ApprovedAt          opt.Optional[time.Time]
	DeniedAt            opt.Optional[time.Time]
}

type DeviceAuthorisation struct {
	ID                  DeviceAuthorisationID
	CreatedAt           time.Time
	ClientID            ClientID
	DeviceCodeHash      string
	UserCodeHash        string
	UserCodeDisplay     string
	Scope               string
	ExpiresAt           time.Time
	PollIntervalSeconds int
	LastPolledAt        opt.Optional[time.Time]
	ClaimedByAccountID  opt.Optional[account.AccountID]
	ApprovedByAccountID opt.Optional[account.AccountID]
	ApprovedAt          opt.Optional[time.Time]
	DeniedAt            opt.Optional[time.Time]
	ConsumedAt          opt.Optional[time.Time]
}

type RefreshToken struct {
	ID                RefreshTokenID
	CreatedAt         time.Time
	ClientID          ClientID
	ClientIdentifier  string
	ClientName        string
	AccountID         account.AccountID
	TokenHash         string
	Scope             string
	ExpiresAt         time.Time
	RevokedAt         opt.Optional[time.Time]
	ReplacedByTokenID opt.Optional[RefreshTokenID]
	LastUsedAt        opt.Optional[time.Time]
}

func MapClient(in *ent.OAuthClient) *Client {
	clientType, _ := NewClientType(in.Type.String())
	scopePolicy, _ := NewScopePolicy(in.ScopePolicy.String())

	return &Client{
		ID:                      ClientID(in.ID),
		CreatedAt:               in.CreatedAt,
		UpdatedAt:               in.UpdatedAt,
		AccountID:               opt.NewPtrMap(in.AccountID, func(id xid.ID) account.AccountID { return account.AccountID(id) }),
		ClientID:                in.ClientID,
		ClientSecretHash:        opt.NewPtr(in.ClientSecretHash),
		Name:                    in.Name,
		Type:                    clientType,
		ScopePolicy:             scopePolicy,
		TokenEndpointAuthMethod: in.TokenEndpointAuthMethod,
		RedirectURIs:            in.RedirectUris,
		AllowedScopes:           in.AllowedScopes,
		AllowedGrants:           in.AllowedGrants,
	}
}

func MapAuthorisationCode(in *ent.OAuthAuthorisationCode) *AuthorisationCode {
	return &AuthorisationCode{
		ID:                  AuthorisationCodeID(in.ID),
		CreatedAt:           in.CreatedAt,
		ClientID:            ClientID(in.ClientID),
		AccountID:           account.AccountID(in.AccountID),
		CodeHash:            in.CodeHash,
		RedirectURI:         in.RedirectURI,
		Scope:               in.Scope,
		CodeChallenge:       in.CodeChallenge,
		CodeChallengeMethod: in.CodeChallengeMethod.String(),
		ExpiresAt:           in.ExpiresAt,
		ConsumedAt:          opt.NewPtr(in.ConsumedAt),
	}
}

func MapAuthorisationRequest(in *ent.OAuthAuthorisationRequest) *AuthorisationRequest {
	return &AuthorisationRequest{
		ID:                  AuthorisationRequestID(in.ID),
		CreatedAt:           in.CreatedAt,
		ClientID:            ClientID(in.ClientID),
		AccountID:           account.AccountID(in.AccountID),
		RequestIDHash:       in.RequestIDHash,
		RedirectURI:         in.RedirectURI,
		Scope:               in.Scope,
		State:               opt.NewPtr(in.State),
		CodeChallenge:       in.CodeChallenge,
		CodeChallengeMethod: in.CodeChallengeMethod.String(),
		ExpiresAt:           in.ExpiresAt,
		ApprovedAt:          opt.NewPtr(in.ApprovedAt),
		DeniedAt:            opt.NewPtr(in.DeniedAt),
	}
}

func MapDeviceAuthorisation(in *ent.OAuthDeviceAuthorisation) *DeviceAuthorisation {
	return &DeviceAuthorisation{
		ID:                  DeviceAuthorisationID(in.ID),
		CreatedAt:           in.CreatedAt,
		ClientID:            ClientID(in.ClientID),
		DeviceCodeHash:      in.DeviceCodeHash,
		UserCodeHash:        in.UserCodeHash,
		UserCodeDisplay:     in.UserCodeDisplay,
		Scope:               in.Scope,
		ExpiresAt:           in.ExpiresAt,
		PollIntervalSeconds: in.PollIntervalSeconds,
		LastPolledAt:        opt.NewPtr(in.LastPolledAt),
		ClaimedByAccountID:  opt.NewPtrMap(in.ClaimedByAccountID, func(id xid.ID) account.AccountID { return account.AccountID(id) }),
		ApprovedByAccountID: opt.NewPtrMap(in.ApprovedByAccountID, func(id xid.ID) account.AccountID { return account.AccountID(id) }),
		ApprovedAt:          opt.NewPtr(in.ApprovedAt),
		DeniedAt:            opt.NewPtr(in.DeniedAt),
		ConsumedAt:          opt.NewPtr(in.ConsumedAt),
	}
}

func MapRefreshToken(in *ent.OAuthRefreshToken) *RefreshToken {
	out := &RefreshToken{
		ID:                RefreshTokenID(in.ID),
		CreatedAt:         in.CreatedAt,
		ClientID:          ClientID(in.ClientID),
		AccountID:         account.AccountID(in.AccountID),
		TokenHash:         in.TokenHash,
		Scope:             in.Scope,
		ExpiresAt:         in.ExpiresAt,
		RevokedAt:         opt.NewPtr(in.RevokedAt),
		ReplacedByTokenID: opt.NewPtrMap(in.ReplacedByTokenID, func(id xid.ID) RefreshTokenID { return RefreshTokenID(id) }),
		LastUsedAt:        opt.NewPtr(in.LastUsedAt),
	}

	if in.Edges.Client != nil {
		out.ClientIdentifier = in.Edges.Client.ClientID
		out.ClientName = in.Edges.Client.Name
	}

	return out
}

func (i ClientID) XID() xid.ID {
	return xid.ID(i)
}

func (i AuthorisationCodeID) XID() xid.ID {
	return xid.ID(i)
}

func (i AuthorisationRequestID) XID() xid.ID {
	return xid.ID(i)
}

func (i DeviceAuthorisationID) XID() xid.ID {
	return xid.ID(i)
}

func (i RefreshTokenID) XID() xid.ID {
	return xid.ID(i)
}
