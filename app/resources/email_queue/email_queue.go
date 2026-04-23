package email_queue

import (
	"net/mail"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	entschema "github.com/Southclaws/storyden/internal/ent/schema"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

type ID xid.ID

type Attempt struct {
	Timestamp time.Time
	Status    Status
	Error     opt.Optional[string]
}

type Email struct {
	ID               ID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RecipientAddress string
	RecipientName    string
	Subject          string
	ContentPlain     string
	ContentHTML      string
	Status           Status
	Attempts         []*Attempt
	ProcessedAt      opt.Optional[time.Time]
}

func Map(in *ent.EmailQueue) (*Email, error) {
	status, err := NewStatus(in.Status.String())
	if err != nil {
		return nil, fault.Wrap(err)
	}

	attempts, err := dt.MapErr(in.Attempts, mapAttempt)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Email{
		ID:               ID(in.ID),
		CreatedAt:        in.CreatedAt,
		UpdatedAt:        in.UpdatedAt,
		RecipientAddress: in.RecipientAddress,
		RecipientName:    in.RecipientName,
		Subject:          in.Subject,
		ContentPlain:     in.ContentPlain,
		ContentHTML:      in.ContentHTML,
		Status:           status,
		Attempts:         attempts,
		ProcessedAt:      opt.NewPtr(in.ProcessedAt),
	}, nil
}

func mapAttempt(in entschema.EmailAttempt) (*Attempt, error) {
	status, err := NewStatus(in.Status)
	if err != nil {
		return nil, err
	}

	return &Attempt{
		Timestamp: in.Timestamp,
		Status:    status,
		Error:     opt.NewPtr(in.Error),
	}, nil
}

func (i ID) String() string {
	return xid.ID(i).String()
}

func (e *Email) Message() (*mailer.Message, error) {
	return mailer.NewMessage(
		mail.Address{
			Name:    e.RecipientName,
			Address: e.RecipientAddress,
		},
		e.RecipientName,
		e.Subject,
		mailer.Content{
			HTML:  e.ContentHTML,
			Plain: e.ContentPlain,
		},
	)
}
