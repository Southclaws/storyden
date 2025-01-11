package asker

import (
	"context"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/question"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type cachedAsker struct {
	logger    *zap.Logger
	asker     semdex.Asker
	questions *question.Repository
}

func newCachedAsker(
	logger *zap.Logger,
	asker semdex.Asker,
	questions *question.Repository,
) (semdex.Asker, error) {
	return &cachedAsker{
		logger:    logger,
		asker:     asker,
		questions: questions,
	}, nil
}

func (a *cachedAsker) Ask(ctx context.Context, q string) (func(yield func(string, error) bool), error) {
	cached, err := a.questions.GetByQuerySlug(ctx, q)
	if err == nil {
		return a.cachedResult(ctx, cached)
	}

	return a.livePrompt(ctx, q)
}

func (a *cachedAsker) cachedResult(ctx context.Context, q *question.Question) (func(yield func(string, error) bool), error) {
	md, err := htmltomarkdown.ConvertNode(q.Result.HTMLTree())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	chunks := strings.SplitAfter(string(md), " ")

	return func(yield func(string, error) bool) {
		for _, ch := range chunks {
			select {
			case <-ctx.Done():
				return

			default:
				if !yield(ch, nil) {
					return
				}
			}
			time.Sleep(time.Millisecond * 10)
		}
	}, nil
}

func (a *cachedAsker) livePrompt(ctx context.Context, q string) (func(yield func(string, error) bool), error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	iter, err := a.asker.Ask(ctx, q)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return func(yield func(string, error) bool) {
		acc := []string{}

		defer func() {
			err := a.cacheResult(ctx, accountID, q, acc)
			if err != nil {
				a.logger.Error("failed to cache result", zap.Error(err))
			}
		}()

		for chunk, err := range iter {
			if err != nil {
				yield("", err)
				return
			}

			acc = append(acc, chunk)
			if !yield(chunk, nil) {
				return
			}
		}
	}, nil
}

func (a *cachedAsker) cacheResult(ctx context.Context, accountID account.AccountID, q string, chunks []string) error {
	result := strings.Join(chunks, "")

	acc, err := datagraph.NewRichTextFromMarkdown(result)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = a.questions.Store(ctx, accountID, q, acc)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
