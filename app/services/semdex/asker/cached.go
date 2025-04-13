package asker

import (
	"context"
	"log/slog"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/question"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type cachedAsker struct {
	logger    *slog.Logger
	asker     semdex.Asker
	questions *question.Repository
}

func newCachedAsker(
	logger *slog.Logger,
	asker semdex.Asker,
	questions *question.Repository,
) (semdex.Asker, error) {
	return &cachedAsker{
		logger:    logger,
		asker:     asker,
		questions: questions,
	}, nil
}

func (a *cachedAsker) Ask(ctx context.Context, q string, parent opt.Optional[xid.ID]) (semdex.AskResponseIterator, error) {
	cached, err := a.questions.GetByQuerySlug(ctx, q)
	if err == nil {
		return a.cachedResult(ctx, cached)
	}

	return a.livePrompt(ctx, q, parent)
}

func (a *cachedAsker) cachedResult(ctx context.Context, q *question.Question) (semdex.AskResponseIterator, error) {
	md, err := htmltomarkdown.ConvertNode(q.Result.HTMLTree())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	chunks := strings.SplitAfter(string(md), " ")

	// NOTE: Stream extractor is only run on cached results here, the live
	// prompter will run the stream extractor itself.
	return streamExtractor(func(yield func(string, error) bool) {
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
	}), nil
}

func (a *cachedAsker) livePrompt(ctx context.Context, q string, parentQuestionID opt.Optional[xid.ID]) (semdex.AskResponseIterator, error) {
	iter, err := a.asker.Ask(ctx, q, parentQuestionID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return func(yield func(semdex.AskResponseChunk, error) bool) {
		acc := []string{}

		defer func() {
			err := a.cacheResult(ctx, q, acc, parentQuestionID)
			if err != nil {
				a.logger.Error("failed to cache result", slog.String("error", err.Error()))
			}
		}()

		for chunk, err := range iter {
			if err != nil {
				yield(nil, err)
				return
			}

			if t, ok := chunk.(*semdex.AskResponseChunkText); ok {
				acc = append(acc, t.Chunk)
			}

			if !yield(chunk, nil) {
				return
			}
		}
	}, nil
}

func (a *cachedAsker) cacheResult(ctx context.Context, q string, chunks []string, parentQuestionID opt.Optional[xid.ID]) error {
	accountID := session.GetOptAccountID(ctx)

	result := strings.Join(chunks, "")

	acc, err := datagraph.NewRichTextFromMarkdown(result)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = a.questions.Store(ctx, q, acc, accountID, parentQuestionID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
