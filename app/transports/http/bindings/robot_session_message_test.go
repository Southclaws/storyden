package bindings

import (
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestSerialiseRobotSessionMessageIncludesImageAssets(t *testing.T) {
	firstID := asset.AssetID(xid.New())
	secondID := asset.AssetID(xid.New())
	message := &robot.Message{
		ID:        robot.MessageID(xid.New()),
		CreatedAt: time.Date(2026, time.September, 20, 20, 0, 0, 0, time.UTC),
		Assets: []*asset.Asset{
			{
				ID:       firstID,
				Name:     asset.NewExistingFilename(firstID, "first.png"),
				MIME:     mime.New("image/png"),
				Metadata: asset.Metadata{"width": float64(800), "height": float64(600)},
			},
			{
				ID:       secondID,
				Name:     asset.NewExistingFilename(secondID, "second.webp"),
				MIME:     mime.New("image/webp"),
				Metadata: asset.Metadata{"width": float64(1200), "height": float64(900)},
			},
		},
		Event: adksession.Event{
			Author: "user",
			LLMResponse: model.LLMResponse{
				Content: &genai.Content{Role: genai.RoleUser, Parts: []*genai.Part{{Text: "Compare these images."}}},
			},
		},
	}

	result, err := serialiseRobotSessionMessage(message, map[string]bool{}, nil)
	require.NoError(t, err)
	require.Len(t, result.Assets, 2)
	assert.Equal(t, firstID.String(), result.Assets[0].Id)
	assert.Equal(t, "/api/assets/"+message.Assets[0].Name.String(), result.Assets[0].Path)
	assert.Equal(t, float32(800), result.Assets[0].Width)
	assert.Equal(t, secondID.String(), result.Assets[1].Id)
	assert.Equal(t, "/api/assets/"+message.Assets[1].Name.String(), result.Assets[1].Path)
	assert.Equal(t, float32(900), result.Assets[1].Height)
}
