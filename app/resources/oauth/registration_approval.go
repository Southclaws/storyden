package oauth

import (
	"encoding/json"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type RegistrationApprovalID xid.ID

func (i RegistrationApprovalID) XID() xid.ID { return xid.ID(i) }

type RegistrationApproval struct {
	ID                  RegistrationApprovalID
	CreatedAt           time.Time
	VerificationCode    string
	Metadata            json.RawMessage
	ExpiresAt           time.Time
	NextPollAt          time.Time
	PollIntervalSeconds int
	ApprovedByAccountID opt.Optional[account.AccountID]
	ApprovedAt          opt.Optional[time.Time]
	DeniedAt            opt.Optional[time.Time]
	CancelledAt         opt.Optional[time.Time]
	ConsumedAt          opt.Optional[time.Time]
}

func MapRegistrationApproval(in *ent.OAuthRegistrationApproval) *RegistrationApproval {
	return &RegistrationApproval{
		ID:                  RegistrationApprovalID(in.ID),
		CreatedAt:           in.CreatedAt,
		VerificationCode:    in.VerificationCodeDisplay,
		Metadata:            in.Metadata,
		ExpiresAt:           in.ExpiresAt,
		NextPollAt:          in.NextPollAt,
		PollIntervalSeconds: in.PollIntervalSeconds,
		ApprovedByAccountID: opt.NewPtrMap(in.ApprovedByAccountID, func(id xid.ID) account.AccountID { return account.AccountID(id) }),
		ApprovedAt:          opt.NewPtr(in.ApprovedAt),
		DeniedAt:            opt.NewPtr(in.DeniedAt),
		CancelledAt:         opt.NewPtr(in.CancelledAt),
		ConsumedAt:          opt.NewPtr(in.ConsumedAt),
	}
}
