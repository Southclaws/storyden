package trail_action_registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/trail"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/session_coordinator"
	"github.com/Southclaws/storyden/internal/ent"
	entinput "github.com/Southclaws/storyden/internal/ent/robotsessioninput"
	entmessage "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	entturn "github.com/Southclaws/storyden/internal/ent/robotsessionturn"
)

type RobotActionAdapter struct {
	db          *ent.Client
	accounts    *account_repo.Repository
	agent       *robotservice.Agent
	coordinator *session_coordinator.Coordinator
	trails      *trail.Repository
}

func NewRobotActionAdapter(
	db *ent.Client,
	accounts *account_repo.Repository,
	agent *robotservice.Agent,
	coordinator *session_coordinator.Coordinator,
	trails *trail.Repository,
) *RobotActionAdapter {
	return &RobotActionAdapter{db: db, accounts: accounts, agent: agent, coordinator: coordinator, trails: trails}
}

func (a *RobotActionAdapter) Kind() trail.ActionKind { return trail.ActionKindRobotRun }

func (a *RobotActionAdapter) Validate(ctx context.Context, creator *account.Account, spec trail.ActionSpec) error {
	if err := creator.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionUseRobots); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	var config trail.RobotRunConfig
	if err := json.Unmarshal(spec.Config, &config); err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if config.Type != (trail.ActionKind{}) && config.Type != trail.ActionKindRobotRun {
		return fault.New("robot action type does not match its binding", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if strings.TrimSpace(config.RobotRef) == "" || strings.TrimSpace(config.Instruction) == "" {
		return fault.New("robot action requires a robot and an unattended instruction", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	return nil
}

func (a *RobotActionAdapter) Start(ctx context.Context, run *trail.ActionRun) error {
	if run.Trail == nil {
		return errors.New("trail action context is missing trail")
	}

	creator, err := a.accounts.GetRefByID(ctx, run.Trail.Creator.ID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := creator.RejectSuspended(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := creator.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionUseRobots); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	var config trail.RobotRunConfig
	if err := json.Unmarshal(run.Config, &config); err != nil {
		return fmt.Errorf("decode robot action: %w", err)
	}

	deterministicID := xid.ID(run.ID)
	sessionID := robotresource.SessionID(deterministicID)
	runNumber, err := a.trails.RunNumber(ctx, run.RunID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := a.agent.PrepareSession(
		ctx,
		config.RobotRef,
		creator.ID.String(),
		sessionID.String(),
		robot_session.WithName(trailRobotSessionName(run.Trail.Name, runNumber)),
		robot_session.WithTrailOrigin(deterministicID),
	); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if _, err := a.coordinator.Enqueue(
		ctx,
		robotresource.InputID(deterministicID),
		config.RobotRef,
		creator.ID.String(),
		sessionID.String(),
		genai.NewContentFromText(config.Instruction, genai.RoleUser),
		trailInvocationContext(run),
		robotservice.RunOptions{Mode: robotservice.ModeUnattended, Source: robotservice.SourceScheduled},
	); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	target, err := encodeRobotInvocation(deterministicID)
	if err != nil {
		return err
	}

	return a.trails.SetActionRunTarget(ctx, run.ID, *run.LeaseToken, target)
}

func trailInvocationContext(run *trail.ActionRun) robotservice.InvocationContext {
	invocation := robotservice.InvocationContext{
		robotservice.InvocationContextKeyTrailID:          run.Trail.ID.String(),
		robotservice.InvocationContextKeyTrailRunID:       run.RunID.String(),
		robotservice.InvocationContextKeyTrailActionRunID: run.ID.String(),
	}

	if run.Run != nil && len(run.Run.TriggerPayload) > 0 {
		invocation[robotservice.InvocationContextKeyTrailTrigger] = json.RawMessage(run.Run.TriggerPayload)

		var event trail.TriggerEvent
		if json.Unmarshal(run.Run.TriggerPayload, &event) == nil {
			invocation[robotservice.InvocationContextKeyTrailTriggerKind] = event.Kind.String()
		}
	}

	return invocation
}

func trailRobotSessionName(trailName string, runNumber int) string {
	return fmt.Sprintf("%s (Run %d)", trailName, runNumber)
}

func (a *RobotActionAdapter) Reconcile(ctx context.Context, run *trail.ActionRun) (*ActionResult, error) {
	input, err := a.db.RobotSessionInput.Query().Where(entinput.IDEQ(xid.ID(run.ID))).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if input.TurnID == nil {
		return nil, nil
	}

	turn, err := a.db.RobotSessionTurn.Query().Where(entturn.IDEQ(*input.TurnID)).Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sessionID := xid.ID(run.ID)
	target, err := encodeRobotInvocation(sessionID)
	if err != nil {
		return nil, err
	}

	if err := a.trails.SetActionRunTarget(ctx, run.ID, *run.LeaseToken, target); err != nil {
		return nil, err
	}

	switch turn.Status {
	case entturn.StatusQueued, entturn.StatusRunning:
		return nil, nil
	case entturn.StatusBlocked:
		return &ActionResult{Status: trail.ActionRunStatusBlocked, ErrorText: "Robot run needs input or approval"}, nil
	case entturn.StatusFailed:
		message := "Robot run failed"
		if turn.ErrorText != nil && strings.TrimSpace(*turn.ErrorText) != "" {
			message = *turn.ErrorText
		}

		return &ActionResult{Status: trail.ActionRunStatusFailed, ErrorText: message}, nil
	case entturn.StatusCancelled:
		return &ActionResult{Status: trail.ActionRunStatusCancelled}, nil
	case entturn.StatusCompleted:
		output, status, err := a.readFinish(ctx, sessionID, *input.TurnID)
		if err != nil {
			return &ActionResult{Status: trail.ActionRunStatusFailed, ErrorText: err.Error()}, nil
		}

		encoded, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}

		return &ActionResult{Status: status, Output: encoded}, nil
	default:
		return &ActionResult{Status: trail.ActionRunStatusFailed, ErrorText: "Robot run has an unknown status"}, nil
	}
}

func (a *RobotActionAdapter) readFinish(ctx context.Context, sessionID, turnID xid.ID) (trail.RobotInvocationOutput, trail.ActionRunStatus, error) {
	messages, err := a.db.RobotSessionMessage.Query().Where(
		entmessage.SessionIDEQ(sessionID), entmessage.TurnIDEQ(turnID),
	).Order(entmessage.BySequence()).All(ctx)
	if err != nil {
		return trail.RobotInvocationOutput{}, trail.ActionRunStatus{}, err
	}

	for i := len(messages) - 1; i >= 0; i-- {
		event := messages[i].EventData.ADK()
		if event.LLMResponse.Content == nil {
			continue
		}
		for _, part := range event.LLMResponse.Content.Parts {
			if part == nil || part.FunctionResponse == nil || part.FunctionResponse.Name != robotservice.UnattendedFinishToolName() {
				continue
			}
			encoded, err := json.Marshal(part.FunctionResponse.Response)
			if err != nil {
				return trail.RobotInvocationOutput{}, trail.ActionRunStatus{}, fmt.Errorf("encode robot_run_finish: %w", err)
			}

			var output trail.RobotInvocationOutput
			if err := json.Unmarshal(encoded, &output); err != nil || strings.TrimSpace(output.Summary) == "" {
				return trail.RobotInvocationOutput{}, trail.ActionRunStatus{}, errors.New("robot completed without a valid structured result")
			}

			switch output.Status {
			case trail.RobotInvocationOutputStatusCompleted:
				return output, trail.ActionRunStatusCompleted, nil
			case trail.RobotInvocationOutputStatusBlocked:
				return output, trail.ActionRunStatusBlocked, nil
			case trail.RobotInvocationOutputStatusFailed:
				return output, trail.ActionRunStatusFailed, nil
			default:
				return trail.RobotInvocationOutput{}, trail.ActionRunStatus{}, errors.New("robot completed with an invalid structured status")
			}
		}
	}

	return trail.RobotInvocationOutput{}, trail.ActionRunStatus{}, errors.New("robot completed without robot_run_finish")
}

func (a *RobotActionAdapter) Cancel(ctx context.Context, run *trail.ActionRun) error {
	input, err := a.db.RobotSessionInput.Query().Where(entinput.IDEQ(xid.ID(run.ID))).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return a.trails.FinishCancellation(ctx, run.ID)
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if input.TurnID != nil {
		if err := a.coordinator.Cancel(ctx, robotresource.SessionID(xid.ID(run.ID)), robotresource.TurnID(*input.TurnID)); err != nil {
			return err
		}
	}

	return a.trails.FinishCancellation(ctx, run.ID)
}

func encodeRobotInvocation(sessionID xid.ID) (json.RawMessage, error) {
	target := trail.RobotInvocation{
		Type:           trail.ActionKindRobotRun,
		RobotSessionID: sessionID.String(),
	}

	return json.Marshal(target)
}
