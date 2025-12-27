package word_checker

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

type WordChecker struct {
	settingsRepo *settings.SettingsRepository
	bus          *pubsub.Bus

	mu                   sync.RWMutex
	blockListNormalized  map[string]string
	reportListNormalized map[string]string
	enabled              bool
}

func NewWordChecker(
	lc fx.Lifecycle,
	settingsRepo *settings.SettingsRepository,
	bus *pubsub.Bus,
) *WordChecker {
	w := &WordChecker{
		settingsRepo:         settingsRepo,
		bus:                  bus,
		blockListNormalized:  make(map[string]string),
		reportListNormalized: make(map[string]string),
		enabled:              true,
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		if err := w.loadSettings(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err := pubsub.Subscribe(ctx, bus, "word_checker_settings_update", w.handleSettingsUpdate)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		return nil
	}))

	return w
}

func (w *WordChecker) loadSettings(ctx context.Context) error {
	s, err := w.settingsRepo.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if services, ok := s.Services.Get(); ok {
		if moderation, ok := services.Moderation.Get(); ok {
			if blockList, ok := moderation.WordBlockList.Get(); ok {
				w.blockListNormalized = normalizeWordList(blockList)
			}
			if reportList, ok := moderation.WordReportList.Get(); ok {
				w.reportListNormalized = normalizeWordList(reportList)
			}
		}
	}

	return nil
}

func normalizeWordList(words []string) map[string]string {
	normalized := make(map[string]string)
	for _, word := range words {
		norm := strings.ToLower(strings.TrimSpace(word))
		if norm != "" {
			normalized[norm] = word
		}
	}
	return normalized
}

func (w *WordChecker) handleSettingsUpdate(ctx context.Context, event *message.EventSettingsUpdated) error {
	if err := w.loadSettings(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (w *WordChecker) Name() string {
	return "word_checker"
}

func (w *WordChecker) Enabled() bool {
	return w.enabled
}

func (w *WordChecker) Check(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*checker.Result, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.blockListNormalized) == 0 && len(w.reportListNormalized) == 0 {
		return &checker.Result{
			Action: checker.ActionAllow,
		}, nil
	}

	plaintext := strings.ToLower(name + " " + content.Plaintext())

	for normalized := range w.blockListNormalized {
		if strings.Contains(plaintext, normalized) {
			return &checker.Result{
				Action: checker.ActionReject,
				Reason: "Content violates community guidelines",
			}, nil
		}
	}

	for normalized, original := range w.reportListNormalized {
		if strings.Contains(plaintext, normalized) {
			return &checker.Result{
				Action: checker.ActionReport,
				Reason: "Content contains blocked word: " + original,
			}, nil
		}
	}

	return &checker.Result{
		Action: checker.ActionAllow,
	}, nil
}
