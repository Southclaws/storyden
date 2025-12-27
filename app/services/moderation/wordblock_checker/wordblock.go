package wordblock_checker

import (
	"context"
	"strings"
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

type WordBlockChecker struct {
	settingsRepo *settings.SettingsRepository
	bus          *pubsub.Bus

	mu            sync.RWMutex
	blockedWords  []string
	normalizedMap map[string]string
	enabled       bool
}

func NewWordBlockChecker(
	lc fx.Lifecycle,
	settingsRepo *settings.SettingsRepository,
	bus *pubsub.Bus,
) *WordBlockChecker {
	w := &WordBlockChecker{
		settingsRepo:  settingsRepo,
		bus:           bus,
		blockedWords:  []string{},
		normalizedMap: make(map[string]string),
		enabled:       true,
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		if err := w.loadSettings(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err := pubsub.Subscribe(ctx, bus, "wordblock_checker_settings_update", w.handleSettingsUpdate)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		return nil
	}))

	return w
}

func (w *WordBlockChecker) loadSettings(ctx context.Context) error {
	s, err := w.settingsRepo.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if services, ok := s.Services.Get(); ok {
		if moderation, ok := services.Moderation.Get(); ok {
			if wordBlacklist, ok := moderation.WordBlocklist.Get(); ok {
				w.blockedWords = wordBlacklist
				w.buildNormalizedMap()
			}
		}
	}

	return nil
}

func (w *WordBlockChecker) buildNormalizedMap() {
	w.normalizedMap = make(map[string]string)
	for _, word := range w.blockedWords {
		normalized := strings.ToLower(strings.TrimSpace(word))
		if normalized != "" {
			w.normalizedMap[normalized] = word
		}
	}
}

func (w *WordBlockChecker) handleSettingsUpdate(ctx context.Context, event *message.EventSettingsUpdated) error {
	if err := w.loadSettings(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (w *WordBlockChecker) Name() string {
	return "wordblock_checker"
}

func (w *WordBlockChecker) Enabled() bool {
	return w.enabled
}

func (w *WordBlockChecker) Check(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*checker.Result, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.normalizedMap) == 0 {
		return &checker.Result{
			RequiresReview: false,
		}, nil
	}

	// merge name with content, since we're just doing a dumb .contains() lol
	plaintext := strings.ToLower(name + " " + content.Plaintext())

	for normalized, original := range w.normalizedMap {
		if strings.Contains(plaintext, normalized) {
			return &checker.Result{
				RequiresReview: true,
				Reason:         "Content contains blocked word: " + original,
			}, nil
		}
	}

	return &checker.Result{
		RequiresReview: false,
	}, nil
}
