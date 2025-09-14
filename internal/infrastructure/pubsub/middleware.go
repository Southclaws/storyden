package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

const (
	accountIDKey      = "storyden-account-id"
	rolesKey          = "storyden-roles"
	securitySchemeKey = "storyden-security-scheme"
)

// propagates session context to message subscribers.
type sessionContextMiddleware struct {
	logger *slog.Logger
}

func newSessionContextMiddleware(logger *slog.Logger) message.HandlerMiddleware {
	m := &sessionContextMiddleware{logger: logger}
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			ctx, err := m.extractSessionContext(msg.Context(), msg)
			if err != nil {
				m.logger.Error("failed to extract session context", slog.String("error", err.Error()))
				ctx = session.WithInternal(msg.Context())
			}
			msg.SetContext(ctx)
			return h(msg)
		}
	}
}

func publisherContextMiddleware(pub message.Publisher) message.Publisher {
	return &sessionContextPublisher{
		publisher: pub,
	}
}

type sessionContextPublisher struct {
	publisher message.Publisher
}

func (p *sessionContextPublisher) Publish(topic string, messages ...*message.Message) error {
	for _, msg := range messages {
		injectSessionContext(msg.Context(), msg)
	}
	return p.publisher.Publish(topic, messages...)
}

func (p *sessionContextPublisher) Close() error {
	return p.publisher.Close()
}

func (m *sessionContextMiddleware) extractSessionContext(ctx context.Context, msg *message.Message) (context.Context, error) {
	accountIDStr := msg.Metadata.Get(accountIDKey)
	rolesStr := msg.Metadata.Get(rolesKey)
	securityScheme := msg.Metadata.Get(securitySchemeKey)

	if accountIDStr == "" && rolesStr == "" {
		return session.WithInternal(ctx), nil
	}

	var roles role.Roles
	if rolesStr != "" {
		if err := json.Unmarshal([]byte(rolesStr), &roles); err != nil {
			return nil, err
		}
	}

	if accountIDStr == "" {
		return session.WithGuest(ctx, roles), nil
	}

	xidID, err := xid.FromString(accountIDStr)
	if err != nil {
		return nil, err
	}
	accountID := account.AccountID(xidID)

	acc := account.Account{
		ID: accountID,
	}

	if securityScheme == "access_key" {
		return session.WithAccessKey(ctx, acc, roles), nil
	}

	return session.WithAccount(ctx, acc, roles), nil
}

func injectSessionContext(ctx context.Context, msg *message.Message) {
	optAccountID := session.GetOptAccountID(ctx)

	accountID, ok := optAccountID.Get()
	if !ok {
		return
	}

	msg.Metadata.Set(accountIDKey, accountID.String())

	roles := session.GetRoles(ctx)
	if rolesJSON, err := json.Marshal(roles); err == nil {
		msg.Metadata.Set(rolesKey, string(rolesJSON))
	}

	if scheme, err := session.GetSecurityScheme(ctx); err == nil {
		msg.Metadata.Set(securitySchemeKey, scheme)
	}
}
