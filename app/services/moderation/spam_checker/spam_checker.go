package spam_checker

import (
	"context"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
)

type SpamChecker struct {
	detector Detector
	enabled  bool
}

func NewSpamChecker(detector Detector) *SpamChecker {
	return &SpamChecker{
		detector: detector,
		enabled:  true,
	}
}

func (s *SpamChecker) Name() string {
	return "spam_detector"
}

func (s *SpamChecker) Enabled() bool {
	return s.enabled
}

func (s *SpamChecker) Check(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*checker.Result, error) {
	isSpam, err := s.detector.Detect(ctx, strings.NewReader(content.Plaintext()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if isSpam {
		return &checker.Result{
			Action: checker.ActionReport,
			Reason: "Content flagged as potential spam by spam detector",
		}, nil
	}

	return &checker.Result{
		Action: checker.ActionAllow,
	}, nil
}
