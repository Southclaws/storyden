package content_policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/moderation/spam"
)

var (
	ErrContentTooLong     = fault.New("content too long", ftag.With(ftag.InvalidArgument))
	ErrContentFlaggedSpam = fault.New("content flagged as spam", ftag.With(ftag.InvalidArgument))
)

type Manager struct {
	spamDetector spam.Detector
}

func New(d spam.Detector) *Manager {
	return &Manager{
		spamDetector: d,
	}
}

func (m *Manager) CheckContent(ctx context.Context, c datagraph.Content) error {
	if len(c.Plaintext()) > post.MaxPostLength {
		message := fmt.Sprintf("Content must be less than %d characters", post.MaxPostLength)

		return fault.Wrap(ErrContentTooLong,
			fctx.With(ctx),
			fmsg.WithDesc("too long", message))
	}

	isSpam, err := m.spamDetector.Detect(ctx, strings.NewReader(c.Plaintext()))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if isSpam {
		message := "Your post has been flagged as potential spam."

		return fault.Wrap(ErrContentFlaggedSpam,
			fctx.With(ctx),
			fmsg.WithDesc("flagged", message))

	}

	return nil
}
