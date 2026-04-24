package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type EmailQueue struct {
	ent.Schema
}

type EmailAttempt struct {
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Error     *string   `json:"error,omitempty"`
}

func (EmailQueue) Mixin() []ent.Mixin {
	return []ent.Mixin{Identifier{}, CreatedAt{}, UpdatedAt{}}
}

func (EmailQueue) Fields() []ent.Field {
	return []ent.Field{
		field.String("recipient_address").
			Comment("The destination email address."),

		field.String("recipient_name").
			Comment("The recipient display name used in the email headers."),

		field.String("subject").
			Comment("The email subject line."),

		field.String("content_plain").
			Default("").
			Comment("The plain text email body."),

		field.String("content_html").
			Default("").
			Comment("The HTML email body."),

		field.Enum("status").
			Values("pending", "processing", "sent", "failed").
			Default("pending"),

		field.JSON("attempts", []EmailAttempt{}).
			Default([]EmailAttempt{}).
			Comment("Attempt records, appended for each delivery attempt."),

		field.Time("processed_at").
			Optional().
			Nillable(),

		field.Time("available_at").
			Default(time.Now),
	}
}

func (EmailQueue) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "available_at"),
	}
}
