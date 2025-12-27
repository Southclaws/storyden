package length_checker

import (
	"context"
	"fmt"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type LengthChecker struct {
	settingsRepo *settings.SettingsRepository
	bus          *pubsub.Bus

	mu                  sync.RWMutex
	maxThreadBodyLength int
	maxReplyBodyLength  int
	enabled             bool
}

func NewLengthChecker(
	lc fx.Lifecycle,
	settingsRepo *settings.SettingsRepository,
	bus *pubsub.Bus,
) *LengthChecker {
	l := &LengthChecker{
		settingsRepo:        settingsRepo,
		bus:                 bus,
		maxThreadBodyLength: 60000,
		maxReplyBodyLength:  10000,
		enabled:             true,
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		if err := l.loadSettings(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err := pubsub.Subscribe(ctx, bus, "length_checker_settings_update", l.handleSettingsUpdate)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		return nil
	}))

	return l
}

func (l *LengthChecker) loadSettings(ctx context.Context) error {
	s, err := l.settingsRepo.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if services, ok := s.Services.Get(); ok {
		if moderation, ok := services.Moderation.Get(); ok {
			if maxThread, ok := moderation.MaxThreadBodyLength.Get(); ok {
				l.maxThreadBodyLength = maxThread
			}
			if maxReply, ok := moderation.MaxReplyBodyLength.Get(); ok {
				l.maxReplyBodyLength = maxReply
			}
		}
	}

	return nil
}

func (l *LengthChecker) handleSettingsUpdate(ctx context.Context, event *message.EventSettingsUpdated) error {
	if err := l.loadSettings(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (l *LengthChecker) Name() string {
	return "length_checker"
}

func (l *LengthChecker) Enabled() bool {
	return l.enabled
}

func (l *LengthChecker) Check(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*checker.Result, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var maxLength int
	switch targetKind {
	case datagraph.KindThread:
		maxLength = l.maxThreadBodyLength
	case datagraph.KindReply:
		maxLength = l.maxReplyBodyLength
	default:
		return nil, fault.Wrap(
			fault.Newf("unsupported kind: %s", targetKind),
			fctx.With(ctx),
		)
	}

	if len(content.Plaintext()) > maxLength {
		return &checker.Result{
			RequiresReview: true,
			Reason:         fmt.Sprintf("Content exceeds maximum length of %d characters", maxLength),
		}, nil
	}

	return &checker.Result{
		RequiresReview: false,
	}, nil
}
