package bindings

import (
	"fmt"

	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

// serialiseRobotSessionMessage converts a robot.Message (containing ADK Event)
// into the Vercel AI SDK UIMessage format for the frontend.
func serialiseRobotSessionMessage(m *robot.Message) (openapi.RobotSessionMessage, error) {
	parts, err := serialiseADKEventToParts(m.Event)
	if err != nil {
		return openapi.RobotSessionMessage{}, err
	}

	role := determineMessageRole(m.Event)

	msg := openapi.RobotSessionMessage{
		Id:        m.ID.String(),
		Role:      openapi.RobotSessionMessageRole(role),
		Parts:     parts,
		CreatedAt: m.CreatedAt,
		Robot:     opt.Map(m.Robot, func(r *robot.Robot) openapi.Robot { return serialiseRobot(r) }).Ptr(),
		Author:    opt.Map(m.Author, func(a *account.Account) openapi.ProfileReference { return serialiseProfileReferenceFromAccount(*a) }).Ptr(),
	}

	return msg, nil
}

// determineMessageRole figures out the role based on the ADK Event structure.
// - If the event has no content or only function responses, it's from "user"
// - If the event has LLM-generated content (text, function calls), it's from "assistant"
func determineMessageRole(event session.Event) string {
	if event.LLMResponse.Content == nil {
		return "user"
	}

	for _, part := range event.LLMResponse.Content.Parts {
		if part != nil {
			if part.Text != "" || part.FunctionCall != nil {
				return "assistant"
			}
		}
	}

	return "user"
}

// serialiseADKEventToParts converts ADK Event.LLMResponse.Content.Parts into UIMessagePart[]
// Since UIMessagePart is a discriminated union represented as json.RawMessage,
// we need to marshal each concrete type to JSON.
func serialiseADKEventToParts(event session.Event) ([]openapi.UIMessagePart, error) {
	if event.LLMResponse.Content == nil {
		return []openapi.UIMessagePart{}, nil
	}

	var parts []openapi.UIMessagePart

	for _, adkPart := range event.LLMResponse.Content.Parts {
		if adkPart == nil {
			continue
		}

		if adkPart.Text != "" {
			textPart := openapi.TextUIPart{
				Type:  openapi.TextUIPartType("text"),
				Text:  adkPart.Text,
				State: ptr(openapi.TextUIPartState("done")), // Historical messages are always "done"
			}
			var uiPart openapi.UIMessagePart
			if err := uiPart.FromTextUIPart(textPart); err != nil {
				return nil, fmt.Errorf("create text part: %w", err)
			}
			parts = append(parts, uiPart)
		}

		if adkPart.FunctionCall != nil {
			uiPart, err := serialiseFunctionCallToPart(adkPart.FunctionCall)
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}

		if adkPart.FunctionResponse != nil {
			uiPart, err := serialiseFunctionResponseToPart(adkPart.FunctionResponse)
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}
	}

	return parts, nil
}

// serialiseFunctionCallToPart converts an ADK FunctionCall into a UIMessagePart (input-available state)
func serialiseFunctionCallToPart(fc *genai.FunctionCall) (openapi.UIMessagePart, error) {
	inputAvailable := openapi.ToolUIPartInputAvailable{
		Type:       openapi.ToolUIPartInputAvailableType("dynamic-tool"),
		ToolCallId: fc.ID,
		ToolName:   fc.Name,
		State:      openapi.ToolUIPartInputAvailableState("input-available"),
		Input:      fc.Args,
	}

	var toolPart openapi.ToolUIPart
	if err := toolPart.FromToolUIPartInputAvailable(inputAvailable); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create tool input part: %w", err)
	}

	var uiPart openapi.UIMessagePart
	if err := uiPart.FromToolUIPart(toolPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create UI message part from tool part: %w", err)
	}

	return uiPart, nil
}

// serialiseFunctionResponseToPart converts an ADK FunctionResponse into a UIMessagePart (output-available state)
func serialiseFunctionResponseToPart(fr *genai.FunctionResponse) (openapi.UIMessagePart, error) {
	outputAvailable := openapi.ToolUIPartOutputAvailable{
		Type:       openapi.ToolUIPartOutputAvailableType("dynamic-tool"),
		ToolCallId: fr.ID,
		ToolName:   fr.Name,
		State:      openapi.ToolUIPartOutputAvailableState("output-available"),
		Input:      fr.Response, // ADK stores the original input separately if needed
		Output:     fr.Response,
	}

	var toolPart openapi.ToolUIPart
	if err := toolPart.FromToolUIPartOutputAvailable(outputAvailable); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create tool output part: %w", err)
	}

	var uiPart openapi.UIMessagePart
	if err := uiPart.FromToolUIPart(toolPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create UI message part from tool part: %w", err)
	}

	return openapi.UIMessagePart(uiPart), nil
}

// ptr is a helper to get a pointer to a value
func ptr[T any](v T) *T {
	return &v
}
