package robotprojection

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	adksession "google.golang.org/adk/session"
	"google.golang.org/adk/tool/toolconfirmation"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/presentation"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func HiddenConfirmationToolCallIDs(messages []*robot.Message) map[string]bool {
	ids := make(map[string]bool)

	for _, message := range messages {
		if message == nil || message.Event.LLMResponse.Content == nil {
			continue
		}

		for _, part := range message.Event.LLMResponse.Content.Parts {
			if part == nil || part.FunctionResponse == nil || part.FunctionResponse.ID == "" {
				continue
			}

			if IsConfirmationBlockedResponse(part.FunctionResponse) {
				ids[part.FunctionResponse.ID] = true
			}
		}
	}

	return ids
}

// ADKEventToUIMessageParts projects persisted ADK event content into AI SDK
// UIMessage parts. It mirrors the live SSE projection so hydration cannot
// expose internal confirmation wrapper events that the live stream hides.
func ADKEventToUIMessageParts(event adksession.Event, hiddenToolCallIDs map[string]bool) ([]openapi.UIMessagePart, error) {
	if event.LLMResponse.Content == nil {
		return []openapi.UIMessagePart{}, nil
	}

	var parts []openapi.UIMessagePart
	confirmationOriginalIDs := confirmationOriginalCallIDs(event.LLMResponse.Content.Parts)
	confirmationBlockedIDs := confirmationBlockedCallIDs(event.LLMResponse.Content.Parts)

	for _, adkPart := range event.LLMResponse.Content.Parts {
		if adkPart == nil {
			continue
		}

		if adkPart.Text != "" {
			textParts, err := TextToPresentationUIParts(event, adkPart.Text)
			if err != nil {
				return nil, err
			}
			parts = append(parts, textParts...)
		}

		if adkPart.FunctionCall != nil {
			if hiddenToolCallIDs[adkPart.FunctionCall.ID] {
				continue
			}

			if confirmationBlockedIDs[adkPart.FunctionCall.ID] {
				continue
			}

			if confirmationOriginalIDs[adkPart.FunctionCall.ID] && adkPart.FunctionCall.Name != toolconfirmation.FunctionCallName {
				continue
			}

			uiPart, err := FunctionCallToUIPart(UnwrapConfirmationFunctionCall(adkPart.FunctionCall))
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}

		if adkPart.FunctionResponse != nil {
			if hiddenToolCallIDs[adkPart.FunctionResponse.ID] {
				continue
			}

			if confirmationBlockedIDs[adkPart.FunctionResponse.ID] {
				continue
			}

			if confirmationOriginalIDs[adkPart.FunctionResponse.ID] {
				continue
			}

			uiPart, err := FunctionResponseToUIPart(adkPart.FunctionResponse)
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}
	}

	return parts, nil
}

func TextToPresentationUIParts(event adksession.Event, text string) ([]openapi.UIMessagePart, error) {
	if event.Author == "user" {
		part, err := TextUIPart(text)
		if err != nil {
			return nil, err
		}
		return []openapi.UIMessagePart{part}, nil
	}

	presentationParts := presentation.Parse(text)
	if len(presentationParts) == 0 {
		return nil, nil
	}

	parts := []openapi.UIMessagePart{}
	for _, presentationPart := range presentationParts {
		switch presentationPart.Kind {
		case presentation.PartText:
			if presentationPart.Text == "" {
				continue
			}
			part, err := TextUIPart(presentationPart.Text)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)

		case presentation.PartRenderCard:
			data := presentation.NewRenderCardData(presentationPart.Ref)
			part, err := DataUIPart(presentation.DataRenderCard, data)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		}
	}

	return parts, nil
}

func TextUIPart(text string) (openapi.UIMessagePart, error) {
	textPart := openapi.TextUIPart{
		Type:  openapi.TextUIPartType("text"),
		Text:  text,
		State: ptr(openapi.TextUIPartState("done")),
	}
	var uiPart openapi.UIMessagePart
	if err := uiPart.FromTextUIPart(textPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create text part: %w", err)
	}

	return uiPart, nil
}

func DataUIPart(partType string, data any) (openapi.UIMessagePart, error) {
	dataPart := openapi.DataPart{
		Type: partType,
		Data: data,
	}
	var uiPart openapi.UIMessagePart
	if err := uiPart.FromDataPart(dataPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create data part: %w", err)
	}
	uiPart.Type = openapi.UIMessagePartType(partType)

	return uiPart, nil
}

