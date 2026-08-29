package robot

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
)

func TestProjectMessageSpeakerPreservesStoredContent(t *testing.T) {
	original := ContentWithSpeaker(genai.RoleUser, "hello from Discord", `Alice <Admin> & "owner"`)

	projected := projectMessageSpeaker(original)
	projectedAgain := projectMessageSpeaker(projected)

	require.Len(t, original.Parts, 1)
	assert.Equal(t, "hello from Discord", original.Parts[0].Text)
	assert.Equal(t, `Alice <Admin> & "owner"`, original.Parts[0].PartMetadata[MessageSpeakerMetadataKey])
	require.Len(t, projected.Parts, 1)
	assert.Equal(t, `<speaker>
Alice &lt;Admin&gt; &amp; "owner"
</speaker>
hello from Discord`, projected.Parts[0].Text)
	assert.Equal(t, projected.Parts[0].Text, projectedAgain.Parts[0].Text)
	assert.NotSame(t, original, projected)
	assert.NotSame(t, original.Parts[0], projected.Parts[0])
}

func TestProjectMessageSpeakerLeavesOrdinaryContentUntouched(t *testing.T) {
	original := genai.NewContentFromText("ordinary message", genai.RoleModel)

	assert.Same(t, original, projectMessageSpeaker(original))
}

func TestProjectMessageSpeakerLeavesAttributedAssistantContentUntouched(t *testing.T) {
	original := ContentWithSpeaker(
		genai.RoleModel,
		"Southclaws is proof that robots need supervision.",
		"makie the friendly robot (@makie the friendly robot)",
	)

	projected := projectMessageSpeaker(original)

	assert.Same(t, original, projected)
	assert.Equal(t, "Southclaws is proof that robots need supervision.", projected.Parts[0].Text)
	assert.Equal(
		t,
		"makie the friendly robot (@makie the friendly robot)",
		projected.Parts[0].PartMetadata[MessageSpeakerMetadataKey],
	)
}

func TestStorydenSpeakerUsesDisplayNameAndHandle(t *testing.T) {
	original := genai.NewContentFromText("hello from Storyden", genai.RoleUser)
	author := &account.Account{
		Kind:   account.AccountKindHuman,
		Name:   "Barney",
		Handle: "southclaws",
	}

	attributed := ContentWithStorydenSpeaker(original, author)
	projected := projectMessageSpeaker(attributed)

	assert.Equal(t, "hello from Storyden", original.Parts[0].Text)
	assert.Nil(t, original.Parts[0].PartMetadata)
	assert.Equal(t, "Barney (@southclaws)", attributed.Parts[0].PartMetadata[MessageSpeakerMetadataKey])
	assert.Equal(t, `<speaker>
Barney (@southclaws)
</speaker>
hello from Storyden`, projected.Parts[0].Text)
}

func TestStorydenSpeakerDoesNotOverrideExternalSpeaker(t *testing.T) {
	original := ContentWithSpeaker(genai.RoleUser, "hello", "Discord Barney (@southclaws)")
	author := &account.Account{
		Kind:   account.AccountKindHuman,
		Name:   "Plugin Bot",
		Handle: "discord-bridge",
	}

	attributed := ContentWithStorydenSpeaker(original, author)

	assert.Same(t, original, attributed)
	assert.Equal(t, "Discord Barney (@southclaws)", attributed.Parts[0].PartMetadata[MessageSpeakerMetadataKey])
}

func TestStorydenSpeakerIgnoresBotAccounts(t *testing.T) {
	original := genai.NewContentFromText("automated input", genai.RoleUser)
	author := &account.Account{
		Kind:   account.AccountKindBot,
		Name:   "Bridge",
		Handle: "bridge",
	}

	assert.Same(t, original, ContentWithStorydenSpeaker(original, author))
}

func TestMapToADKEventsAddsStorydenSpeakerWithoutMutatingStoredEvent(t *testing.T) {
	storedContent := genai.NewContentFromText("hello from a shared session", genai.RoleUser)
	messages := []*robotresource.Message{
		{
			Author: opt.New(&account.Account{
				Kind:   account.AccountKindHuman,
				Name:   "Barney",
				Handle: "southclaws",
			}),
			Event: newSessionEvent(storedContent),
		},
	}

	events := mapToADKEventsFromMessages(messages)

	require.Equal(t, 1, events.Len())
	hydrated := events.At(0).LLMResponse.Content
	assert.Equal(t, "hello from a shared session", storedContent.Parts[0].Text)
	assert.Nil(t, storedContent.Parts[0].PartMetadata)
	assert.Equal(t, "Barney (@southclaws)", hydrated.Parts[0].PartMetadata[MessageSpeakerMetadataKey])
	assert.NotSame(t, storedContent, hydrated)
}

func newSessionEvent(content *genai.Content) adksession.Event {
	return adksession.Event{
		LLMResponse: model.LLMResponse{
			Content: content,
		},
	}
}
