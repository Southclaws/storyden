package robot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/pagination"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

type SessionNamer struct {
	logger      *slog.Logger
	prompter    ai.Prompter
	sessionRepo *robot_session.Repository
}

func NewSessionNamer(
	logger *slog.Logger,
	prompter ai.Prompter,
	sessionRepo *robot_session.Repository,
) *SessionNamer {
	return &SessionNamer{
		logger:      logger,
		prompter:    prompter,
		sessionRepo: sessionRepo,
	}
}

type SessionNameResponse struct {
	HasEnoughInfo bool   `json:"has_enough_info"`
	Name          string `json:"name"`
}

const namingPromptTemplate = `Given a user's first message in a chat session, decide if there is enough meaningful context to generate a concise session name.

Rules:
- Only set has_enough_info to true if the message contains a clear task, question, or topic
- Generic greetings like "hi", "hello", "hey" are NOT enough context
- Vague messages like "help me" without specifics are NOT enough context
- The name should be 2-5 words that capture the core topic or task
- Use lowercase for the name
- Be specific and descriptive

Examples of sufficient context:
- "How do I deploy my app to production?" -> {has_enough_info: true, name: "app deployment help"}
- "Debug memory leak in Go service" -> {has_enough_info: true, name: "go memory leak debugging"}
- "Explain how authentication works" -> {has_enough_info: true, name: "authentication explanation"}

Examples of insufficient context:
- "hi" -> {has_enough_info: false, name: ""}
- "hello there" -> {has_enough_info: false, name: ""}
- "help me" -> {has_enough_info: false, name: ""}

User's message:
%s
`

func (s *SessionNamer) MaybeNameSession(
	ctx context.Context,
	sessionID robotresource.SessionID,
	userMessage string,
) {
	sess, _, err := s.sessionRepo.Get(ctx, sessionID, pagination.NewPageParams(1, 5))
	if err != nil {
		s.logger.WarnContext(ctx, "failed to get session for naming",
			slog.String("session_id", xid.ID(sessionID).String()),
			slog.Any("error", err))
		return
	}

	if !strings.HasPrefix(sess.Name, "Untitled") {
		return
	}

	prompt := fmt.Sprintf(namingPromptTemplate, userMessage)

	result, err := ai.PromptObject(
		ctx,
		s.prompter,
		"Generate a session name based on the user's first message",
		prompt,
		SessionNameResponse{},
	)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to generate session name",
			slog.String("session_id", xid.ID(sessionID).String()),
			slog.Any("error", err))
		return
	}

	if !result.HasEnoughInfo || result.Name == "" {
		s.logger.DebugContext(ctx, "insufficient context for session naming",
			slog.String("session_id", xid.ID(sessionID).String()))
		return
	}

	if err := s.sessionRepo.UpdateName(ctx, sessionID, result.Name); err != nil {
		s.logger.WarnContext(ctx, "failed to update session name",
			slog.String("session_id", xid.ID(sessionID).String()),
			slog.String("proposed_name", result.Name),
			slog.Any("error", err))
		return
	}

	s.logger.InfoContext(ctx, "automatically named session",
		slog.String("session_id", xid.ID(sessionID).String()),
		slog.String("name", result.Name))
}
