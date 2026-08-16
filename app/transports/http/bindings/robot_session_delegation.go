package bindings

import (
	"context"
	"log/slog"
	"strings"

	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
)

const delegationDataPartType = "data-delegation"

type delegationStatus string

const (
	delegationStatusRunning   delegationStatus = "running"
	delegationStatusCompleted delegationStatus = "completed"
	delegationStatusFailed    delegationStatus = "failed"
)

type delegationData struct {
	CallID   string                 `json:"callId"`
	Robot    openapi.RobotReference `json:"robot"`
	Request  string                 `json:"request"`
	Status   delegationStatus       `json:"status"`
	Messages []openapi.UIMessage    `json:"messages"`
	Error    string                 `json:"error,omitempty"`
}

type delegationStream struct {
	ctx          context.Context
	robotQuerier *robot_querier.Querier
	toolMetadata robotprojection.ToolMetadataResolver
	emitter      partEmitter
	logger       *slog.Logger
	active       map[string]*delegationData
}

func newDelegationStream(
	ctx context.Context,
	robotQuerier *robot_querier.Querier,
	toolMetadata robotprojection.ToolMetadataResolver,
	emitter partEmitter,
	logger *slog.Logger,
) *delegationStream {
	return &delegationStream{
		ctx:          ctx,
		robotQuerier: robotQuerier,
		toolMetadata: toolMetadata,
		emitter:      emitter,
		logger:       logger,
		active:       make(map[string]*delegationData),
	}
}

func (d *delegationStream) Start(call *genai.FunctionCall) bool {
	if call == nil || call.ID == "" {
		return false
	}

	robotID, ok := robot_ref.IDFromAgentName(call.Name)
	if !ok {
		return false
	}

	request, _ := call.Args["request"].(string)
	data := &delegationData{
		CallID:   call.ID,
		Robot:    d.robotReference(robotID),
		Request:  strings.TrimSpace(request),
		Status:   delegationStatusRunning,
		Messages: []openapi.UIMessage{},
	}
	d.active[call.ID] = data
	d.emit(data)
	return true
}

func (d *delegationStream) AppendEvent(event *adksession.Event) bool {
	if event == nil || event.IsolationScope == "" {
		return false
	}

	data, ok := d.active[event.IsolationScope]
	if !ok {
		robotID, isRobot := robot_ref.IDFromAgentName(event.Author)
		if !isRobot {
			return true
		}
		data = &delegationData{
			CallID:   event.IsolationScope,
			Robot:    d.robotReference(robotID),
			Status:   delegationStatusRunning,
			Messages: []openapi.UIMessage{},
		}
		d.active[event.IsolationScope] = data
	}

	projectedEvent := delegationPresentationEvent(event)
	parts, err := robotprojection.ADKEventToUIMessageParts(projectedEvent, nil, d.toolMetadata)
	if err != nil {
		d.logger.Warn("failed to project delegated Robot event",
			slog.String("call_id", event.IsolationScope),
			slog.String("error", err.Error()))
		return true
	}
	if len(parts) == 0 {
		return true
	}

	role := openapi.UIMessageRoleAssistant
	if event.Author == "user" {
		role = openapi.UIMessageRoleUser
	}
	data.Messages = append(data.Messages, openapi.UIMessage{
		Id:    event.ID,
		Role:  role,
		Parts: parts,
	})
	d.emit(data)
	return true
}

func delegationPresentationEvent(event *adksession.Event) adksession.Event {
	projected := *event
	if event.LLMResponse.Content == nil {
		return projected
	}
	content := *event.LLMResponse.Content
	content.Parts = make([]*genai.Part, 0, len(event.LLMResponse.Content.Parts))
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil {
			continue
		}
		if part.FunctionCall != nil && part.FunctionCall.Name == robot.UnattendedFinishToolName() {
			continue
		}
		if part.FunctionResponse != nil && part.FunctionResponse.Name == robot.UnattendedFinishToolName() {
			continue
		}
		content.Parts = append(content.Parts, part)
	}
	projected.LLMResponse.Content = &content
	return projected
}

func (d *delegationStream) Complete(response *genai.FunctionResponse) bool {
	if response == nil || response.ID == "" {
		return false
	}
	data, ok := d.active[response.ID]
	if !ok {
		robotID, isRobot := robot_ref.IDFromAgentName(response.Name)
		if !isRobot {
			return false
		}
		data = &delegationData{
			CallID: response.ID, Robot: d.robotReference(robotID),
			Status: delegationStatusRunning, Messages: []openapi.UIMessage{},
		}
		d.active[response.ID] = data
	}

	status, _ := response.Response["status"].(string)
	if status == "pending" {
		d.emit(data)
		return true
	}
	if status == "failed" || status == "cancelled" || status == "blocked" {
		data.Status = delegationStatusFailed
		data.Error, _ = response.Response["summary"].(string)
		d.emit(data)
		return true
	}
	data.Status = delegationStatusCompleted
	d.emit(data)
	return true
}

func (d *delegationStream) Fail(message string) {
	for _, data := range d.active {
		if data.Status != delegationStatusRunning {
			continue
		}
		data.Status = delegationStatusFailed
		data.Error = message
		d.emit(data)
	}
}

func (d *delegationStream) CompleteRunning() {
	for _, data := range d.active {
		if data.Status != delegationStatusRunning {
			continue
		}
		data.Status = delegationStatusCompleted
		d.emit(data)
	}
}

func (d *delegationStream) robotReference(id robot_ref.ID) openapi.RobotReference {
	reference := openapi.RobotReference{
		Id:   openapi.Identifier(id.String()),
		Name: "Delegated Robot",
	}

	resolved, err := d.robotQuerier.Get(d.ctx, id)
	if err != nil {
		d.logger.Warn("failed to resolve delegated Robot",
			slog.String("robot_id", id.String()),
			slog.String("error", err.Error()))
		return reference
	}
	reference.Name = resolved.Name
	return reference
}

func (d *delegationStream) emit(data *delegationData) {
	_ = d.emitter.Send(robotprojection.DataStreamPartWithID(delegationDataPartType, data.CallID, data))
}
