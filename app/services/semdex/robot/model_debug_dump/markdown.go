package model_debug_dump

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

func render(input snapshot, created time.Time) (string, error) {
	if input.request == nil {
		return "", fmt.Errorf("model debug dump requires a request")
	}
	var out strings.Builder
	metadata := map[string]any{
		"session_id": input.sessionID, "invocation_id": input.invocationID,
		"agent": input.agentName, "model": input.request.Model, "created_at": created.Format(time.RFC3339Nano),
	}
	if err := writeJSONSection(&out, "DUMP_METADATA", metadata); err != nil {
		return "", err
	}
	var system *genai.Content
	if input.request.Config != nil {
		if len(input.request.Config.Tools) > 0 {
			if err := writeJSONSection(&out, "TOOL_DECLARATIONS", input.request.Config.Tools); err != nil {
				return "", err
			}
		}
		system = input.request.Config.SystemInstruction
	}
	if err := writeContent(&out, "SYSTEM_PROMPT", system); err != nil {
		return "", err
	}
	for i, content := range input.request.Contents {
		if err := writeContent(&out, fmt.Sprintf("MESSAGE %d", i+1), content); err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

func writeContent(out *strings.Builder, section string, content *genai.Content) error {
	writeMarker(out, "BEGIN "+section)
	if content == nil {
		writeMarker(out, "NO_CONTENT")
		writeMarker(out, "END "+section)
		return nil
	}
	writeMarker(out, fmt.Sprintf("ADK_CONTENT_ROLE %q", content.Role))
	for i, part := range content.Parts {
		writeMarker(out, fmt.Sprintf("BEGIN PART %d", i+1))
		if err := writePart(out, part); err != nil {
			return err
		}
		writeMarker(out, fmt.Sprintf("END PART %d", i+1))
	}
	writeMarker(out, "END "+section)
	return nil
}

func writePart(out *strings.Builder, part *genai.Part) error {
	if part == nil {
		writeMarker(out, "NO_PART")
		return nil
	}
	if part.Text != "" {
		writeFence(out, "text", part.Text)
	}
	metadata := *part
	metadata.Text = ""
	raw, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	if string(raw) != "{}" {
		writeFence(out, "json", string(raw))
	}
	return nil
}

func writeJSONSection(out *strings.Builder, section string, value any) error {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	writeMarker(out, "BEGIN "+section)
	writeFence(out, "json", string(raw))
	writeMarker(out, "END "+section)
	return nil
}

func writeMarker(out *strings.Builder, marker string) {
	fmt.Fprintf(out, "<| ROBOT_MODEL_DEBUG_DUMP %s |>\n\n", marker)
}

func writeFence(out *strings.Builder, language, text string) {
	fence := "```"
	for strings.Contains(text, fence) {
		fence += "`"
	}
	fmt.Fprintf(out, "%s%s\n%s\n%s\n\n", fence, language, text, fence)
}
