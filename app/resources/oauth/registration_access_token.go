package oauth

import (
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type DynamicRegistrationAccessTokenID xid.ID

func (i DynamicRegistrationAccessTokenID) XID() xid.ID {
	return xid.ID(i)
}

const DynamicRegistrationAccessTokenPrefix = "sdiat_"

type DynamicRegistrationAccessToken struct {
	ID                DynamicRegistrationAccessTokenID
	CreatedAt         time.Time
	Label             string
	CreatorAccountID  account.AccountID
	ExpiresAt         time.Time
	MaxRegistrations  int
	RegistrationCount int
	RevokedAt         opt.Optional[time.Time]
}

func MapDynamicRegistrationAccessToken(in *ent.OAuthDynamicRegistrationAccessTokens) *DynamicRegistrationAccessToken {
	return &DynamicRegistrationAccessToken{
		ID:                DynamicRegistrationAccessTokenID(in.ID),
		CreatedAt:         in.CreatedAt,
		Label:             in.Label,
		CreatorAccountID:  account.AccountID(in.CreatorAccountID),
		ExpiresAt:         in.ExpiresAt,
		MaxRegistrations:  in.MaxRegistrations,
		RegistrationCount: in.RegistrationCount,
		RevokedAt:         opt.NewPtr(in.RevokedAt),
	}
}
