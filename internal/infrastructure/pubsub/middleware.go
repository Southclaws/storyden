package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

const (
	accountIDKey      = "storyden-account-id"
	permissionsKey    = "storyden-permissions"
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
				return nil, err
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
		// Message handlers should inherit caller identity, but not the caller's
		// cancellation or deadline. Otherwise short-lived HTTP request contexts
		// can cancel queued work before consumers process it.
		ctx := context.WithoutCancel(msg.Context())
		msg.SetContext(ctx)
		if err := injectSessionContext(ctx, msg); err != nil {
			return err
		}
	}
	return p.publisher.Publish(topic, messages...)
}

func (p *sessionContextPublisher) Close() error {
	return p.publisher.Close()
}

func (m *sessionContextMiddleware) extractSessionContext(ctx context.Context, msg *message.Message) (context.Context, error) {
	accountIDStr := msg.Metadata.Get(accountIDKey)
	permissionsStr := msg.Metadata.Get(permissionsKey)
	securityScheme := msg.Metadata.Get(securitySchemeKey)

	if accountIDStr == "" && permissionsStr == "" {
		return session.WithInternal(ctx), nil
	}

	permissions := rbac.NewList()
	if permissionsStr != "" {
		var list rbac.PermissionList
		if err := json.Unmarshal([]byte(permissionsStr), &list); err != nil {
			return nil, fault.Wrap(
				err,
				fctx.With(ctx),
				fmsg.With("invalid permissions metadata"),
				ftag.With(ftag.Unauthenticated),
			)
		}

		permissions = rbac.NewList(list...)
	}

	if accountIDStr == "" {
		return session.WithGuestPermissions(ctx, permissions), nil
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
		return session.WithAccessKeyPermissions(ctx, acc, permissions), nil
	}

	return session.WithAccountPermissions(ctx, acc, permissions), nil
}

func injectSessionContext(ctx context.Context, msg *message.Message) error {
	if !session.HasContext(ctx) {
		return nil
	}

	permissions, err := session.GetPermissions(ctx)
	if err != nil {
		return err
	}

	permissionsJSON, err := json.Marshal(permissions.List())
	if err != nil {
		return fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("failed to marshal permissions metadata"),
		)
	}
	msg.Metadata.Set(permissionsKey, string(permissionsJSON))

	optAccountID := session.GetOptAccountID(ctx)

	accountID, ok := optAccountID.Get()
	if !ok {
		return nil
	}

	msg.Metadata.Set(accountIDKey, accountID.String())

	if scheme, err := session.GetSecurityScheme(ctx); err == nil {
		msg.Metadata.Set(securitySchemeKey, scheme)
	}

	return nil
}

func newChaosDelayMiddleware(maxDelay time.Duration, logger *slog.Logger) message.HandlerMiddleware {
	if maxDelay == 0 {
		return func(h message.HandlerFunc) message.HandlerFunc {
			return h
		}
	}

	logger.Info("chaos delay: slow message consumption enabled",
		slog.Duration("max_delay", maxDelay),
	)

	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			delay := time.Duration(rand.Int63n(int64(maxDelay)))
			if delay > 0 {
				logger.Debug("chaos delay: delaying message consumption",
					slog.String("message_uuid", msg.UUID),
					slog.String("message_type", msg.Metadata.Get("name")),
					slog.Duration("delay", delay),
				)
				time.Sleep(delay)
			}
			return h(msg)
		}
	}
}
