package robot

import (
	"encoding/xml"
	"strings"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
)

// MessageSpeakerMetadataKey identifies the untrusted display label projected
// into model context for a message's speaker.
const MessageSpeakerMetadataKey = "storyden.message.speaker"

// ContentWithSpeaker builds model content with a caller-formatted speaker label.
func ContentWithSpeaker(role genai.Role, content, speaker string) *genai.Content {
	return withMessageSpeaker(genai.NewContentFromText(content, role), speaker)
}

// ContentWithStorydenSpeaker attributes content to a human Storyden account.
func ContentWithStorydenSpeaker(content *genai.Content, author *account.Account) *genai.Content {
	if author == nil || author.Kind != account.AccountKindHuman {
		return content
	}
	return withMessageSpeaker(content, storydenAccountSpeaker(author))
}

func projectMessageSpeakersBeforeModel() llmagent.BeforeModelCallback {
	return func(_ agent.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
		if req == nil {
			return nil, nil
		}

		contents := make([]*genai.Content, len(req.Contents))
		for i, content := range req.Contents {
			contents[i] = projectMessageSpeaker(content)
		}
		req.Contents = contents
		return nil, nil
	}
}

func withMessageSpeaker(content *genai.Content, speaker string) *genai.Content {
	if content == nil || strings.TrimSpace(speaker) == "" {
		return content
	}

	for _, part := range content.Parts {
		if part == nil {
			continue
		}
		if existing, ok := part.PartMetadata[MessageSpeakerMetadataKey].(string); ok && strings.TrimSpace(existing) != "" {
			return content
		}
	}

	parts := make([]*genai.Part, len(content.Parts))
	copy(parts, content.Parts)
	for i, part := range parts {
		if part == nil || part.Text == "" {
			continue
		}

		copyPart := *part
		copyPart.PartMetadata = make(map[string]any, len(part.PartMetadata)+1)
		for key, value := range part.PartMetadata {
			copyPart.PartMetadata[key] = value
		}
		copyPart.PartMetadata[MessageSpeakerMetadataKey] = speaker
		parts[i] = &copyPart

		copyContent := *content
		copyContent.Parts = parts
		return &copyContent
	}

	return content
}

func projectMessageSpeaker(content *genai.Content) *genai.Content {
	if content == nil || content.Role != genai.RoleUser {
		return content
	}

	parts := make([]*genai.Part, len(content.Parts))
	changed := false
	for i, part := range content.Parts {
		parts[i] = part
		if part == nil || part.Text == "" {
			continue
		}
		speaker, ok := part.PartMetadata[MessageSpeakerMetadataKey].(string)
		if !ok || strings.TrimSpace(speaker) == "" {
			continue
		}

		header := speakerHeader(speaker)
		if strings.HasPrefix(part.Text, header) {
			continue
		}
		copyPart := *part
		copyPart.Text = header + part.Text
		parts[i] = &copyPart
		changed = true
	}
	if !changed {
		return content
	}

	copyContent := *content
	copyContent.Parts = parts
	return &copyContent
}

func speakerHeader(speaker string) string {
	var escaped strings.Builder
	_ = xml.EscapeText(&escaped, []byte(speaker))
	escapedSpeaker := strings.NewReplacer("&#34;", `"`, "&#39;", `'`).Replace(escaped.String())

	var output strings.Builder
	output.WriteString("<speaker>\n")
	output.WriteString(escapedSpeaker)
	output.WriteString("\n</speaker>\n")
	return output.String()
}

func storydenAccountSpeaker(author *account.Account) string {
	if author == nil {
		return ""
	}

	name := strings.TrimSpace(author.Name)
	handle := strings.TrimSpace(author.Handle)
	switch {
	case name != "" && handle != "":
		return name + " (@" + handle + ")"
	case handle != "":
		return "@" + handle
	default:
		return name
	}
}