func FunctionCallToUIPart(fc *genai.FunctionCall) (openapi.UIMessagePart, error) {
	inputAvailable := openapi.ToolUIPartInputAvailable{
		ToolCallId: fc.ID,
		ToolName:   fc.Name,
		State:      openapi.InputAvailable,
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

	uiPart.Type = openapi.UIMessagePartType("tool-" + fc.Name)

	return uiPart, nil
}

func FunctionResponseToUIPart(fr *genai.FunctionResponse) (openapi.UIMessagePart, error) {
	outputAvailable := openapi.ToolUIPartOutputAvailable{
		ToolCallId: fr.ID,
		ToolName:   fr.Name,
		State:      openapi.OutputAvailable,
		Input:      map[string]interface{}{},
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

	uiPart.Type = openapi.UIMessagePartType("tool-" + fr.Name)

	return openapi.UIMessagePart(uiPart), nil
}

func PresentationStreamParts(event *adksession.Event, fallbackTextID string) []openapi.StreamPart {
	if event == nil || event.LLMResponse.Content == nil {
		return nil
	}

	var streamParts []openapi.StreamPart
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || strings.TrimSpace(part.Text) == "" {
			continue
		}

		for _, presentationPart := range presentation.Parse(part.Text) {
			switch presentationPart.Kind {
			case presentation.PartText:
				if presentationPart.Text == "" {
					continue
				}
				textID := uuid.NewString()
				if fallbackTextID != "" {
					textID = fallbackTextID
					fallbackTextID = ""
				}
				streamParts = append(streamParts, TextStreamParts(textID, presentationPart.Text)...)

			case presentation.PartRenderCard:
				data := presentation.NewRenderCardData(presentationPart.Ref)
				streamParts = append(streamParts, DataStreamPart(presentation.DataRenderCard, data))
			}
		}
	}

	return streamParts
}

func TextStreamParts(textID string, text string) []openapi.StreamPart {
	textStartPart := openapi.StreamPart{}
	_ = textStartPart.FromTextStartPart(openapi.TextStartPart{
		Id: textID,
	})

	textDeltaPart := openapi.StreamPart{}
	if err := textDeltaPart.FromTextDeltaPart(openapi.TextDeltaPart{
		Id:    textID,
		Delta: text,
	}); err != nil {
		return []openapi.StreamPart{textStartPart}
	}

	textEndPart := openapi.StreamPart{}
	_ = textEndPart.FromTextEndPart(openapi.TextEndPart{
		Id: textID,
	})

	return []openapi.StreamPart{textStartPart, textDeltaPart, textEndPart}
}

func DataStreamPart(partType string, data any) openapi.StreamPart {
	dataPart := openapi.StreamPart{}
	_ = dataPart.FromDataPart(openapi.DataPart{
		Type: partType,
		Data: data,
	})
	dataPart.Type = partType
	return dataPart
}

func FunctionCallStreamParts(fc *genai.FunctionCall) []openapi.StreamPart {
	if fc == nil {
		return nil
	}

	toolInputStartPart := openapi.StreamPart{}
	_ = toolInputStartPart.FromToolInputStartPart(openapi.ToolInputStartPart{
		ToolCallId: fc.ID,
		ToolName:   fc.Name,
	})

	parts := []openapi.StreamPart{toolInputStartPart}

	argsJSON, err := json.Marshal(fc.Args)
	if err == nil {
		toolInputDeltaPart := openapi.StreamPart{}
		_ = toolInputDeltaPart.FromToolInputDeltaPart(openapi.ToolInputDeltaPart{
			ToolCallId:     fc.ID,
			InputTextDelta: string(argsJSON),
		})
		parts = append(parts, toolInputDeltaPart)
	}

	toolInputAvailablePart := openapi.StreamPart{}
	_ = toolInputAvailablePart.FromToolInputAvailablePart(openapi.ToolInputAvailablePart{
		ToolCallId: fc.ID,
		ToolName:   fc.Name,
		Input:      fc.Args,
	})
	parts = append(parts, toolInputAvailablePart)

	return parts
}

func FunctionResponseStreamPart(fr *genai.FunctionResponse) (openapi.StreamPart, bool) {
	if fr == nil {
		return openapi.StreamPart{}, false
	}

	toolOutputAvailablePart := openapi.StreamPart{}
	_ = toolOutputAvailablePart.FromToolOutputAvailablePart(openapi.ToolOutputAvailablePart{
		ToolCallId: fr.ID,
		Output:     fr.Response,
	})
	return toolOutputAvailablePart, true
}

func UnwrapConfirmationFunctionCall(fc *genai.FunctionCall) *genai.FunctionCall {
	if fc == nil || fc.Name != toolconfirmation.FunctionCallName {
		return fc
	}

	original, err := toolconfirmation.OriginalCallFrom(fc)
	if err != nil {
		return fc
	}

	return &genai.FunctionCall{
		ID:   fc.ID,
		Name: original.Name,
		Args: original.Args,
	}
}

func IsConfirmationBlockedResponse(fr *genai.FunctionResponse) bool {
	if fr == nil || fr.Response == nil {
		return false
	}

	errValue, ok := fr.Response["error"]
	return ok && strings.Contains(fmt.Sprint(errValue), "requires confirmation")
}

func confirmationOriginalCallIDs(parts []*genai.Part) map[string]bool {
	ids := make(map[string]bool)

	for _, part := range parts {
		if part == nil || part.FunctionCall == nil || part.FunctionCall.Name != toolconfirmation.FunctionCallName {
			continue
		}

		original, err := toolconfirmation.OriginalCallFrom(part.FunctionCall)
		if err != nil || original.ID == "" {
			continue
		}

		ids[original.ID] = true
	}

	return ids
}

func confirmationBlockedCallIDs(parts []*genai.Part) map[string]bool {
	ids := make(map[string]bool)

	for _, part := range parts {
		if part == nil || part.FunctionResponse == nil || part.FunctionResponse.ID == "" {
			continue
		}

		if IsConfirmationBlockedResponse(part.FunctionResponse) {
			ids[part.FunctionResponse.ID] = true
		}
	}

	return ids
}

func ptr[T any](v T) *T {
	return &v
}
