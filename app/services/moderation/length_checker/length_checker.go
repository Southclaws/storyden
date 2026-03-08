package length_checker

import (
	"context"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type LengthChecker struct {
	settingsRepo *settings.SettingsRepository
	bus          *pubsub.Bus

	mu                  sync.RWMutex
	threadBodyLengthMax int
	replyBodyLengthMax  int
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
		threadBodyLengthMax: 60000,
		replyBodyLengthMax:  10000,
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
			if maxThread, ok := moderation.ThreadBodyLengthMax.Get(); ok {
				l.threadBodyLengthMax = maxThread
			}
			if maxReply, ok := moderation.ReplyBodyLengthMax.Get(); ok {
				l.replyBodyLengthMax = maxReply
			}
		}
	}

	return nil
}

func (l *LengthChecker) handleSettingsUpdate(ctx context.Context, event *rpc.EventSettingsUpdated) error {
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
		maxLength = l.threadBodyLengthMax
	case datagraph.KindReply:
		maxLength = l.replyBodyLengthMax
	default:
		return nil, fault.Wrap(
			fault.Newf("unsupported kind: %s", targetKind),
			fctx.With(ctx),
		)
	}

	if len(content.Plaintext()) > maxLength {
		return &checker.Result{
			Action: checker.ActionReject,
			Reason: "Content exceeds maximum allowed length",
		}, nil
	}

	return &checker.Result{
		Action: checker.ActionAllow,
	}, nil
}
